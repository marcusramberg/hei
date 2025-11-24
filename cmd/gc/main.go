package gc

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "gc",
	ArgsUsage: "[]",
	Usage:     "Garbage collect & optimize nix store",
	Action:    gcAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "all",
			Aliases: []string{"a"},
			Usage:   "Collect all garbage",
		},
		&cli.BoolFlag{
			Name:    "system",
			Aliases: []string{"s"},
			Usage:   "Collect system garbage",
		},
	},
}

func gcAction(ctx context.Context, c *cli.Command) error {
	if c.Bool("all") || c.Bool("system") {
		if err := utils.ExecWithStdio(c, "sudo", []string{"nix-collect-garbage", "-d"}); err != nil {
			return err
		}
		if err := utils.ExecWithStdio(c, "sudo", []string{"nix-store", "--optimize"}); err != nil {
			return err
		}
		if err := utils.ExecWithStdio(c, "sudo", []string{"nix-env", "--delete-generations old", "--profile", "/nix/var/nix/profiles/system"}); err != nil {
			return err
		}
	}
	if c.Bool("all") || !c.Bool("system") {
		if err := utils.ExecWithStdio(c, "nix-collect-garbage", []string{"-d"}); err != nil {
			return err
		}
	}
	return nil
}
