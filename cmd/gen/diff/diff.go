// Package diff implements the "diff" command for comparing two system generations.
package diff

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var (
	errToolMissing = errors.New("nvd tool must be installed for diffs")
	errArgMissing  = errors.New("required arguments are missing")
)

var Command = &cli.Command{
	Name:      "diff",
	ArgsUsage: "[gen1] [gen2]",
	Usage:     "Diff two generations",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	var nvd string
	var err error
	if c.Bool("dry-run") {
		nvd = "nvd"
	} else {
		nvd, err = exec.LookPath("nvd")
		if err != nil {
			return fmt.Errorf("nvd tool must be installed for diffs: %w", errToolMissing)
		}
	}
	if c.Args().Len() != 2 {
		return fmt.Errorf("you must provide 2 argument, from and to generation: %w", errArgMissing)
	}
	return utils.ExecWithStdio(c, nvd, []string{
		"diff",
		fmt.Sprintf("/nix/var/nix/profiles/system-%s-link", c.Args().Get(0)),
		fmt.Sprintf("/nix/var/nix/profiles/system-%s-link", c.Args().Get(1)),
	})
}
