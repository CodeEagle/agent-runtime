package policy

import (
	"fmt"
	"path"
	"time"
)

type Policy struct {
	SubjectID                 string
	TenantID                  string
	AllowedTools              []string
	AllowedWorkspaces         []string
	AllowedCredentialProfiles []string
	AllowTerminal             bool
	MaxJobDuration            time.Duration
}

type JobRequest struct {
	TenantID          string
	Tool              string
	WorkspaceID       string
	CredentialProfile string
	RequestedDuration time.Duration
	WantsTerminal     bool
}

type TerminalRequest struct {
	TenantID          string
	WorkspaceID       string
	CredentialProfile string
}

func (p Policy) AuthorizeJob(req JobRequest) error {
	if req.TenantID != p.TenantID {
		return fmt.Errorf("tenant %q is not allowed for subject %q", req.TenantID, p.SubjectID)
	}
	if !matchesAny(p.AllowedTools, req.Tool) {
		return fmt.Errorf("tool %q is not allowed for subject %q", req.Tool, p.SubjectID)
	}
	if !matchesAny(p.AllowedWorkspaces, req.WorkspaceID) {
		return fmt.Errorf("workspace %q is not allowed for subject %q", req.WorkspaceID, p.SubjectID)
	}
	if !matchesAny(p.AllowedCredentialProfiles, req.CredentialProfile) {
		return fmt.Errorf("credential profile %q is not allowed for subject %q", req.CredentialProfile, p.SubjectID)
	}
	if req.WantsTerminal && !p.AllowTerminal {
		return fmt.Errorf("terminal access is not allowed for subject %q", p.SubjectID)
	}
	if p.MaxJobDuration > 0 && req.RequestedDuration > p.MaxJobDuration {
		return fmt.Errorf("requested duration %s exceeds max duration %s", req.RequestedDuration, p.MaxJobDuration)
	}
	return nil
}

func (p Policy) AuthorizeTerminal(req TerminalRequest) error {
	if !p.AllowTerminal {
		return fmt.Errorf("terminal access is not allowed for subject %q", p.SubjectID)
	}
	if req.TenantID != p.TenantID {
		return fmt.Errorf("tenant %q is not allowed for subject %q", req.TenantID, p.SubjectID)
	}
	if !matchesAny(p.AllowedWorkspaces, req.WorkspaceID) {
		return fmt.Errorf("workspace %q is not allowed for subject %q", req.WorkspaceID, p.SubjectID)
	}
	if !matchesAny(p.AllowedCredentialProfiles, req.CredentialProfile) {
		return fmt.Errorf("credential profile %q is not allowed for subject %q", req.CredentialProfile, p.SubjectID)
	}
	return nil
}

func matchesAny(patterns []string, value string) bool {
	for _, pattern := range patterns {
		if pattern == "*" || pattern == value {
			return true
		}
		matched, err := path.Match(pattern, value)
		if err == nil && matched {
			return true
		}
	}
	return false
}
