package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type Tool struct {
	Name             string `json:"name"`
	Path             string `json:"path"`
	Version          string `json:"version"`
	CredentialEnv    string `json:"credential_env,omitempty"`
	CredentialSubdir string `json:"credential_subdir,omitempty"`
}

type Registry struct {
	mu        sync.RWMutex
	byName    map[string]Tool
	storePath string
}

func NewRegistry(initial []Tool) *Registry {
	registry := &Registry{byName: make(map[string]Tool)}
	for _, tool := range initial {
		if err := Validate(tool); err != nil {
			continue
		}
		registry.byName[tool.Name] = tool
	}
	return registry
}

func NewPersistentRegistry(initial []Tool, storePath string) (*Registry, error) {
	registry := NewRegistry(initial)
	registry.storePath = storePath
	if storePath == "" {
		return registry, nil
	}
	raw, err := os.ReadFile(storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return registry, nil
		}
		return nil, fmt.Errorf("read tool registry: %w", err)
	}
	var stored persistedRegistry
	if err := json.Unmarshal(raw, &stored); err != nil {
		return nil, fmt.Errorf("parse tool registry: %w", err)
	}
	defaults := make(map[string]Tool, len(initial))
	for _, tool := range initial {
		if err := Validate(tool); err == nil {
			defaults[tool.Name] = tool
		}
	}
	for _, tool := range stored.Tools {
		if err := Validate(tool); err == nil {
			if fallback, ok := defaults[tool.Name]; ok && isLegacyEnvPlaceholder(tool) && fallback.Path != tool.Path {
				tool = fallback
			}
			registry.byName[tool.Name] = tool
		}
	}
	return registry, nil
}

func (r *Registry) Resolve(name string) (Tool, bool) {
	if r == nil {
		return Tool{}, false
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	tool, ok := r.byName[name]
	return tool, ok
}

func (r *Registry) List() []Tool {
	if r == nil {
		return nil
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.listLocked()
}

func (r *Registry) listLocked() []Tool {
	out := make([]Tool, 0, len(r.byName))
	for _, tool := range r.byName {
		out = append(out, tool)
	}
	sort.Slice(out, func(i int, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func (r *Registry) Upsert(tool Tool) error {
	if r == nil {
		return fmt.Errorf("tool registry is not configured")
	}
	if err := Validate(tool); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byName[tool.Name] = tool
	return r.saveLocked()
}

func (r *Registry) Delete(name string) error {
	if r == nil {
		return fmt.Errorf("tool registry is not configured")
	}
	if !safeName(name) {
		return fmt.Errorf("tool not found")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.byName[name]; !ok {
		return fmt.Errorf("tool not found")
	}
	delete(r.byName, name)
	return r.saveLocked()
}

func (r *Registry) saveLocked() error {
	if r.storePath == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(r.storePath), 0o755); err != nil {
		return fmt.Errorf("create tool registry directory: %w", err)
	}
	raw, err := json.MarshalIndent(persistedRegistry{Tools: r.listLocked()}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode tool registry: %w", err)
	}
	tmpPath := r.storePath + ".tmp"
	if err := os.WriteFile(tmpPath, append(raw, '\n'), 0o644); err != nil {
		return fmt.Errorf("write tool registry: %w", err)
	}
	if err := os.Rename(tmpPath, r.storePath); err != nil {
		return fmt.Errorf("replace tool registry: %w", err)
	}
	return nil
}

func Validate(tool Tool) error {
	if !safeName(tool.Name) {
		return fmt.Errorf("tool name is required and must not contain path separators")
	}
	if strings.TrimSpace(tool.Path) == "" {
		return fmt.Errorf("tool path is required")
	}
	if strings.Contains(tool.CredentialSubdir, "..") || strings.ContainsAny(tool.CredentialSubdir, `/\`) {
		return fmt.Errorf("credential subdir must be a simple directory name")
	}
	return nil
}

func safeName(name string) bool {
	name = strings.TrimSpace(name)
	return name != "" && name != "." && name != ".." && !strings.ContainsAny(name, `/\`)
}

func isLegacyEnvPlaceholder(tool Tool) bool {
	return tool.Path == "/usr/bin/env" && (tool.Version == "" || tool.Version == "placeholder" || tool.Version == "local-env")
}

type persistedRegistry struct {
	Tools []Tool `json:"tools"`
}
