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
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Do 'git pull' before updating the flake",
		},
	},
	Action: updateAction,
}

func updateAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if c.Bool("pull") {
		return utils.ExecWithStdio(c, "git", []string{"-C", flake, "pull"})
	}
	return utils.ExecWithStdio(c, "nix", append([]string{"flake", "update", "--flake", flake}, c.Args().Slice()...))
}
