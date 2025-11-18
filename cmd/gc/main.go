package gc

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "gc",
	ArgsUsage: "[]",
	Usage:     "Garbe collection",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
