// Package search allows you to search for packages in nixpkgs
package search

import (
	"context"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "search",
	ArgsUsage: "[query]",
	Usage:     "Search nixpkgs for packages",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	return utils.ExecWithStdio(ctx, c, "nix", append([]string{"search", "nixpkgs"}, c.Args().Slice()...))
}
