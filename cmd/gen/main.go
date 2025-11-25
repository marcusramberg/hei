// Package gen hooks up subcommands for managing nix generations
package gen

import (
	"code.bas.es/marcus/hei/cmd/gen/deletegen"
	"code.bas.es/marcus/hei/cmd/gen/diff"
	"code.bas.es/marcus/hei/cmd/gen/list"
	"code.bas.es/marcus/hei/cmd/gen/switchgen"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "gen",
	Usage: "Manage nix generations",
	Commands: []*cli.Command{
		list.Command,
		deletegen.Command,
		diff.Command,
		switchgen.Command,
	},
}
