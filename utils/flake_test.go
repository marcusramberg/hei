package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHome(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("Could not get user home directory")
	}

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path",
			input:    "/etc/nixos",
			expected: "/etc/nixos",
		},
		{
			name:     "relative path",
			input:    "foo/bar",
			expected: "foo/bar",
		},
		{
			name:     "home",
			input:    "~",
			expected: home,
		},
		{
			name:     "subdirectory in home",
			input:    "~/foo/bar",
			expected: filepath.Join(home, "foo/bar"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := ExpandHome(tt.input)
			if got != tt.expected {
				t.Errorf("ExpandHome(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
