// Package utils provides utility functions for hei
package utils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

var defaultFlakePaths = []string{"/etc/nixos", "~/.config/nix-darwin", "~/.config/nix-config"}

func expandHome(path string) string {
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
	flake := expandHome(c.String("flake"))
	if flake != "" {
		s, err := os.Stat(fmt.Sprintf("%s/flake.nix", flake))
		if err == nil && !s.IsDir() {
			return flake
		}
		log.Fatalf("Provided flake path is invalid: %s", flake)
	}
	for _, flake := range defaultFlakePaths {
		flake = expandHome(flake)
		s, err := os.Stat(fmt.Sprintf("%s/flake.nix", flake))
		if err == nil && !s.IsDir() {
			log.Printf("Using flake: %s", flake)
			return flake
		}
	}
	log.Fatalf("No valid flake found in candidates: %v", candidates)
	return ""
}

func ExecWithStdio(c *cli.Command, cmd string, args []string) error {
	if c.Bool("dry-run") {
		slog.Info("dry-run:", "cmd", cmd, "args", args)
		return nil
	}
	// FIXME: This should probably take context from caller, but meh...
	ec := exec.CommandContext(context.Background(), cmd, args...)
	ec.Stdin = os.Stdin
	ec.Stdout = os.Stdout
	ec.Stderr = os.Stderr
	if err := ec.Run(); err != nil {
		return err
	}

	return nil
}

func ExecGetOutput(c *cli.Command, cmd string, args []string) ([]byte, error) {
	if c.Bool("dry-run") {
		slog.Info("dry-run:", "cmd", cmd, "args", args)
		return nil, nil
	}
	ec := exec.CommandContext(context.Background(), cmd, args...)
	out, err := ec.Output()
	if err != nil {
		return nil, err
	}

	return out, nil
}
