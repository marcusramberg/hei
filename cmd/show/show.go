// Package show implements the 'show' command to show your nix flake info
package show

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "show",
	ArgsUsage: "[flake-path...]",
	Usage:     "Run nix flake show on the given flake paths or the default ones if none are provided",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	return utils.ExecWithStdio(c, "nix", []string{"flake", "show", flake})
}
