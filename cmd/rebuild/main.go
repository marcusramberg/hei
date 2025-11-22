package rebuild

import (
	"context"
	"errors"
	"runtime"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var Command = &cli.Command{
	Name:      "rebuild",
	ArgsUsage: "[switch|boot|..]",
	Usage:     "Rebuild your nix configuration",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "fast",
			Aliases: []string{"f"},
			Usage:   "Build in fast mode",
		},
		&cli.BoolFlag{
			Name:    "offline",
			Aliases: []string{"o"},
			Usage:   "Build in offline mode",
		},
		&cli.BoolFlag{
			Name:    "rollback",
			Aliases: []string{"r"},
			Usage:   "Build in fast mode",
		},
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Pull on nix flake before rebuilding",
		},
	},
	Action: rebuildAction,
}

func rebuildAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if c.Bool("update") {
		if err := utils.ExecWithStdout(c, "git", []string{"-C", flake, "pull"}); err != nil {
			return err
		}
	}

	args := []string{"nixos-rebuild"}
	if runtime.GOOS == "darwin" {
		args = []string{"darwin-rebuild"}
	}
	args = append(args, "--flake", flake)

	if !c.Args().Present() {
		return errors.New("rebuild called without required arguments (switch/boot/build/...)")
	}
	if c.Bool("rollback") {
		args = append(args, "--rollback")
	}
	if c.Bool("fast") {
		args = append(args, "--fast")
	}
	if c.Bool("offline") {
		args = append(args, "--option", "substitute", "false")
	}
	return utils.ExecWithStdout(c, "sudo", append(args, c.Args().Slice()...))
}
