package utils

import (
	"fmt"
	"log"
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

func ExecWithStdout(cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}

	return nil
}
