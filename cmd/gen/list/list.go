// Package list implements the 'gen list' command to list Nix generations.
package list

import (
	"context"
	"log/slog"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:   "list",
	Usage:  "List nix generations",
	Action: listAction,
}

func listAction(ctx context.Context, c *cli.Command) error {
	j, err := utils.ListGenerations(ctx, c)
	if err != nil {
		return err
	}
	slog.Info("Listed generations", "count", len(*j))
	utils.PrintGenerations(*j)
	return nil
}
