package files

import (
	"fmt"
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

func NewExplorer(tenantsDir string) Explorer {
	return Explorer{tenantsDir: tenantsDir}
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
