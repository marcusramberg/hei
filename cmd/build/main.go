package build

import (
	"context"
	"log"
	"os/exec"

	"bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "build",
	ArgsUsage: "[flake-path...]",
	Usage:     "Build the given flake paths or the default ones if none are provided",
	Action:    buildAction,
}

func buildAction(ctx context.Context, c *cli.Command) error {
	log.Printf("Starting build action for %v", c.Args())
	builder, err := exec.LookPath("nom")
	if err != nil {
		builder, err = exec.LookPath("nix")
		if err != nil {
			log.Fatalf("Cannot build, neither 'nom' nor 'nix' found in PATH")
		}
	}
	if c.Args().Present() {
		return utils.ExecWithStdout(builder, append([]string{"build"}, c.Args().Slice()...))
	}
	flake := utils.GetFlake(c)
	return utils.ExecWithStdout(builder, []string{"build", flake})
}
