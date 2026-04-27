package tools

import (
	"path/filepath"
	"strings"
)

func RuntimePath(home string, base string) string {
	parts := []string{}
	if home != "" {
		parts = append(parts,
			filepath.Join(home, ".local", "bin"),
			filepath.Join(home, "bin"),
			filepath.Join(home, ".npm-global", "bin"),
			filepath.Join(home, ".bun", "bin"),
		)
	}
	parts = append(parts,
		"/data/bin",
		"/data/npm-global/bin",
	)
	if base != "" {
		parts = append(parts, base)
	}
	return strings.Join(parts, ":")
}
