// Package repl implements the "repl" command to open a Nix REPL in the specified flake.
package repl

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "repl",
	ArgsUsage: "[flake-path...]",
	Usage:     "open a repl in your nix config",
	Action:    replAction,
}

func replAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	return utils.ExecWithStdio(c, "nix", []string{"repl", flake})
}
