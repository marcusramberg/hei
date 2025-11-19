package main

import (
	"context"
	"log"
	"os"

	"bas.es/marcus/hei/cmd/build"
	"bas.es/marcus/hei/cmd/check"
	"bas.es/marcus/hei/cmd/gc"
	"bas.es/marcus/hei/cmd/gen"
	"bas.es/marcus/hei/cmd/p"
	"bas.es/marcus/hei/cmd/rebuild"
	"bas.es/marcus/hei/cmd/repl"
	"bas.es/marcus/hei/cmd/rollback"
	"bas.es/marcus/hei/cmd/search"
	"bas.es/marcus/hei/cmd/show"
	"bas.es/marcus/hei/cmd/ssh"
	"bas.es/marcus/hei/cmd/test"
	"bas.es/marcus/hei/cmd/update"
	"bas.es/marcus/hei/cmd/upgrade"
	"github.com/urfave/cli/v3"
)

var Version string = "dev"

func main() {
	hei := cli.Command{}
	hei.Name = "hei"
	hei.Description = "A simple consistent command wrapper for nix"
	hei.Version = Version
	hei.Usage = "A simple consistent command wrapper for nix"
	hei.ConfigureShellCompletionCommand = func(c *cli.Command) {
		c.Hidden = false
		c.Usage = "Generate shell completion scripts"
	}
	hei.Suggest = true
	hei.Commands = []*cli.Command{
		build.Command,
		check.Command,
		gc.Command,
		gen.Command,
		p.Command,
		rebuild.Command,
		repl.Command,
		rollback.Command,
		search.Command,
		show.Command,
		ssh.Command,
		test.Command,
		upgrade.Command,
		update.Command,
	}
	hei.Flags = setupFlags()

	if err := hei.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func setupFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:    "flake",
			Aliases: []string{"f"},
			Usage:   "Path to flake. Will default to auto-detect",
			Sources: cli.EnvVars("HEI_FLAKE"),
		},
		&cli.BoolFlag{
			Name:    "debug",
			Aliases: []string{"D"},
			Usage:   "Enable debug logging",
			Sources: cli.EnvVars("HEI_DEBUG"),
		},
		&cli.BoolFlag{
			Name:    "dry-run",
			Aliases: []string{"d"},
			Usage:   "Perform a dry run without making any changes",
			Sources: cli.EnvVars("HEI_DRY_RUN"),
		},
	}
}
