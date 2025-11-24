package utils

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
)

var defaultFlakePaths = []string{"/etc/nixos", "~/.config/nix-darwin", "~/.config/nix-config"}

func GetFlake(c *cli.Command) string {
	candidates := defaultFlakePaths
	flake := c.String("flake")
	if flake != "" {
		s, err := os.Stat(fmt.Sprintf("%s/flake.nix", flake))
		if err == nil && !s.IsDir() {
			return flake
		}
		log.Fatalf("Provided flake path is invalid: %s", flake)
	}
	for _, flake := range defaultFlakePaths {
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
