package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Resolver struct {
	tenantsDir string
}

func NewResolver(tenantsDir string) Resolver {
	return Resolver{tenantsDir: tenantsDir}
}

func (r Resolver) ResolveProfile(tenantID string, profileID string) (string, error) {
	if !safeID(tenantID) {
		return "", fmt.Errorf("unsafe tenant id %q", tenantID)
	}
	if !safeID(profileID) {
		return "", fmt.Errorf("unsafe credential profile id %q", profileID)
	}
	root, err := filepath.Abs(r.tenantsDir)
	if err != nil {
		return "", fmt.Errorf("resolve tenants dir: %w", err)
	}
	path := filepath.Join(root, tenantID, "homes", profileID)
	if err := os.MkdirAll(path, 0o700); err != nil {
		return "", fmt.Errorf("create credential profile: %w", err)
	}
	return path, nil
}

func safeID(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	return !strings.ContainsRune(value, filepath.Separator) && !strings.Contains(value, "/") && !strings.Contains(value, "\\")
}
