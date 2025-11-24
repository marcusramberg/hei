package test

import (
	"context"
	"errors"
	"fmt"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var errArgMissing = errors.New("missing required argument")

var Command = &cli.Command{
	Name:      "test",
	ArgsUsage: "[test target]",
	Usage:     "Run a nix check",
	Action:    testAction,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "interactive",
			Aliases: []string{"i"},
			Usage:   "Run the test with the interctive test driver",
		},
	},
}

func testAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if c.Bool("interactive") {
		if c.Args().Len() != 1 {
			return fmt.Errorf("%w,takes one argument, the nix test to run", errArgMissing)
		}
		err := utils.ExecWithStdio(c, "nix", []string{"build", fmt.Sprintf("%s#%s.driverInteractive", flake, c.Args().First())})
		if err != nil {
			return fmt.Errorf("failed to build the interactive test driver: %w", err)
		}
		return utils.ExecWithStdio(c, "./result/bin/nixos-test-driver", nil)
	}
	return utils.ExecWithStdio(c, "nom", []string{"build", fmt.Sprintf("%s#%s", flake, c.Args().First())})
}
