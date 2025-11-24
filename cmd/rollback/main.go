package rollback

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "rollback",
	ArgsUsage: "[flake-path...]",
	Usage:     "Roll back to previous generation of nixos. See gen list for the current generations.",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	return utils.ExecWithStdio(c, "sudo", []string{"nixos-rebuild", "--rollback", "--flake", flake, "switch"})
}
