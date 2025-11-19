package search

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "search",
	ArgsUsage: "[query]",
	Usage:     "Search nixpkgs for packages",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
