package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"slices"
	"strings"
	"testing"
)

type slogLine struct {
	Args []string `json:"args"`
	Cmd  string   `json:"cmd"`
}

func TestCommands(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedLines int
		expectedCmd   string
		expectedArgs  []string
		checkFirst    bool
	}{
		{
			name:          "rebuild",
			args:          []string{"hei", "-d", "-f", ".", "rebuild"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--flake", ".", "switch"},
		},
		{
			name:          "rebuild-boot-default",
			args:          []string{"hei", "-d", "-f", ".", "rebuild", "boot"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--flake", ".", "boot"},
		},
		{
			name:          "show",
			args:          []string{"hei", "-d", "-f", ".", "show"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"flake", "show", "."},
		},
		{
			name:          "gc-system",
			args:          []string{"hei", "-d", "gc", "--system"},
			expectedLines: 3,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nix-env", "--delete-generations", "old", "--profile", "/nix/var/nix/profiles/system"},
		},
		{
			name:          "check",
			args:          []string{"hei", "-d", "-f", ".", "check"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"flake", "check", "."},
		},
		{
			name:          "update",
			args:          []string{"hei", "-d", "-f", ".", "update"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"flake", "update", "--flake", "."},
		},
		{
			name:          "update-pull",
			args:          []string{"hei", "-d", "-f", ".", "update", "--pull"},
			expectedLines: 2,
			expectedCmd:   "nix",
			expectedArgs:  []string{"flake", "update", "--flake", "."},
		},
		{
			name:          "gen-list",
			args:          []string{"hei", "-d", "gen", "list"},
			expectedLines: 2,
			expectedCmd:   "nixos-rebuild",
			expectedArgs:  []string{"list-generations", "--json"},
			checkFirst:    true,
		},
		{
			name:          "gen-delete",
			args:          []string{"hei", "-d", "gen", "delete", "42"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nix-env", "--delete-generations", "--profile", "/nix/var/nix/profiles/system", "42"},
		},
		{
			name:          "gen-delete-range",
			args:          []string{"hei", "-d", "gen", "delete", "10-15"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nix-env", "--delete-generations", "--profile", "/nix/var/nix/profiles/system", "10", "11", "12", "13", "14", "15"},
		},
		{
			name:          "search",
			args:          []string{"hei", "-d", "search", "hello"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"search", "nixpkgs", "hello"},
		},
		{
			name:          "rollback",
			args:          []string{"hei", "-d", "-f", ".", "rollback"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--rollback", "--flake", ".", "switch"},
		},
		{
			name:          "upgrade",
			args:          []string{"hei", "-d", "-f", ".", "upgrade"},
			expectedLines: 2,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--flake", ".", "switch"},
		},
		{
			name:          "repl",
			args:          []string{"hei", "-d", "-f", ".", "repl"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"repl", "."},
		},
		{
			name:          "p",
			args:          []string{"hei", "-d", "p", "list"},
			expectedLines: 1,
			expectedCmd:   "nix",
			expectedArgs:  []string{"profile", "list"},
		},
		{
			name:          "test-simple",
			args:          []string{"hei", "-d", "-f", ".", "test", "mytest"},
			expectedLines: 1,
			expectedCmd:   "nom",
			expectedArgs:  []string{"build", ".#mytest"},
		},
		{
			name:          "build",
			args:          []string{"hei", "-d", "-f", ".", "build"},
			expectedLines: 1,
			expectedCmd:   "nom",
			expectedArgs:  []string{"build", "."},
		},
		{
			name:          "ssh",
			args:          []string{"hei", "-d", "ssh", "myhost", "show"},
			expectedLines: 1,
			expectedCmd:   "ssh",
			expectedArgs:  []string{"myhost", "hei", "show"},
		},
		{
			name:          "gen-diff",
			args:          []string{"hei", "-d", "gen", "diff", "1", "2"},
			expectedLines: 1,
			expectedCmd:   "nvd",
			expectedArgs:  []string{"diff", "/nix/var/nix/profiles/system-1-link", "/nix/var/nix/profiles/system-2-link"},
		},
		{
			name:          "rebuild-unknown-flags",
			args:          []string{"hei", "-d", "-f", ".", "rebuild", "--show-trace", "--custom"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--flake", ".", "--show-trace", "--custom", "switch"},
		},
		{
			name:          "rebuild-flag-with-arg",
			args:          []string{"hei", "-d", "-f", ".", "rebuild", "--build-host", "user@host"},
			expectedLines: 1,
			expectedCmd:   "sudo",
			expectedArgs:  []string{"nixos-rebuild", "--flake", ".", "--build-host", "user@host", "switch"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("Testing command: %v", tt.args)
			jsonLogs := captureStderr(t, func() {
				run(tt.args)
			})
			if len(jsonLogs) != tt.expectedLines {
				t.Fatalf("expected %d log line, got: %d", tt.expectedLines, len(jsonLogs))
			}
			if tt.expectedLines == 0 {
				return
			}
			logToCheck := jsonLogs[len(jsonLogs)-1]
			if tt.checkFirst {
				logToCheck = jsonLogs[0]
			}
			if tt.expectedCmd != "" && logToCheck.Cmd != tt.expectedCmd {
				t.Errorf("expected cmd to be %q, got: %s", tt.expectedCmd, logToCheck.Cmd)
			}
			if !slices.Equal(logToCheck.Args, tt.expectedArgs) {
				t.Errorf("expected args to be %v, got: %v", tt.expectedArgs, logToCheck.Args)
			}
		})
	}
}

func parseJSONLog(t *testing.T, log string) []slogLine {
	t.Helper()
	logs := make([]slogLine, 0)
	for line := range strings.SplitSeq(strings.TrimSpace(log), "\n") {
		if line == "" {
			continue
		}
		var j slogLine
		if err := json.Unmarshal([]byte(line), &j); err != nil {
			t.Fatalf("failed to unmarshal log line %q: %v", line, err)
		}
		logs = append(logs, j)
	}
	return logs
}

func captureStderr(t *testing.T, fn func()) []slogLine {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	slog.SetDefault(slog.New(slog.NewJSONHandler(w, nil)))
	defer func() {
		_ = w.Close()
		_ = r.Close()
	}()
	fn()
	_ = w.Close()
	testOutput, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("failed to read from pipe: %v", err)
	}
	return parseJSONLog(t, string(testOutput))
}
