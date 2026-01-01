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
			lastLog := jsonLogs[len(jsonLogs)-1]
			if tt.expectedCmd != "" && lastLog.Cmd != tt.expectedCmd {
				t.Errorf("expected cmd to be %q, got: %s", tt.expectedCmd, jsonLogs[0].Cmd)
			}
			if !slices.Equal(lastLog.Args, tt.expectedArgs) {
				t.Errorf("expected args to be %v, got: %v", tt.expectedArgs, jsonLogs[0].Args)
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
