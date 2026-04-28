package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"agent-runtime/internal/files"
	"agent-runtime/internal/jobs"
	"agent-runtime/internal/policy"
	"agent-runtime/internal/tenants"
	"agent-runtime/internal/tools"

	"github.com/gorilla/websocket"
)

type ToolManager interface {
	List() []tools.Tool
	Upsert(tools.Tool) error
	Delete(name string) error
}

type TenantManager interface {
	Lookup(token string) (policy.Policy, bool)
	AuthenticateUser(username string, password string) (string, policy.Policy, bool)
	List() []tenants.Summary
	ListFor(policy.Policy) []tenants.Summary
	ListUsers() []tenants.UserSummary
	RegisterUser(tenants.UserRequest) (string, tenants.UserSummary, error)
	UpsertUser(tenants.UserRequest) (tenants.UserSummary, error)
	DeleteUser(id string) error
	ListTokens() []tenants.TokenSummary
	UpsertToken(tenants.TokenRequest) (tenants.TokenSummary, error)
	DeleteToken(id string) error
}

type Options struct {
	Jobs                     *jobs.Manager
	Tools                    ToolManager
	Tenants                  TenantManager
	Terminal                 http.Handler
	AppServer                http.Handler
	Files                    files.Explorer
	ResolveCredentialProfile func(tenantID string, profileID string) (string, error)
}

func NewServer(options Options) http.Handler {
	server := &server{
		jobs:                     options.Jobs,
		tools:                    options.Tools,
		tenants:                  options.Tenants,
		terminal:                 options.Terminal,
		appServer:                options.AppServer,
		files:                    options.Files,
		resolveCredentialProfile: options.ResolveCredentialProfile,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", server.webUI)
	mux.HandleFunc("GET /docs", server.swaggerUI)
	mux.HandleFunc("GET /api/health", server.health)
	mux.HandleFunc("GET /api/ready", server.ready)
	mux.HandleFunc("GET /api/status", server.status)
	mux.HandleFunc("GET /openapi.json", server.openapi)
	mux.HandleFunc("POST /api/register", server.register)
	mux.HandleFunc("POST /api/login", server.login)
	mux.HandleFunc("GET /api/session", server.session)
	mux.HandleFunc("GET /api/tools", server.listTools)
	mux.HandleFunc("POST /api/tools", server.upsertTool)
	mux.HandleFunc("DELETE /api/tools/{name}", server.deleteTool)
	mux.HandleFunc("GET /api/tenants", server.listTenants)
	mux.HandleFunc("GET /api/users", server.listUsers)
	mux.HandleFunc("POST /api/users", server.upsertUser)
	mux.HandleFunc("DELETE /api/users/{id}", server.deleteUser)
	mux.HandleFunc("GET /api/tokens", server.listTokens)
	mux.HandleFunc("POST /api/tokens", server.upsertToken)
	mux.HandleFunc("DELETE /api/tokens/{id}", server.deleteToken)
	mux.HandleFunc("GET /api/files", server.listFiles)
	mux.HandleFunc("GET /api/files/raw", server.readFile)
	mux.HandleFunc("GET /api/terminal", server.terminalSession)
	mux.HandleFunc("GET /api/v1/terminal/ws", server.terminalSession)
	mux.HandleFunc("GET /api/app-server/{tool}", server.appServerSession)
	mux.HandleFunc("POST /api/jobs", server.createJob)
	mux.HandleFunc("/api/jobs/", server.jobByID)
	mux.Handle("GET /assets/", http.FileServerFS(assetsFS))
	return mux
}

type server struct {
	jobs                     *jobs.Manager
	tools                    ToolManager
	tenants                  TenantManager
	terminal                 http.Handler
	appServer                http.Handler
	files                    files.Explorer
	resolveCredentialProfile func(tenantID string, profileID string) (string, error)
}

func (s *server) webUI(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(webUIHTML))
}

func (s *server) swaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerUIHTML))
}

func (s *server) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *server) ready(w http.ResponseWriter, r *http.Request) {
	if s.jobs == nil || s.tools == nil {
		writeError(w, http.StatusServiceUnavailable, "runtime is not ready")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

func (s *server) status(w http.ResponseWriter, r *http.Request) {
	toolCount := 0
	if s.tools != nil {
		toolCount = len(s.tools.List())
	}
	tenantCount := 0
	userCount := 0
	if s.tenants != nil {
		tenantCount = len(s.tenants.List())
		userCount = len(s.tenants.ListUsers())
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":   "ok",
		"tools":    toolCount,
		"tenants":  tenantCount,
		"users":    userCount,
		"terminal": s.terminal != nil,
	})
}

func (s *server) openapi(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"openapi": "3.1.0",
		"info": map[string]any{
			"title":   "Agent Runtime API",
			"version": "0.1.0",
		},
		"paths": map[string]any{
			"/api/status": map[string]any{
				"get": map[string]any{
					"summary": "Runtime status",
					"responses": map[string]any{
						"200": map[string]any{"description": "Runtime status"},
					},
				},
			},
			"/api/register": map[string]any{
				"post": map[string]any{"summary": "Register a tenant user and receive a bearer token", "responses": map[string]any{"201": map[string]any{"description": "User registered"}}},
			},
			"/api/login": map[string]any{
				"post": map[string]any{"summary": "Login and receive the current user's bearer token", "responses": map[string]any{"200": map[string]any{"description": "Session token"}}},
			},
			"/api/tools": map[string]any{
				"get":  map[string]any{"summary": "List CLI tools", "responses": map[string]any{"200": map[string]any{"description": "CLI tool list"}}},
				"post": map[string]any{"summary": "Register or update CLI tool", "responses": map[string]any{"201": map[string]any{"description": "CLI tool saved"}}},
			},
			"/api/jobs": map[string]any{
				"post": map[string]any{"summary": "Create a CLI job", "responses": map[string]any{"202": map[string]any{"description": "Job accepted"}}},
			},
			"/api/jobs/{id}": map[string]any{
				"get": map[string]any{"summary": "Fetch job status", "responses": map[string]any{"200": map[string]any{"description": "Job status"}}},
			},
			"/api/jobs/{id}/events": map[string]any{
				"get": map[string]any{"summary": "Read job events as Server-Sent Events", "responses": map[string]any{"200": map[string]any{"description": "Job events"}}},
			},
			"/api/jobs/{id}/events/ws": map[string]any{
				"get": map[string]any{"summary": "Read job events as a WebSocket stream", "responses": map[string]any{"101": map[string]any{"description": "WebSocket upgrade"}}},
			},
			"/api/app-server/{tool}": map[string]any{
				"get": map[string]any{"summary": "Open a tool app-server WebSocket proxy", "responses": map[string]any{"101": map[string]any{"description": "WebSocket upgrade"}}},
			},
			"/api/terminal": map[string]any{
				"get": map[string]any{"summary": "Open an interactive PTY WebSocket", "responses": map[string]any{"101": map[string]any{"description": "WebSocket upgrade"}}},
			},
		},
	})
}

func (s *server) listTools(w http.ResponseWriter, r *http.Request) {
	if s.tools == nil {
		writeError(w, http.StatusServiceUnavailable, "tool registry is not configured")
		return
	}
	searchPath := tools.RuntimePath("", os.Getenv("PATH"))
	query := r.URL.Query()
	if query.Get("tenant") != "" && query.Get("credential_profile") != "" {
		p, ok := s.authenticate(r)
		if !ok {
			writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
			return
		}
		tenantID := query.Get("tenant")
		profileID := query.Get("credential_profile")
		if err := p.AuthorizeTenantProfile(tenantID, profileID); err != nil {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}
		if s.resolveCredentialProfile != nil {
			credentialRoot, err := s.resolveCredentialProfile(tenantID, profileID)
			if err != nil {
				writeError(w, http.StatusBadRequest, fmt.Sprintf("resolve credential profile: %v", err))
				return
			}
			searchPath = tools.RuntimePath(credentialRoot, os.Getenv("PATH"))
		}
	}
	items := s.tools.List()
	out := make([]toolResponse, 0, len(items))
	for _, item := range items {
		health := tools.CheckHealth(item, searchPath)
		out = append(out, toolResponse{
			Name:             item.Name,
			Path:             item.Path,
			Version:          item.Version,
			CredentialEnv:    item.CredentialEnv,
			CredentialSubdir: item.CredentialSubdir,
			Available:        health.Available,
			Health:           health.Health,
			ResolvedPath:     health.ResolvedPath,
			DetectedVersion:  health.DetectedVersion,
			Error:            health.Error,
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{"tools": out})
}

func (s *server) upsertTool(w http.ResponseWriter, r *http.Request) {
	if s.tools == nil {
		writeError(w, http.StatusServiceUnavailable, "tool registry is not configured")
		return
	}
	if !s.requireAdmin(w, r) {
		return
	}
	var tool tools.Tool
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&tool); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid tool request: %v", err))
		return
	}
	if err := s.tools.Upsert(tool); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, tool)
}

func (s *server) deleteTool(w http.ResponseWriter, r *http.Request) {
	if s.tools == nil {
		writeError(w, http.StatusServiceUnavailable, "tool registry is not configured")
		return
	}
	if !s.requireAdmin(w, r) {
		return
	}
	name := r.PathValue("name")
	if err := s.tools.Delete(name); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) listTenants(w http.ResponseWriter, r *http.Request) {
	if s.tenants == nil {
		writeError(w, http.StatusServiceUnavailable, "tenant store is not configured")
		return
	}
	p, ok := s.authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tenants": s.tenants.ListFor(p)})
}

func (s *server) login(w http.ResponseWriter, r *http.Request) {
	if s.tenants == nil {
		writeError(w, http.StatusServiceUnavailable, "tenant store is not configured")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid login request: %v", err))
		return
	}
	token, p, ok := s.tenants.AuthenticateUser(req.Username, req.Password)
	if !ok {
		writeError(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"token":   token,
		"session": sessionBody(p),
	})
}

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	if s.tenants == nil {
		writeError(w, http.StatusServiceUnavailable, "tenant store is not configured")
		return
	}
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid register request: %v", err))
		return
	}

	allowedTools := []string{"*"}
	if s.tools != nil {
		toolsList := s.tools.List()
		allowedTools = make([]string, 0, len(toolsList))
		for _, tool := range toolsList {
			allowedTools = append(allowedTools, tool.Name)
		}
		if len(allowedTools) == 0 {
			allowedTools = []string{"*"}
		}
	}
	token, user, err := s.tenants.RegisterUser(tenants.UserRequest{
		Username:                  req.Username,
		Password:                  req.Password,
		AllowedTools:              allowedTools,
		AllowedWorkspaces:         []string{"repo-*"},
		AllowedCredentialProfiles: []string{"team-default"},
		AllowTerminal:             true,
		MaxJobSeconds:             900,
	})
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(err.Error(), "already exists") {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}
	if s.files.Configured() {
		if err := s.files.EnsureTenant(user.TenantID, user.AllowedCredentialProfiles, user.AllowedWorkspaces); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("create tenant folders: %v", err))
			return
		}
	}
	p, _ := s.tenants.Lookup(token)
	writeJSON(w, http.StatusCreated, map[string]any{
		"token":   token,
		"session": sessionBody(p),
		"user":    user,
	})
}

func (s *server) session(w http.ResponseWriter, r *http.Request) {
	p, ok := s.authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return
	}
	writeJSON(w, http.StatusOK, sessionBody(p))
}

func (s *server) listUsers(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"users": s.tenants.ListUsers()})
}

func (s *server) upsertUser(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	var req tenants.UserRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid user request: %v", err))
		return
	}
	user, err := s.tenants.UpsertUser(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if s.files.Configured() {
		if err := s.files.EnsureTenant(user.TenantID, user.AllowedCredentialProfiles, user.AllowedWorkspaces); err != nil {
			writeError(w, http.StatusInternalServerError, fmt.Sprintf("create tenant folders: %v", err))
			return
		}
	}
	writeJSON(w, http.StatusCreated, user)
}

func (s *server) deleteUser(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if err := s.tenants.DeleteUser(r.PathValue("id")); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) listTokens(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tokens": s.tenants.ListTokens()})
}

func (s *server) upsertToken(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	var req tenants.TokenRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid token request: %v", err))
		return
	}
	token, err := s.tenants.UpsertToken(req)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, token)
}

func (s *server) deleteToken(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdmin(w, r) {
		return
	}
	if err := s.tenants.DeleteToken(r.PathValue("id")); err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeError(w, status, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *server) listFiles(w http.ResponseWriter, r *http.Request) {
	p, ok := s.authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return
	}
	query := r.URL.Query()
	tenantID := query.Get("tenant")
	if tenantID == "" {
		tenantID = p.TenantID
	}
	if err := p.AuthorizeTenant(tenantID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	listing, err := s.files.List(tenantID, query.Get("space"), query.Get("path"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, listing)
}

func (s *server) readFile(w http.ResponseWriter, r *http.Request) {
	p, ok := s.authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return
	}
	query := r.URL.Query()
	tenantID := query.Get("tenant")
	if tenantID == "" {
		tenantID = p.TenantID
	}
	if err := p.AuthorizeTenant(tenantID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	content, err := s.files.Read(tenantID, query.Get("space"), query.Get("path"), 2*1024*1024)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, content)
}

func (s *server) terminalSession(w http.ResponseWriter, r *http.Request) {
	if s.terminal == nil {
		writeError(w, http.StatusServiceUnavailable, "terminal is not configured")
		return
	}
	s.terminal.ServeHTTP(w, r)
}

func (s *server) appServerSession(w http.ResponseWriter, r *http.Request) {
	if s.appServer == nil {
		writeError(w, http.StatusServiceUnavailable, "app-server proxy is not configured")
		return
	}
	s.appServer.ServeHTTP(w, r)
}

func (s *server) createJob(w http.ResponseWriter, r *http.Request) {
	if s.jobs == nil {
		writeError(w, http.StatusServiceUnavailable, "job manager is not configured")
		return
	}
	token, ok := bearerToken(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing bearer token")
		return
	}

	var req jobs.CreateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("invalid job request: %v", err))
		return
	}
	req.Token = token

	job, err := s.jobs.Create(r.Context(), req)
	if err != nil {
		status := http.StatusForbidden
		if strings.Contains(err.Error(), "not registered") || strings.Contains(err.Error(), "not configured") {
			status = http.StatusBadRequest
		}
		writeError(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusAccepted, job)
}

func (s *server) jobByID(w http.ResponseWriter, r *http.Request) {
	rest := strings.TrimPrefix(r.URL.Path, "/api/jobs/")
	if rest == "" || strings.Contains(rest, "/../") {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if strings.HasSuffix(rest, "/events/ws") {
		id := strings.TrimSuffix(rest, "/events/ws")
		s.jobEventsWS(w, r, id)
		return
	}
	if strings.HasSuffix(rest, "/events") {
		id := strings.TrimSuffix(rest, "/events")
		s.jobEvents(w, r, id)
		return
	}
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	job, ok := s.jobs.Get(rest)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	writeJSON(w, http.StatusOK, job)
}

func (s *server) jobEventsWS(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	job, ok := s.jobs.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if !s.authorizeJobEvents(w, r, job) {
		return
	}
	events, ch, unsubscribe, ok := s.jobs.Subscribe(id)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	defer unsubscribe()

	upgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}()

	for _, event := range events {
		if err := conn.WriteJSON(event); err != nil {
			return
		}
		if event.Type == jobs.EventExit {
			return
		}
	}
	for {
		select {
		case event := <-ch:
			if err := conn.WriteJSON(event); err != nil {
				return
			}
			if event.Type == jobs.EventExit {
				return
			}
		case <-done:
			return
		case <-r.Context().Done():
			return
		}
	}
}

func (s *server) jobEvents(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	job, ok := s.jobs.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "job not found")
		return
	}
	if !s.authorizeJobEvents(w, r, job) {
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	encoder := json.NewEncoder(w)
	for _, event := range s.jobs.Events(id) {
		if _, err := fmt.Fprintf(w, "event: %s\n", event.Type); err != nil {
			return
		}
		if _, err := fmt.Fprint(w, "data: "); err != nil {
			return
		}
		if err := encoder.Encode(event); err != nil {
			return
		}
		if _, err := fmt.Fprint(w, "\n"); err != nil {
			return
		}
	}
}

func (s *server) authorizeJobEvents(w http.ResponseWriter, r *http.Request, job jobs.Job) bool {
	p, ok := s.authenticate(r)
	if !ok {
		token := strings.TrimSpace(r.URL.Query().Get("token"))
		if token != "" && s.tenants != nil {
			p, ok = s.tenants.Lookup(token)
		}
	}
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return false
	}
	if err := p.AuthorizeTenant(job.TenantID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return false
	}
	return true
}

func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}

func (s *server) authenticate(r *http.Request) (policy.Policy, bool) {
	if s.tenants == nil {
		return policy.Policy{}, false
	}
	token, ok := bearerToken(r)
	if !ok {
		return policy.Policy{}, false
	}
	return s.tenants.Lookup(token)
}

func (s *server) requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	p, ok := s.authenticate(r)
	if !ok {
		writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
		return false
	}
	if !p.IsAdmin() {
		writeError(w, http.StatusForbidden, "admin role is required")
		return false
	}
	return true
}

type toolResponse struct {
	Name             string `json:"name"`
	Path             string `json:"path"`
	Version          string `json:"version"`
	CredentialEnv    string `json:"credential_env,omitempty"`
	CredentialSubdir string `json:"credential_subdir,omitempty"`
	Available        bool   `json:"available"`
	Health           string `json:"health"`
	ResolvedPath     string `json:"resolved_path,omitempty"`
	DetectedVersion  string `json:"detected_version,omitempty"`
	Error            string `json:"error,omitempty"`
}

func sessionBody(p policy.Policy) map[string]any {
	return map[string]any{
		"subject":                     p.SubjectID,
		"tenant":                      p.TenantID,
		"role":                        p.Role,
		"admin":                       p.IsAdmin(),
		"allowed_tools":               p.AllowedTools,
		"allowed_workspaces":          p.AllowedWorkspaces,
		"allowed_credential_profiles": p.AllowedCredentialProfiles,
		"allow_terminal":              p.AllowTerminal,
		"max_job_seconds":             int(p.MaxJobDuration.Seconds()),
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
