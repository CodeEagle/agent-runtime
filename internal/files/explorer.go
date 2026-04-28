package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Explorer struct {
	tenantsDir string
}

type Listing struct {
	TenantID string  `json:"tenant"`
	Space    string  `json:"space"`
	Path     string  `json:"path"`
	AbsPath  string  `json:"abs_path"`
	Entries  []Entry `json:"entries"`
}

type Entry struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Kind     string    `json:"kind"`
	Size     int64     `json:"size"`
	Modified time.Time `json:"modified"`
}

type Content struct {
	TenantID  string    `json:"tenant"`
	Space     string    `json:"space"`
	Path      string    `json:"path"`
	AbsPath   string    `json:"abs_path"`
	Size      int64     `json:"size"`
	Modified  time.Time `json:"modified"`
	Content   string    `json:"content"`
	Truncated bool      `json:"truncated"`
}

func NewExplorer(tenantsDir string) Explorer {
	return Explorer{tenantsDir: tenantsDir}
}

func (e Explorer) Configured() bool {
	return strings.TrimSpace(e.tenantsDir) != ""
}

func (e Explorer) EnsureTenant(tenantID string, credentialProfiles []string, workspacePatterns []string) error {
	if !safeID(tenantID) {
		return fmt.Errorf("unsafe tenant id %q", tenantID)
	}
	if !e.Configured() {
		return fmt.Errorf("tenants directory is not configured")
	}
	root, err := filepath.Abs(e.tenantsDir)
	if err != nil {
		return fmt.Errorf("resolve tenants dir: %w", err)
	}
	tenantRoot := filepath.Join(root, tenantID)
	homesRoot := filepath.Join(tenantRoot, "homes")
	workspacesRoot := filepath.Join(tenantRoot, "workspaces")
	if err := os.MkdirAll(homesRoot, 0o755); err != nil {
		return fmt.Errorf("create tenant homes root: %w", err)
	}
	if err := os.MkdirAll(workspacesRoot, 0o755); err != nil {
		return fmt.Errorf("create tenant workspaces root: %w", err)
	}
	for _, profileID := range concreteIDs(credentialProfiles, "team-default") {
		if err := os.MkdirAll(filepath.Join(homesRoot, profileID), 0o700); err != nil {
			return fmt.Errorf("create credential profile %q: %w", profileID, err)
		}
	}
	for _, workspaceID := range workspaceIDs(workspacePatterns) {
		if err := os.MkdirAll(filepath.Join(workspacesRoot, workspaceID), 0o755); err != nil {
			return fmt.Errorf("create workspace %q: %w", workspaceID, err)
		}
	}
	return nil
}

func (e Explorer) List(tenantID string, space string, requestedPath string) (Listing, error) {
	if !safeID(tenantID) {
		return Listing{}, fmt.Errorf("unsafe tenant id %q", tenantID)
	}
	space = normalizeSpace(space)
	if space == "" {
		return Listing{}, fmt.Errorf("space must be workspaces or homes")
	}
	root, err := filepath.Abs(e.tenantsDir)
	if err != nil {
		return Listing{}, fmt.Errorf("resolve tenants dir: %w", err)
	}
	base := filepath.Join(root, tenantID, space)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return Listing{}, fmt.Errorf("create tenant %s root: %w", space, err)
	}

	relPath, err := cleanRelativePath(requestedPath)
	if err != nil {
		return Listing{}, err
	}
	target := filepath.Join(base, relPath)
	resolved, err := filepath.Abs(target)
	if err != nil {
		return Listing{}, fmt.Errorf("resolve path: %w", err)
	}
	if resolved != base && !strings.HasPrefix(resolved, base+string(os.PathSeparator)) {
		return Listing{}, fmt.Errorf("path escapes tenant root")
	}

	info, err := os.Stat(resolved)
	if err != nil {
		return Listing{}, fmt.Errorf("stat path: %w", err)
	}
	if !info.IsDir() {
		return Listing{}, fmt.Errorf("path is not a directory")
	}

	items, err := os.ReadDir(resolved)
	if err != nil {
		return Listing{}, fmt.Errorf("read directory: %w", err)
	}
	entries := make([]Entry, 0, len(items))
	for _, item := range items {
		itemInfo, err := item.Info()
		if err != nil {
			continue
		}
		kind := "file"
		if itemInfo.IsDir() {
			kind = "directory"
		}
		entryRel := filepath.ToSlash(filepath.Join(relPath, item.Name()))
		if entryRel == "." {
			entryRel = ""
		}
		entries = append(entries, Entry{
			Name:     item.Name(),
			Path:     entryRel,
			Kind:     kind,
			Size:     itemInfo.Size(),
			Modified: itemInfo.ModTime(),
		})
	}
	sort.Slice(entries, func(i int, j int) bool {
		if entries[i].Kind != entries[j].Kind {
			return entries[i].Kind == "directory"
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return Listing{
		TenantID: tenantID,
		Space:    space,
		Path:     filepath.ToSlash(relPath),
		AbsPath:  resolved,
		Entries:  entries,
	}, nil
}

func (e Explorer) Read(tenantID string, space string, requestedPath string, maxBytes int64) (Content, error) {
	if maxBytes <= 0 {
		maxBytes = 2 * 1024 * 1024
	}
	_, relPath, resolved, err := e.resolve(tenantID, space, requestedPath)
	if err != nil {
		return Content{}, err
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return Content{}, fmt.Errorf("stat path: %w", err)
	}
	if info.IsDir() {
		return Content{}, fmt.Errorf("path is a directory")
	}
	file, err := os.Open(resolved)
	if err != nil {
		return Content{}, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	limited := io.LimitReader(file, maxBytes+1)
	raw, err := io.ReadAll(limited)
	if err != nil {
		return Content{}, fmt.Errorf("read file: %w", err)
	}
	truncated := int64(len(raw)) > maxBytes
	if truncated {
		raw = raw[:maxBytes]
	}
	return Content{
		TenantID:  tenantID,
		Space:     normalizeSpace(space),
		Path:      filepath.ToSlash(relPath),
		AbsPath:   resolved,
		Size:      info.Size(),
		Modified:  info.ModTime(),
		Content:   string(raw),
		Truncated: truncated || info.Size() > int64(len(raw)),
	}, nil
}

func (e Explorer) resolve(tenantID string, space string, requestedPath string) (string, string, string, error) {
	if !safeID(tenantID) {
		return "", "", "", fmt.Errorf("unsafe tenant id %q", tenantID)
	}
	space = normalizeSpace(space)
	if space == "" {
		return "", "", "", fmt.Errorf("space must be workspaces or homes")
	}
	root, err := filepath.Abs(e.tenantsDir)
	if err != nil {
		return "", "", "", fmt.Errorf("resolve tenants dir: %w", err)
	}
	base := filepath.Join(root, tenantID, space)
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", "", "", fmt.Errorf("create tenant %s root: %w", space, err)
	}

	relPath, err := cleanRelativePath(requestedPath)
	if err != nil {
		return "", "", "", err
	}
	target := filepath.Join(base, relPath)
	resolved, err := filepath.Abs(target)
	if err != nil {
		return "", "", "", fmt.Errorf("resolve path: %w", err)
	}
	if resolved != base && !strings.HasPrefix(resolved, base+string(os.PathSeparator)) {
		return "", "", "", fmt.Errorf("path escapes tenant root")
	}
	return base, relPath, resolved, nil
}

func normalizeSpace(space string) string {
	switch strings.TrimSpace(space) {
	case "", "workspaces", "workspace":
		return "workspaces"
	case "homes", "home", "credentials", "credential":
		return "homes"
	default:
		return ""
	}
}

func cleanRelativePath(path string) (string, error) {
	path = strings.TrimSpace(strings.TrimPrefix(path, "/"))
	if path == "" {
		return ".", nil
	}
	cleaned := filepath.Clean(filepath.FromSlash(path))
	if cleaned == "." {
		return ".", nil
	}
	if cleaned == ".." || strings.HasPrefix(cleaned, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes tenant root")
	}
	return cleaned, nil
}

func safeID(id string) bool {
	id = strings.TrimSpace(id)
	return id != "" && id != "." && id != ".." && !strings.ContainsAny(id, `/\`)
}

func concreteIDs(values []string, fallback string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if concreteID(value) {
			out = appendUnique(out, value)
		}
	}
	if len(out) == 0 {
		out = append(out, fallback)
	}
	return out
}

func workspaceIDs(patterns []string) []string {
	out := make([]string, 0, len(patterns))
	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if concreteID(pattern) {
			out = appendUnique(out, pattern)
			continue
		}
		if strings.HasSuffix(pattern, "*") && !strings.ContainsAny(strings.TrimSuffix(pattern, "*"), "*?[]") {
			candidate := strings.TrimSuffix(pattern, "*") + "main"
			if concreteID(candidate) {
				out = appendUnique(out, candidate)
			}
		}
	}
	if len(out) == 0 {
		out = append(out, "repo-main")
	}
	return out
}

func concreteID(id string) bool {
	return safeID(id) && !strings.ContainsAny(id, "*?[]")
}

func appendUnique(values []string, value string) []string {
	for _, existing := range values {
		if existing == value {
			return values
		}
	}
	return append(values, value)
}
