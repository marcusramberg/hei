package check

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "check",
	ArgsUsage: "[flaKe-path]",
	Usage:     "Run checks on the given flake paths or the default ones if none are provided",
	Action:    checkAction,
}

func checkAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Present() {
		return utils.ExecWithStdio(c, "nix", append([]string{"flake", "check"}, c.Args().Slice()...))
	}
	flake := utils.GetFlake(c)
	return utils.ExecWithStdio(c, "nix", []string{"flake", "check", flake})
}
