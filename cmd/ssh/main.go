package ssh

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "ssh",
	ArgsUsage: "[host, command]",
	Usage:     "Run a hei command on a remote NixOS system",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	return utils.ExecWithStdout(c, "ssh", append([]string{c.Args().First(), "hei"}, c.Args().Tail()...))
}
