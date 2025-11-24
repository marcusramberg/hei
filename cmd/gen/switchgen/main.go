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
		Name:      "switch",
		ArgsUsage: "[]",
		Usage:     "Switch generation",
		Action:    switchAction,
	}
)

func switchAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("%w: you must provide 1 argument, the generation to switch to", errArgMissing)
	}
	gen := fmt.Sprintf("/nix/var/nix/profiles/system-%s-link/bin/switch-to-configuration", c.Args().First())
	_, err := os.Stat(gen)
	if err != nil {
		return fmt.Errorf("generation %s could not be found: %w", c.Args().First(), err)
	}
	return utils.ExecWithStdio(c, "sudo", []string{gen, "switch"})
}
