package update

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "update",
	ArgsUsage: "[flake-path...]",
	Usage:     "Update the given flake paths or the default ones if none are provided",
	Action:    updateAction,
}

func updateAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	return utils.ExecWithStdout(c, "git", []string{"-C", flake, "pull"})
}
