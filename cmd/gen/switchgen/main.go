package switchgen

import (
	"context"
	"errors"
	"fmt"
	"os"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "switch",
	ArgsUsage: "[]",
	Usage:     "Switch generation",
	Action:    switchAction,
}

func switchAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return errors.New("you must provide 1 argument, the generation to switch to")
	}
	gen := fmt.Sprintf("/nix/var/nix/profiles/system-%s-link/bin/switch-to-configuration", c.Args().First())
	finfo, err := os.Stat(gen)
	if err != nil {
		return fmt.Errorf("generation %s could not be found: %w", c.Args().First(), err)
	}
	if finfo.Mode()&0o111 == 0 {
		return fmt.Errorf("gen switch script is not executable: %s", gen)
	}
	return utils.ExecWithStdout(c, "sudo", []string{gen, "switch"})
}
