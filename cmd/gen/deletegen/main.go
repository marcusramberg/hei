package deletegen

import (
	"context"
	"errors"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "delete",
	ArgsUsage: "[gen]",
	Usage:     "Build the given flake paths or the default ones if none are provided",
	Action:    delAction,
}

func delAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return errors.New("you must provide 1 argument, the generation to delete")
	}
	return utils.ExecWithStdout(c, "sudo", []string{"nix-env", "--delete-generations", "--profile", "/nix/var/nix/profiles/system", c.Args().Get(0)})
}
