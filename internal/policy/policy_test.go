package policy_test

import (
	"strings"
	"testing"
	"time"

	"agent-runtime/internal/policy"
)

func TestAuthorizeJobAllowsMatchingRequest(t *testing.T) {
	p := policy.Policy{
		SubjectID:                 "service-account:cc-connect",
		TenantID:                  "team-a",
		AllowedTools:              []string{"codex", "claude"},
		AllowedWorkspaces:         []string{"repo-*", "shared"},
		AllowedCredentialProfiles: []string{"team-default"},
		AllowTerminal:             false,
		MaxJobDuration:            15 * time.Minute,
	}

	err := p.AuthorizeJob(policy.JobRequest{
		TenantID:          "team-a",
		Tool:              "codex",
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
		RequestedDuration: 5 * time.Minute,
	})

	if err != nil {
		t.Fatalf("expected request to be authorized, got %v", err)
	}
}

func TestAuthorizeJobRejectsOutOfPolicyRequests(t *testing.T) {
	p := policy.Policy{
		SubjectID:                 "service-account:worker",
		TenantID:                  "team-a",
		AllowedTools:              []string{"codex"},
		AllowedWorkspaces:         []string{"repo-main"},
		AllowedCredentialProfiles: []string{"team-default"},
		AllowTerminal:             false,
		MaxJobDuration:            10 * time.Minute,
	}

	tests := []struct {
		name string
		req  policy.JobRequest
		want string
	}{
		{
			name: "tenant mismatch",
			req: policy.JobRequest{
				TenantID:          "team-b",
				Tool:              "codex",
				WorkspaceID:       "repo-main",
				CredentialProfile: "team-default",
				RequestedDuration: time.Minute,
			},
			want: "tenant",
		},
		{
			name: "tool not allowed",
			req: policy.JobRequest{
				TenantID:          "team-a",
				Tool:              "claude",
				WorkspaceID:       "repo-main",
				CredentialProfile: "team-default",
				RequestedDuration: time.Minute,
			},
			want: "tool",
		},
		{
			name: "workspace not allowed",
			req: policy.JobRequest{
				TenantID:          "team-a",
				Tool:              "codex",
				WorkspaceID:       "repo-other",
				CredentialProfile: "team-default",
				RequestedDuration: time.Minute,
			},
			want: "workspace",
		},
		{
			name: "credential profile not allowed",
			req: policy.JobRequest{
				TenantID:          "team-a",
				Tool:              "codex",
				WorkspaceID:       "repo-main",
				CredentialProfile: "personal",
				RequestedDuration: time.Minute,
			},
			want: "credential",
		},
		{
			name: "terminal not allowed",
			req: policy.JobRequest{
				TenantID:          "team-a",
				Tool:              "codex",
				WorkspaceID:       "repo-main",
				CredentialProfile: "team-default",
				RequestedDuration: time.Minute,
				WantsTerminal:     true,
			},
			want: "terminal",
		},
		{
			name: "duration too long",
			req: policy.JobRequest{
				TenantID:          "team-a",
				Tool:              "codex",
				WorkspaceID:       "repo-main",
				CredentialProfile: "team-default",
				RequestedDuration: 30 * time.Minute,
			},
			want: "duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := p.AuthorizeJob(tt.req)
			if err == nil {
				t.Fatalf("expected authorization error containing %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("expected error to contain %q, got %v", tt.want, err)
			}
		})
	}
}

func TestAuthorizeTerminalRequiresTenantWorkspaceProfileAndPermission(t *testing.T) {
	p := policy.Policy{
		SubjectID:                 "service-account:web",
		TenantID:                  "team-a",
		AllowedWorkspaces:         []string{"repo-*"},
		AllowedCredentialProfiles: []string{"team-default"},
		AllowTerminal:             true,
	}

	err := p.AuthorizeTerminal(policy.TerminalRequest{
		TenantID:          "team-a",
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
	})
	if err != nil {
		t.Fatalf("expected terminal request to be authorized, got %v", err)
	}

	p.AllowTerminal = false
	err = p.AuthorizeTerminal(policy.TerminalRequest{
		TenantID:          "team-a",
		WorkspaceID:       "repo-main",
		CredentialProfile: "team-default",
	})
	if err == nil || !strings.Contains(err.Error(), "terminal") {
		t.Fatalf("expected terminal permission error, got %v", err)
	}
}
