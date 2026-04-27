package execution

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"agent-runtime/internal/jobs"
)

func TestLocalExecutorResolvesExecutableFromSpecPath(t *testing.T) {
	dir := t.TempDir()
	binDir := filepath.Join(dir, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatal(err)
	}
	commandPath := filepath.Join(binDir, "agent-test-cli")
	if err := os.WriteFile(commandPath, []byte("#!/bin/sh\necho from-spec-path\n"), 0o755); err != nil {
		t.Fatal(err)
	}

	sink := recordingSink{}
	result := LocalExecutor{}.Run(context.Background(), jobs.ExecutionSpec{
		Executable: "agent-test-cli",
		WorkingDir: dir,
		Env: map[string]string{
			"PATH": binDir,
		},
	}, &sink)

	if result.ExitCode != 0 || result.Error != "" {
		t.Fatalf("expected success, got %#v", result)
	}
	if len(sink.events) != 1 || sink.events[0].Message != "from-spec-path\n" {
		t.Fatalf("unexpected events: %#v", sink.events)
	}
}

type recordingSink struct {
	events []jobs.Event
}

func (s *recordingSink) Emit(event jobs.Event) {
	s.events = append(s.events, event)
}
