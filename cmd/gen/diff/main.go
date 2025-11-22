package diff

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "diff",
	ArgsUsage: "[gen1] [gen2]",
	Usage:     "Build the given flake paths or the default ones if none are provided",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	nvd, err := exec.LookPath("nvd")
	if err != nil {
		return errors.New("nvd tool must be installed for diffs")
	}
	if c.Args().Len() != 2 {
		return errors.New("you must provide 2 argument, from and to generation")
	}
	return utils.ExecWithStdout(c, nvd, []string{
		"diff",
		fmt.Sprintf("/nix/var/nix/profiles/system-%s-link", c.Args().Get(0)),
		fmt.Sprintf("/nix/var/nix/profiles/system-%s-link", c.Args().Get(1)),
	})
}
