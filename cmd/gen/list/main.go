package list

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:   "list",
	Usage:  "List nix generations",
	Action: switchAction,
}

func switchAction(ctx context.Context, c *cli.Command) error {
	return utils.ExecWithStdout(c, "sudo", []string{"nix-env", "--list-generations", "--profile", "/nix/var/nix/profiles/system"})
}
