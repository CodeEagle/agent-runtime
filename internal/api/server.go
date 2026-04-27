package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"agent-runtime/internal/jobs"
	"agent-runtime/internal/tools"
)

type ToolLister interface {
	List() []tools.Tool
}

type Options struct {
	Jobs  *jobs.Manager
	Tools ToolLister
}

func NewServer(options Options) http.Handler {
	server := &server{jobs: options.Jobs, tools: options.Tools}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", server.webUI)
	mux.HandleFunc("GET /api/health", server.health)
	mux.HandleFunc("GET /api/ready", server.ready)
	mux.HandleFunc("GET /api/status", server.status)
	mux.HandleFunc("GET /api/tools", server.listTools)
	mux.HandleFunc("POST /api/jobs", server.createJob)
	mux.HandleFunc("/api/jobs/", server.jobByID)
	return mux
}

type server struct {
	jobs  *jobs.Manager
	tools ToolLister
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
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok",
		"tools":  toolCount,
	})
}

func (s *server) listTools(w http.ResponseWriter, r *http.Request) {
	if s.tools == nil {
		writeError(w, http.StatusServiceUnavailable, "tool registry is not configured")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tools": s.tools.List()})
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

func (s *server) jobEvents(w http.ResponseWriter, r *http.Request, id string) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if _, ok := s.jobs.Get(id); !ok {
		writeError(w, http.StatusNotFound, "job not found")
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

func bearerToken(r *http.Request) (string, bool) {
	header := r.Header.Get("Authorization")
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(header, prefix))
	return token, token != ""
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
