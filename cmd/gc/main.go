package gc

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "gc",
	ArgsUsage: "[]",
	Usage:     "Garbage collection",
	Action:    gcAction,
}

func gcAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting garbage collection for %v", c.Args())
	return nil
}
