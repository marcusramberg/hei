package rollback

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "rollback",
	ArgsUsage: "[flake-path...]",
	Usage:     "Roll back to previous generation of nixos. See gen list for the current generations.",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
