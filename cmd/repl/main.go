package repl

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "repl",
	ArgsUsage: "[flake-path...]",
	Usage:     "open a repl in your nix config",
	Action:    replAction,
}

func replAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
