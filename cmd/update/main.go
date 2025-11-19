package update

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "update",
	ArgsUsage: "[flake-path...]",
	Usage:     "Update the given flake paths or the default ones if none are provided",
	Action:    checkAction,
}

func checkAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
