package tenants

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"agent-runtime/internal/policy"
)

type Summary struct {
	ID                 string   `json:"id"`
	Subjects           []string `json:"subjects"`
	AllowedTools       []string `json:"allowed_tools"`
	WorkspacePatterns  []string `json:"workspace_patterns"`
	CredentialProfiles []string `json:"credential_profiles"`
	AllowTerminal      bool     `json:"allow_terminal"`
	TokenCount         int      `json:"token_count"`
}

type TokenRequest struct {
	Token                     string   `json:"token"`
	SubjectID                 string   `json:"subject"`
	TenantID                  string   `json:"tenant"`
	Role                      string   `json:"role"`
	AllowedTools              []string `json:"allowed_tools"`
	AllowedWorkspaces         []string `json:"allowed_workspaces"`
	AllowedCredentialProfiles []string `json:"allowed_credential_profiles"`
	AllowTerminal             bool     `json:"allow_terminal"`
	MaxJobSeconds             int      `json:"max_job_seconds"`
}

type TokenSummary struct {
	ID                        string   `json:"id"`
	TokenPreview              string   `json:"token_preview"`
	SubjectID                 string   `json:"subject"`
	TenantID                  string   `json:"tenant"`
	Role                      string   `json:"role"`
	AllowedTools              []string `json:"allowed_tools"`
	AllowedWorkspaces         []string `json:"allowed_workspaces"`
	AllowedCredentialProfiles []string `json:"allowed_credential_profiles"`
	AllowTerminal             bool     `json:"allow_terminal"`
	MaxJobSeconds             int      `json:"max_job_seconds"`
}

type Store struct {
	mu        sync.RWMutex
	policies  map[string]policy.Policy
	storePath string
}

func NewStore(policies map[string]policy.Policy) *Store {
	copied := make(map[string]policy.Policy, len(policies))
	for token, p := range policies {
		if p.Role == "" {
			p.Role = "tenant"
		}
		copied[token] = p
	}
	return &Store{policies: copied}
}

func NewPersistentStore(policies map[string]policy.Policy, storePath string) (*Store, error) {
	store := NewStore(policies)
	store.storePath = storePath
	if storePath == "" {
		return store, nil
	}

	raw, err := os.ReadFile(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return nil, fmt.Errorf("read tenant registry: %w", err)
	}
	var persisted persistedStore
	if err := json.Unmarshal(raw, &persisted); err != nil {
		return nil, fmt.Errorf("parse tenant registry: %w", err)
	}
	for _, item := range persisted.Tokens {
		if err := store.upsertLocked(item); err != nil {
			return nil, err
		}
	}
	return store, nil
}

func (s *Store) Lookup(token string) (policy.Policy, bool) {
	if s == nil {
		return policy.Policy{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.policies[token]
	return p, ok
}

func (s *Store) List() []Summary {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return summariesFor(s.policies, nil)
}

func (s *Store) ListFor(actor policy.Policy) []Summary {
	if s == nil {
		return nil
	}
	if !actor.IsAdmin() {
		return summariesFor(map[string]policy.Policy{"current": actor}, nil)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return summariesFor(s.policies, &actor)
}

func (s *Store) ListTokens() []TokenSummary {
	if s == nil {
		return nil
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]TokenSummary, 0, len(s.policies))
	for token, p := range s.policies {
		out = append(out, tokenSummary(token, p))
	}
	sort.Slice(out, func(i int, j int) bool {
		if out[i].TenantID == out[j].TenantID {
			return out[i].SubjectID < out[j].SubjectID
		}
		return out[i].TenantID < out[j].TenantID
	})
	return out
}

func (s *Store) UpsertToken(req TokenRequest) (TokenSummary, error) {
	if s == nil {
		return TokenSummary{}, fmt.Errorf("tenant store is not configured")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.upsertLocked(req); err != nil {
		return TokenSummary{}, err
	}
	if err := s.saveLocked(); err != nil {
		return TokenSummary{}, err
	}
	return tokenSummary(req.Token, s.policies[req.Token]), nil
}

func (s *Store) DeleteToken(id string) error {
	if s == nil {
		return fmt.Errorf("tenant store is not configured")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for token := range s.policies {
		if tokenID(token) == id {
			delete(s.policies, token)
			return s.saveLocked()
		}
	}
	return fmt.Errorf("token not found")
}

func (s *Store) upsertLocked(req TokenRequest) error {
	req.Token = strings.TrimSpace(req.Token)
	req.SubjectID = strings.TrimSpace(req.SubjectID)
	req.TenantID = strings.TrimSpace(req.TenantID)
	req.Role = strings.TrimSpace(req.Role)
	if req.Role == "" {
		req.Role = "tenant"
	}
	if req.Role != "admin" && req.Role != "tenant" {
		return fmt.Errorf("role must be admin or tenant")
	}
	if req.Token == "" {
		return fmt.Errorf("token is required")
	}
	if req.SubjectID == "" {
		return fmt.Errorf("subject is required")
	}
	if !safeID(req.TenantID) {
		return fmt.Errorf("tenant must be a safe id")
	}
	if req.MaxJobSeconds < 0 {
		return fmt.Errorf("max_job_seconds must be positive")
	}
	s.policies[req.Token] = policy.Policy{
		SubjectID:                 req.SubjectID,
		TenantID:                  req.TenantID,
		Role:                      req.Role,
		AllowedTools:              cleanList(req.AllowedTools),
		AllowedWorkspaces:         cleanList(req.AllowedWorkspaces),
		AllowedCredentialProfiles: cleanList(req.AllowedCredentialProfiles),
		AllowTerminal:             req.AllowTerminal,
		MaxJobDuration:            time.Duration(req.MaxJobSeconds) * time.Second,
	}
	return nil
}

func (s *Store) saveLocked() error {
	if s.storePath == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(s.storePath), 0o700); err != nil {
		return fmt.Errorf("create tenant registry directory: %w", err)
	}
	tokens := make([]TokenRequest, 0, len(s.policies))
	for token, p := range s.policies {
		tokens = append(tokens, tokenRequest(token, p))
	}
	sort.Slice(tokens, func(i int, j int) bool {
		if tokens[i].TenantID == tokens[j].TenantID {
			return tokens[i].SubjectID < tokens[j].SubjectID
		}
		return tokens[i].TenantID < tokens[j].TenantID
	})
	raw, err := json.MarshalIndent(persistedStore{Tokens: tokens}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode tenant registry: %w", err)
	}
	tmpPath := s.storePath + ".tmp"
	if err := os.WriteFile(tmpPath, append(raw, '\n'), 0o600); err != nil {
		return fmt.Errorf("write tenant registry: %w", err)
	}
	if err := os.Rename(tmpPath, s.storePath); err != nil {
		return fmt.Errorf("replace tenant registry: %w", err)
	}
	return nil
}

func summariesFor(policies map[string]policy.Policy, actor *policy.Policy) []Summary {
	byTenant := make(map[string]*Summary)
	for _, p := range policies {
		if p.TenantID == "" {
			continue
		}
		if actor != nil && !actor.IsAdmin() && p.TenantID != actor.TenantID {
			continue
		}
		summary := byTenant[p.TenantID]
		if summary == nil {
			summary = &Summary{ID: p.TenantID}
		}
		byTenant[p.TenantID] = summary
		summary.Subjects = appendUnique(summary.Subjects, p.SubjectID)
		summary.AllowedTools = appendUniqueAll(summary.AllowedTools, p.AllowedTools)
		summary.WorkspacePatterns = appendUniqueAll(summary.WorkspacePatterns, p.AllowedWorkspaces)
		summary.CredentialProfiles = appendUniqueAll(summary.CredentialProfiles, p.AllowedCredentialProfiles)
		summary.AllowTerminal = summary.AllowTerminal || p.AllowTerminal
		summary.TokenCount++
	}

	out := make([]Summary, 0, len(byTenant))
	for _, summary := range byTenant {
		sort.Strings(summary.Subjects)
		sort.Strings(summary.AllowedTools)
		sort.Strings(summary.WorkspacePatterns)
		sort.Strings(summary.CredentialProfiles)
		out = append(out, *summary)
	}
	sort.Slice(out, func(i int, j int) bool {
		return out[i].ID < out[j].ID
	})
	return out
}

func tokenRequest(token string, p policy.Policy) TokenRequest {
	return TokenRequest{
		Token:                     token,
		SubjectID:                 p.SubjectID,
		TenantID:                  p.TenantID,
		Role:                      roleOrDefault(p.Role),
		AllowedTools:              append([]string(nil), p.AllowedTools...),
		AllowedWorkspaces:         append([]string(nil), p.AllowedWorkspaces...),
		AllowedCredentialProfiles: append([]string(nil), p.AllowedCredentialProfiles...),
		AllowTerminal:             p.AllowTerminal,
		MaxJobSeconds:             int(p.MaxJobDuration / time.Second),
	}
}

func tokenSummary(token string, p policy.Policy) TokenSummary {
	return TokenSummary{
		ID:                        tokenID(token),
		TokenPreview:              previewToken(token),
		SubjectID:                 p.SubjectID,
		TenantID:                  p.TenantID,
		Role:                      roleOrDefault(p.Role),
		AllowedTools:              append([]string(nil), p.AllowedTools...),
		AllowedWorkspaces:         append([]string(nil), p.AllowedWorkspaces...),
		AllowedCredentialProfiles: append([]string(nil), p.AllowedCredentialProfiles...),
		AllowTerminal:             p.AllowTerminal,
		MaxJobSeconds:             int(p.MaxJobDuration / time.Second),
	}
}

func tokenID(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])[:16]
}

func previewToken(token string) string {
	if len(token) <= 6 {
		return "******"
	}
	return token[:3] + "..." + token[len(token)-3:]
}

func roleOrDefault(role string) string {
	if role == "" {
		return "tenant"
	}
	return role
}

func cleanList(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			out = appendUnique(out, value)
		}
	}
	return out
}

func safeID(id string) bool {
	id = strings.TrimSpace(id)
	return id != "" && id != "." && id != ".." && !strings.ContainsAny(id, `/\`)
}

func appendUniqueAll(values []string, additions []string) []string {
	for _, addition := range additions {
		values = appendUnique(values, addition)
	}
	return values
}

func appendUnique(values []string, value string) []string {
	if value == "" {
		return values
	}
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}

type persistedStore struct {
	Tokens []TokenRequest `json:"tokens"`
}
