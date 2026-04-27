package tenants

import (
	"sort"

	"agent-runtime/internal/policy"
)

type Summary struct {
	ID                 string   `json:"id"`
	Subjects           []string `json:"subjects"`
	AllowedTools       []string `json:"allowed_tools"`
	WorkspacePatterns  []string `json:"workspace_patterns"`
	CredentialProfiles []string `json:"credential_profiles"`
	AllowTerminal      bool     `json:"allow_terminal"`
}

type Store struct {
	policies map[string]policy.Policy
}

func NewStore(policies map[string]policy.Policy) Store {
	copied := make(map[string]policy.Policy, len(policies))
	for token, p := range policies {
		copied[token] = p
	}
	return Store{policies: copied}
}

func (s Store) Lookup(token string) (policy.Policy, bool) {
	p, ok := s.policies[token]
	return p, ok
}

func (s Store) List() []Summary {
	byTenant := make(map[string]*Summary)
	for _, p := range s.policies {
		if p.TenantID == "" {
			continue
		}
		summary := byTenant[p.TenantID]
		if summary == nil {
			summary = &Summary{ID: p.TenantID}
			byTenant[p.TenantID] = summary
		}
		summary.Subjects = appendUnique(summary.Subjects, p.SubjectID)
		summary.AllowedTools = appendUniqueAll(summary.AllowedTools, p.AllowedTools)
		summary.WorkspacePatterns = appendUniqueAll(summary.WorkspacePatterns, p.AllowedWorkspaces)
		summary.CredentialProfiles = appendUniqueAll(summary.CredentialProfiles, p.AllowedCredentialProfiles)
		summary.AllowTerminal = summary.AllowTerminal || p.AllowTerminal
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
