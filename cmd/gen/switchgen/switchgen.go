// Package switchgen implements the "switch" command to switch to a specified NixOS generation.
package switchgen

import (
	"context"
	"errors"
	"fmt"
	"os"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var (
	errArgMissing = errors.New("required argument missing")
	Command       = &cli.Command{
		Name:          "switch",
		ArgsUsage:     "[generation]",
		Usage:         "Switch generation",
		Action:        switchAction,
		ShellComplete: completeInputs,
	}
)

func switchAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("you must provide 1 argument, the generation to switch to: %w", errArgMissing)
	}
	gen := fmt.Sprintf("/nix/var/nix/profiles/system-%s-link/bin/switch-to-configuration", c.Args().First())
	_, err := os.Stat(gen)
	if err != nil {
		return fmt.Errorf("generation %s could not be found: %w", c.Args().First(), err)
	}
	return utils.ExecWithStdio(ctx, c, "sudo", []string{gen, "switch"})
}

func completeInputs(ctx context.Context, c *cli.Command) {
	gens, _ := utils.ListGenerations(ctx, c)
	for _, gen := range *gens {
		fmt.Println(gen.Generation)
	}
}
