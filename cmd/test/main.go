package test

import (
	"context"
	"log"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "test",
	ArgsUsage: "[flake-path...]",
	Usage:     "Run a test in an interactive shell",
	Action:    testAction,
}

func testAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	return nil
}
