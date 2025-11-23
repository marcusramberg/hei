package deletegen

import (
	"context"
	"errors"
	"fmt"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var errArgMissing = errors.New("required argument missing")

var Command = &cli.Command{
	Name:      "delete",
	ArgsUsage: "[gen]",
	Usage:     "Build the given flake paths or the default ones if none are provided",
	Action:    delAction,
}

func delAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("%w: you must provide 1 argument, the generation to delete", errArgMissing)
	}
	return utils.ExecWithStdout(c, "sudo", []string{"nix-env", "--delete-generations", "--profile", "/nix/var/nix/profiles/system", c.Args().Get(0)})
}
