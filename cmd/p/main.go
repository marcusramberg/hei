package p

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "p",
	ArgsUsage: "[nix profile arguments]",
	Usage:     "Shortcut for nix profile commands",
	Action:    profileAction,
}

func profileAction(ctx context.Context, c *cli.Command) error {
	return utils.ExecWithStdio(c, "nix", append([]string{"profile"}, c.Args().Slice()...))
}
