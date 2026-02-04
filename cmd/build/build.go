// Package build implements the 'build' command to build Nix flakes with nom
package build

import (
	"context"
	"fmt"
	"os/exec"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "build",
	ArgsUsage: "[flake-path...]",
	Usage:     "Run nix flake check on your flake",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	var builder string
	var err error
	if c.Bool("dry-run") {
		builder = "nom"
	} else {
		builder, err = exec.LookPath("nom")
		if err != nil {
			builder, err = exec.LookPath("nix")
			if err != nil {
				return fmt.Errorf("cannot build, neither 'nom' nor 'nix' found in PATH: %w", err)
			}
		}
	}
	if c.Args().Present() {
		return utils.ExecWithStdio(ctx, c, builder, append([]string{"build"}, c.Args().Slice()...))
	}
	flake := utils.GetFlake(c)
	return utils.ExecWithStdio(ctx, c, builder, []string{"build", flake})
}
