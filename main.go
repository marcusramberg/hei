package main

import (
	"context"
	"log"
	"os"

	"bas.es/marcus/hei/cmd/build"
	"github.com/urfave/cli/v3"
)

var Version string = "dev"

var defaultFlakePaths = []string{"/etc/nixos", "~/.config/nix-darwin", "~/.config/nix-config"}

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
	hei.Commands = []*cli.Command{
		build.Command,
	}
	if err := hei.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("error: %v", err)
	}
}
