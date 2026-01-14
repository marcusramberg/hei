// Package update lets you update a input to your flake, or all of them
package update

import (
	"context"
	"fmt"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "update",
	ArgsUsage: "[inputs...]",
	Usage:     "Update the given flake inputs or all of them if none are provided",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "pull",
			Aliases: []string{"p"},
			Usage:   "Do 'git pull' before updating the flake",
		},
	},
	Action:        updateAction,
	ShellComplete: completeInputs,
}

func updateAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if c.Bool("pull") {
		if err := utils.ExecWithStdio(c, "git", []string{"-C", flake, "pull"}); err != nil {
			return err
		}
	}
	return utils.ExecWithStdio(c, "nix", append([]string{"flake", "update", "--flake", flake}, c.Args().Slice()...))
}

func completeInputs(_ context.Context, c *cli.Command) {
	flake := utils.GetFlake(c)
	for _, input := range utils.GetFlakeInputs(flake) {
		fmt.Println(input)
	}
}
