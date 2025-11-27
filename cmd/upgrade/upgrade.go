// Package upgrade updates all your inputs and rebuilds your system.
package upgrade

import (
	"context"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

var errDeprecated = errors.New("deprecated")

var Command = &cli.Command{
	Name:   "upgrade",
	Usage:  "deprecated, use rebuild -u instead",
	Action: upgradeAction,
}

func upgradeAction(ctx context.Context, c *cli.Command) error {
	return fmt.Errorf("use rebuild -u instead: %w", errDeprecated)
}
