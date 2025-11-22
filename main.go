package main

import (
	"context"
	"log"
	"os"

	"code.bas.es/marcus/hei/cmd/build"
	"code.bas.es/marcus/hei/cmd/check"
	"code.bas.es/marcus/hei/cmd/gc"
	"code.bas.es/marcus/hei/cmd/gen"
	"code.bas.es/marcus/hei/cmd/p"
	"code.bas.es/marcus/hei/cmd/rebuild"
	"code.bas.es/marcus/hei/cmd/repl"
	"code.bas.es/marcus/hei/cmd/rollback"
	"code.bas.es/marcus/hei/cmd/search"
	"code.bas.es/marcus/hei/cmd/show"
	"code.bas.es/marcus/hei/cmd/ssh"
	"code.bas.es/marcus/hei/cmd/swap"
	"code.bas.es/marcus/hei/cmd/test"
	"code.bas.es/marcus/hei/cmd/update"
	"code.bas.es/marcus/hei/cmd/upgrade"
	docs "github.com/urfave/cli-docs/v3"
	"github.com/urfave/cli/v3"
)

var Version string = "dev"

func main() {
	hei := cli.Command{}
	hei.Name = "hei"
	hei.Description = "A simple consistent command wrapper for nix"
	hei.UseShortOptionHandling = true
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
		swap.Command,
		test.Command,
		upgrade.Command,
		update.Command,
		{
			Name:   "gen-docs",
			Hidden: true,
			Action: func(ctx context.Context, c *cli.Command) error {
				md, err := docs.ToMarkdown(&hei)
				if err != nil {
					return err
				}
				err = os.WriteFile("docs.md", []byte(md), 0o644)
				if err != nil {
					return err
				}
				man, err := docs.ToMan(&hei)
				if err != nil {
					return err
				}
				err = os.WriteFile("hei.1.man", []byte(man), 0o644)
				if err != nil {
					return err
				}
				return nil
			},
		},
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
