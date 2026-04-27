package tools

import (
	"sort"
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
	mu     sync.RWMutex
	byName map[string]Tool
}

func NewRegistry(initial []Tool) *Registry {
	registry := &Registry{byName: make(map[string]Tool)}
	for _, tool := range initial {
		if tool.Name == "" {
			continue
		}
		registry.byName[tool.Name] = tool
	}
	return registry
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

	out := make([]Tool, 0, len(r.byName))
	for _, tool := range r.byName {
		out = append(out, tool)
	}
	sort.Slice(out, func(i int, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}
