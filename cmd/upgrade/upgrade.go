// Package upgrade updates all your inputs and rebuilds your system.
package upgrade

import (
	"context"

	"code.bas.es/marcus/hei/cmd/rebuild"
	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:   "upgrade",
	Usage:  "update all flake inputs and rebuild the system",
	Action: upgradeAction,
}

func upgradeAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if err := utils.ExecWithStdio(ctx, c, "nix", append([]string{"flake", "update", "--flake", flake}, c.Args().Slice()...)); err != nil {
		return err
	}
	return rebuild.Command.Action(ctx, c)
}
