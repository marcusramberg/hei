package build

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "build",
	ArgsUsage: "[flake-path...]",
	Usage:     "Build the given flake paths or the default ones if none are provided",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
