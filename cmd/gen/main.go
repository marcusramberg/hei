package gen

import (
	"bas.es/marcus/hei/cmd/gen/deletegen"
	"bas.es/marcus/hei/cmd/gen/diff"
	"bas.es/marcus/hei/cmd/gen/list"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:  "gen",
	Usage: "Manage nix generations",
	Commands: []*cli.Command{
		list.Command,
		deletegen.Command,
		diff.Command,
	},
}
