package tools

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Health struct {
	Available       bool   `json:"available"`
	Health          string `json:"health"`
	ResolvedPath    string `json:"resolved_path,omitempty"`
	DetectedVersion string `json:"detected_version,omitempty"`
	Error           string `json:"error,omitempty"`
}

func CheckHealth(tool Tool, searchPath string) Health {
	resolved, err := resolveExecutable(tool.Path, searchPath)
	if err != nil {
		return Health{Available: false, Health: "missing", Error: err.Error()}
	}
	health := Health{Available: true, Health: "ok", ResolvedPath: resolved}
	if version, err := detectVersion(resolved, searchPath); err == nil {
		health.DetectedVersion = version
	}
	return health
}

func resolveExecutable(command string, searchPath string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("empty command")
	}
	if strings.Contains(command, "/") {
		if isExecutable(command) {
			return command, nil
		}
		return "", fmt.Errorf("%s is not executable", command)
	}
	for _, dir := range filepath.SplitList(searchPath) {
		if dir == "" {
			continue
		}
		candidate := filepath.Join(dir, command)
		if isExecutable(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("%s not found in PATH", command)
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	return info.Mode()&0o111 != 0
}

func detectVersion(path string, searchPath string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1500*time.Millisecond)
	defer cancel()
	cmd := exec.CommandContext(ctx, path, "--version")
	cmd.Env = append(os.Environ(), "PATH="+searchPath)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		return "", err
	}
	version := strings.TrimSpace(output.String())
	if version == "" {
		return "", fmt.Errorf("empty version")
	}
	lines := strings.Split(version, "\n")
	return strings.TrimSpace(lines[0]), nil
}
