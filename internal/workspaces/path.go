package workspaces

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

func (r Resolver) ResolveWorkspace(tenantID string, workspaceID string) (string, error) {
	if !safeID(tenantID) {
		return "", fmt.Errorf("unsafe tenant id %q", tenantID)
	}
	if !safeID(workspaceID) {
		return "", fmt.Errorf("unsafe workspace id %q", workspaceID)
	}
	root, err := filepath.Abs(r.tenantsDir)
	if err != nil {
		return "", fmt.Errorf("resolve tenants dir: %w", err)
	}
	parent := filepath.Join(root, tenantID, "workspaces")
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return "", fmt.Errorf("create workspace parent: %w", err)
	}
	resolved, err := ResolveWithin(root, filepath.Join(tenantID, "workspaces", workspaceID))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(resolved, 0o755); err != nil {
		return "", fmt.Errorf("create workspace: %w", err)
	}
	return resolved, nil
}

func ResolveWithin(root string, requested string) (string, error) {
	if root == "" {
		return "", fmt.Errorf("workspace root is required")
	}

	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	realRoot, err := filepath.EvalSymlinks(absoluteRoot)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root symlinks: %w", err)
	}

	candidate := filepath.Clean(filepath.Join(realRoot, requested))
	if !isWithin(realRoot, candidate) {
		return "", fmt.Errorf("path %q escapes workspace root %q", requested, realRoot)
	}

	realCandidate, err := filepath.EvalSymlinks(candidate)
	if err != nil {
		if os.IsNotExist(err) {
			parent := filepath.Dir(candidate)
			realParent, parentErr := filepath.EvalSymlinks(parent)
			if parentErr != nil {
				return "", fmt.Errorf("resolve parent symlinks for %q: %w", requested, parentErr)
			}
			if !isWithin(realRoot, realParent) {
				return "", fmt.Errorf("path %q escapes workspace root %q", requested, realRoot)
			}
			return candidate, nil
		}
		return "", fmt.Errorf("resolve requested path symlinks: %w", err)
	}
	if !isWithin(realRoot, realCandidate) {
		return "", fmt.Errorf("path %q escapes workspace root %q", requested, realRoot)
	}
	return realCandidate, nil
}

func isWithin(root string, candidate string) bool {
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return false
	}
	return rel == "." || (rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)))
}

func safeID(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	return !strings.ContainsRune(value, filepath.Separator) && !strings.Contains(value, "/") && !strings.Contains(value, "\\")
}
