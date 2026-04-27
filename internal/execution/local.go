package execution

import (
	"bufio"
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
	"time"

	"agent-runtime/internal/jobs"
)

type LocalExecutor struct{}

func (LocalExecutor) Run(ctx context.Context, spec jobs.ExecutionSpec, sink jobs.EventSink) jobs.ExecutionResult {
	cmd := exec.CommandContext(ctx, spec.Executable, spec.Args...)
	cmd.Dir = spec.WorkingDir
	cmd.Env = mergeEnv(os.Environ(), spec.Env)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.WaitDelay = 5 * time.Second
	cmd.Cancel = func() error {
		if cmd.Process == nil {
			return nil
		}
		return syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return jobs.ExecutionResult{ExitCode: -1, Error: err.Error()}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return jobs.ExecutionResult{ExitCode: -1, Error: err.Error()}
	}
	if err := cmd.Start(); err != nil {
		return jobs.ExecutionResult{ExitCode: -1, Error: err.Error()}
	}

	done := make(chan struct{}, 2)
	go scanPipe(stdout, jobs.EventStdout, sink, done)
	go scanPipe(stderr, jobs.EventStderr, sink, done)

	err = cmd.Wait()
	<-done
	<-done

	if err == nil {
		return jobs.ExecutionResult{ExitCode: 0}
	}
	if ctx.Err() != nil {
		return jobs.ExecutionResult{ExitCode: -1, Error: ctx.Err().Error()}
	}

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return jobs.ExecutionResult{ExitCode: exitError.ExitCode(), Error: err.Error()}
	}
	return jobs.ExecutionResult{ExitCode: -1, Error: err.Error()}
}

func scanPipe(pipe any, eventType jobs.EventType, sink jobs.EventSink, done chan<- struct{}) {
	defer func() { done <- struct{}{} }()

	reader, ok := pipe.(interface {
		Read([]byte) (int, error)
	})
	if !ok {
		return
	}
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		sink.Emit(jobs.Event{Type: eventType, Message: scanner.Text() + "\n", At: time.Now()})
	}
	if err := scanner.Err(); err != nil {
		sink.Emit(jobs.Event{Type: jobs.EventStderr, Message: err.Error(), At: time.Now()})
	}
}

func mergeEnv(base []string, extra map[string]string) []string {
	out := make([]string, 0, len(base)+len(extra))
	seen := make(map[string]int, len(base)+len(extra))
	for _, item := range base {
		key := envKey(item)
		seen[key] = len(out)
		out = append(out, item)
	}
	for key, value := range extra {
		item := key + "=" + value
		if index, ok := seen[key]; ok {
			out[index] = item
			continue
		}
		seen[key] = len(out)
		out = append(out, item)
	}
	return out
}

func envKey(item string) string {
	for i, char := range item {
		if char == '=' {
			return item[:i]
		}
	}
	return item
}
