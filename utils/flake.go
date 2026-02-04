// Package utils provides utility functions for hei
package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var defaultFlakePaths = []string{"/etc/nixos", "~/.config/nix-darwin", "~/.config/nix-config"}

func ExpandHome(path string) string {
	if len(path) == 0 || path[0] != '~' {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if len(path) == 1 {
		return home
	}

	return filepath.Join(home, path[1:])
}

func GetFlake(c *cli.Command) string {
	candidates := defaultFlakePaths
	flake := ExpandHome(c.String("flake"))
	if flake != "" {
		s, err := os.Stat(fmt.Sprintf("%s/flake.nix", flake))
		if err == nil && !s.IsDir() {
			return flake
		}
		log.Fatalf("Provided flake path is invalid: %s", flake)
	}
	for _, flake := range defaultFlakePaths {
		flake = ExpandHome(flake)
		s, err := os.Stat(fmt.Sprintf("%s/flake.nix", flake))
		if err == nil && !s.IsDir() {
			log.Printf("Using flake: %s", flake)
			return flake
		}
	}
	log.Fatalf("No valid flake found in candidates: %v", candidates)
	return ""
}

func ExecWithStdio(ctx context.Context, c *cli.Command, cmd string, args []string) error {
	if c.Bool("dry-run") {
		slog.Info("dry-run:", "cmd", cmd, "args", args)
		return nil
	}
	ec := exec.CommandContext(ctx, cmd, args...)
	ec.Stdin = os.Stdin
	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	if err := ec.Run(); err != nil {
		return err
	}

	return nil
}

func ExecGetOutput(ctx context.Context, c *cli.Command, cmd string, args []string) ([]byte, error) {
	if c.Bool("dry-run") {
		slog.Info("dry-run:", "cmd", cmd, "args", args)
		return nil, nil
	}
	ec := exec.CommandContext(ctx, cmd, args...)
	out, err := ec.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}

// flakeLock represents the structure of a flake.lock file.
type flakeLock struct {
	Nodes map[string]flakeLockNode `json:"nodes"`
}

type flakeLockNode struct {
	Inputs map[string]any `json:"inputs"`
}

// GetFlakeInputs reads the flake.lock file and returns the list of input names.
func GetFlakeInputs(flakePath string) []string {
	lockPath := filepath.Join(flakePath, "flake.lock")
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return nil
	}

	var lock flakeLock
	if err := json.Unmarshal(data, &lock); err != nil {
		return nil
	}

	root, ok := lock.Nodes["root"]
	if !ok {
		return nil
	}

	inputs := make([]string, 0, len(root.Inputs))
	for name := range root.Inputs {
		inputs = append(inputs, name)
	}
	return inputs
}
