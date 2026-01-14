package utils

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestGetFlake_Search(t *testing.T) { //nolint:paralleltest // modifies global defaultFlakePaths
	// Create a temp directory to simulate a flake location
	tmpDir, err := os.MkdirTemp("", "hei-flake-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create flake.nix inside
	if err := os.WriteFile(filepath.Join(tmpDir, "flake.nix"), []byte("{}"), 0600); err != nil {
		t.Fatal(err)
	}

	// Backup original defaultFlakePaths and restore after test
	origPaths := defaultFlakePaths
	defer func() { defaultFlakePaths = origPaths }()

	// Set defaultFlakePaths to temp dir
	defaultFlakePaths = []string{tmpDir}

	// To be safe, we can run a dummy app that calls GetFlake.
	app := &cli.Command{
		Name: "test-app",
		Action: func(ctx context.Context, c *cli.Command) error {
			found := GetFlake(c)
			if found != tmpDir {
				t.Errorf("GetFlake returned %q, want %q", found, tmpDir)
			}
			return nil
		},
	}

	if err := app.Run(context.Background(), []string{"app"}); err != nil {
		t.Fatal(err)
	}
}
