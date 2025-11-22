package upgrade

import (
	"context"
	"errors"

	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:   "upgrade",
	Usage:  "deprecated, use rebuild -u instead",
	Action: upgradeAction,
}

func upgradeAction(ctx context.Context, c *cli.Command) error {
	return errors.New("deprecated: use rebuild -u instead")
}
