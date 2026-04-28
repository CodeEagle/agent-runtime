package jobs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"agent-runtime/internal/policy"
	"agent-runtime/internal/tools"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusTimedOut  Status = "timed_out"
	StatusCanceled  Status = "canceled"
)

type EventType string

const (
	EventStatus EventType = "status"
	EventStdout EventType = "stdout"
	EventStderr EventType = "stderr"
	EventExit   EventType = "exit"
)

type Event struct {
	Type    EventType `json:"type"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

type EventSink interface {
	Emit(Event)
}

type ExecutionSpec struct {
	JobID             string
	TenantID          string
	SubjectID         string
	Tool              tools.Tool
	Executable        string
	Args              []string
	WorkingDir        string
	CredentialProfile string
	Env               map[string]string
	Timeout           time.Duration
}

type ExecutionResult struct {
	ExitCode int
	Error    string
}

type Executor interface {
	Run(context.Context, ExecutionSpec, EventSink) ExecutionResult
}

type PolicyStore interface {
	Lookup(token string) (policy.Policy, bool)
}

type StaticPolicyStore map[string]policy.Policy

func (s StaticPolicyStore) Lookup(token string) (policy.Policy, bool) {
	p, ok := s[token]
	return p, ok
}

type ToolResolver interface {
	Resolve(name string) (tools.Tool, bool)
	List() []tools.Tool
}

type Options struct {
	Policies                 PolicyStore
	Tools                    ToolResolver
	ResolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	ResolveCredentialProfile func(tenantID string, profileID string) (string, error)
	Executor                 Executor
}

type CreateRequest struct {
	Token             string            `json:"-"`
	TenantID          string            `json:"tenant"`
	Tool              string            `json:"tool"`
	Args              []string          `json:"args"`
	WorkspaceID       string            `json:"workspace"`
	CredentialProfile string            `json:"credential_profile"`
	Env               map[string]string `json:"env"`
	Timeout           time.Duration     `json:"-"`
	TimeoutSeconds    int               `json:"timeout_seconds,omitempty"`
}

type Job struct {
	ID                string    `json:"id"`
	TenantID          string    `json:"tenant"`
	SubjectID         string    `json:"subject"`
	Tool              string    `json:"tool"`
	ToolVersion       string    `json:"tool_version"`
	Args              []string  `json:"args"`
	WorkspaceID       string    `json:"workspace"`
	WorkingDir        string    `json:"cwd"`
	CredentialProfile string    `json:"credential_profile"`
	Status            Status    `json:"status"`
	ExitCode          int       `json:"exit_code"`
	Error             string    `json:"error,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	StartedAt         time.Time `json:"started_at,omitempty"`
	FinishedAt        time.Time `json:"finished_at,omitempty"`
}

type Manager struct {
	policies                 PolicyStore
	tools                    ToolResolver
	resolveWorkspace         func(tenantID string, workspaceID string) (string, error)
	resolveCredentialProfile func(tenantID string, profileID string) (string, error)
	executor                 Executor

	mu     sync.RWMutex
	jobs   map[string]Job
	events map[string][]Event
	subs   map[string]map[chan Event]struct{}
}

func NewManager(options Options) *Manager {
	return &Manager{
		policies:                 options.Policies,
		tools:                    options.Tools,
		resolveWorkspace:         options.ResolveWorkspace,
		resolveCredentialProfile: options.ResolveCredentialProfile,
		executor:                 options.Executor,
		jobs:                     make(map[string]Job),
		events:                   make(map[string][]Event),
		subs:                     make(map[string]map[chan Event]struct{}),
	}
}

func (m *Manager) Create(ctx context.Context, req CreateRequest) (Job, error) {
	if req.Timeout == 0 && req.TimeoutSeconds > 0 {
		req.Timeout = time.Duration(req.TimeoutSeconds) * time.Second
	}
	if req.Timeout == 0 {
		req.Timeout = 10 * time.Minute
	}

	if m.policies == nil {
		return Job{}, fmt.Errorf("policy store is not configured")
	}
	p, ok := m.policies.Lookup(req.Token)
	if !ok {
		return Job{}, fmt.Errorf("token is not authorized")
	}
	if err := p.AuthorizeJob(policy.JobRequest{
		TenantID:          req.TenantID,
		Tool:              req.Tool,
		WorkspaceID:       req.WorkspaceID,
		CredentialProfile: req.CredentialProfile,
		RequestedDuration: req.Timeout,
	}); err != nil {
		return Job{}, err
	}

	if m.tools == nil {
		return Job{}, fmt.Errorf("tool registry is not configured")
	}
	tool, ok := m.tools.Resolve(req.Tool)
	if !ok {
		return Job{}, fmt.Errorf("tool %q is not registered", req.Tool)
	}
	if tool.Path == "" {
		return Job{}, fmt.Errorf("tool %q has no executable path", req.Tool)
	}
	if m.resolveWorkspace == nil {
		return Job{}, fmt.Errorf("workspace resolver is not configured")
	}
	workingDir, err := m.resolveWorkspace(req.TenantID, req.WorkspaceID)
	if err != nil {
		return Job{}, fmt.Errorf("resolve workspace: %w", err)
	}
	if m.resolveCredentialProfile == nil {
		return Job{}, fmt.Errorf("credential resolver is not configured")
	}
	credentialRoot, err := m.resolveCredentialProfile(req.TenantID, req.CredentialProfile)
	if err != nil {
		return Job{}, fmt.Errorf("resolve credential profile: %w", err)
	}
	if m.executor == nil {
		return Job{}, fmt.Errorf("executor is not configured")
	}

	env := copyEnv(req.Env)
	env["HOME"] = credentialRoot
	env["XDG_CONFIG_HOME"] = filepath.Join(credentialRoot, ".config")
	env["NPM_CONFIG_PREFIX"] = "/data/npm-global"
	env["PATH"] = tools.RuntimePath(credentialRoot, os.Getenv("PATH"))
	if tool.CredentialEnv != "" {
		credentialPath := credentialRoot
		if tool.CredentialSubdir != "" {
			credentialPath = filepath.Join(credentialRoot, tool.CredentialSubdir)
		}
		env[tool.CredentialEnv] = credentialPath
	}
	env["AGENT_RUNTIME_TENANT"] = req.TenantID
	env["AGENT_RUNTIME_SUBJECT"] = p.SubjectID

	job := Job{
		ID:                newID(),
		TenantID:          req.TenantID,
		SubjectID:         p.SubjectID,
		Tool:              tool.Name,
		ToolVersion:       tool.Version,
		Args:              append([]string(nil), req.Args...),
		WorkspaceID:       req.WorkspaceID,
		WorkingDir:        workingDir,
		CredentialProfile: req.CredentialProfile,
		Status:            StatusQueued,
		CreatedAt:         time.Now(),
	}

	m.mu.Lock()
	m.jobs[job.ID] = job
	m.events[job.ID] = []Event{{Type: EventStatus, Message: string(StatusQueued), At: job.CreatedAt}}
	m.mu.Unlock()

	spec := ExecutionSpec{
		JobID:             job.ID,
		TenantID:          req.TenantID,
		SubjectID:         p.SubjectID,
		Tool:              tool,
		Executable:        tool.Path,
		Args:              append([]string(nil), req.Args...),
		WorkingDir:        workingDir,
		CredentialProfile: req.CredentialProfile,
		Env:               env,
		Timeout:           req.Timeout,
	}
	go m.run(spec)

	return job, nil
}

func (m *Manager) Get(id string) (Job, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	job, ok := m.jobs[id]
	return job, ok
}

func (m *Manager) Events(id string) []Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	events := m.events[id]
	return append([]Event(nil), events...)
}

func (m *Manager) Subscribe(id string) ([]Event, <-chan Event, func(), bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.jobs[id]; !ok {
		return nil, nil, func() {}, false
	}
	ch := make(chan Event, 256)
	if m.subs[id] == nil {
		m.subs[id] = make(map[chan Event]struct{})
	}
	m.subs[id][ch] = struct{}{}
	events := append([]Event(nil), m.events[id]...)
	unsubscribe := func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		if subs := m.subs[id]; subs != nil {
			delete(subs, ch)
			if len(subs) == 0 {
				delete(m.subs, id)
			}
		}
	}
	return events, ch, unsubscribe, true
}

func (m *Manager) run(spec ExecutionSpec) {
	started := time.Now()
	m.update(spec.JobID, func(job *Job) {
		job.Status = StatusRunning
		job.StartedAt = started
	})
	m.appendEvent(spec.JobID, Event{Type: EventStatus, Message: string(StatusRunning), At: started})

	ctx, cancel := context.WithTimeout(context.Background(), spec.Timeout)
	defer cancel()

	result := m.executor.Run(ctx, spec, managerSink{manager: m, jobID: spec.JobID})
	finished := time.Now()

	status := StatusSucceeded
	if ctx.Err() == context.DeadlineExceeded {
		status = StatusTimedOut
	} else if result.Error != "" || result.ExitCode != 0 {
		status = StatusFailed
	}

	m.update(spec.JobID, func(job *Job) {
		job.Status = status
		job.ExitCode = result.ExitCode
		job.Error = result.Error
		job.FinishedAt = finished
	})
	m.appendEvent(spec.JobID, Event{Type: EventExit, Message: string(status), At: finished})
}

func (m *Manager) update(id string, update func(*Job)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	job := m.jobs[id]
	update(&job)
	m.jobs[id] = job
}

func (m *Manager) appendEvent(id string, event Event) {
	m.mu.Lock()
	if event.At.IsZero() {
		event.At = time.Now()
	}
	m.events[id] = append(m.events[id], event)
	subs := make([]chan Event, 0, len(m.subs[id]))
	for ch := range m.subs[id] {
		subs = append(subs, ch)
	}
	m.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}
}

type managerSink struct {
	manager *Manager
	jobID   string
}

func (s managerSink) Emit(event Event) {
	s.manager.appendEvent(s.jobID, event)
}

func copyEnv(in map[string]string) map[string]string {
	out := make(map[string]string, len(in)+2)
	for key, value := range in {
		out[key] = value
	}
	return out
}

func newID() string {
	var bytes [8]byte
	if _, err := rand.Read(bytes[:]); err == nil {
		return hex.EncodeToString(bytes[:])
	}
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
