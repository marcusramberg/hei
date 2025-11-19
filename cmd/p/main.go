package p

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "p",
	ArgsUsage: "[nix profile arguments]",
	Usage:     "Shortcut for nix profile commands",
	Action:    profileAction,
}

func profileAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
