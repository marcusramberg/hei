package check

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "check",
	ArgsUsage: "[flaKe-path]",
	Usage:     "Run checks on the given flake paths or the default ones if none are provided",
	Action:    checkAction,
}

func checkAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting check action for %v", c.Args())
	return nil
}
