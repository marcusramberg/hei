package rebuild

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "rebuild",
	ArgsUsage: "[flake-path...]",
	Usage:     "Rebuild your nix configuration",
	Action:    rebuildAction,
}

func rebuildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting rebuild action for %v", c.Args())
	return nil
}
