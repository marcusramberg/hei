// Package rebuild implements the rebuild command
package rebuild

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var errToolMissing = errors.New("nvd tool must be installed for diffs")

var Command = &cli.Command{
	Name:      "rebuild",
	ArgsUsage: "<[switch|boot|..]> - defaults to switch",
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
			Usage:   "Roll back to previous system",
		},
		&cli.BoolFlag{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "Pull on nix flake before rebuilding",
		},
		&cli.BoolFlag{
			Name:    "confirm",
			Aliases: []string{"c"},
			Usage:   "Run a diff and confirm before switching",
		},
	},
	Action: rebuildAction,
}

func rebuildAction(ctx context.Context, c *cli.Command) error {
	flake := utils.GetFlake(c)
	if c.Bool("update") {
		if err := utils.ExecWithStdio(c, "git", []string{"-C", flake, "pull"}); err != nil {
			return err
		}
	}

	args := []string{"nixos-rebuild"}
	if runtime.GOOS == "darwin" {
		args = []string{"darwin-rebuild"}
	}
	args = append(args, "--flake", flake)

	if c.Bool("rollback") {
		args = append(args, "--rollback")
	}
	if c.Bool("fast") {
		args = append(args, "--fast")
	}
	if c.Bool("offline") {
		args = append(args, "--option", "substitute", "false")
	}
	if c.Bool("confirm") {
		return buildConfirm(c, args)
	}

	if !c.Args().Present() {
		args = append(args, "switch")
	}
	return utils.ExecWithStdio(c, "sudo", append(args, c.Args().Slice()...))
}

func buildConfirm(c *cli.Command, args []string) error {
	// setup a temp dir to not pollute cwd with result/
	tmpdir, err := os.MkdirTemp(os.TempDir(), "rebuild-")
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			slog.Warn("Failed to clean up temp build dir", "err", err)
		}
	}()
	if err != nil {
		return fmt.Errorf("couldn't make temp dir for build: %w", err)
	}

	if err = os.Chdir(tmpdir); err != nil {
		return err
	}

	nvd, err := exec.LookPath("nvd")
	if err != nil {
		return fmt.Errorf("nvd tool must be installed for diffs%w ", errToolMissing)
	}

	if err := utils.ExecWithStdio(c, "sudo", append(args, "build")); err != nil {
		return err
	}

	if err := utils.ExecWithStdio(c, nvd, []string{"diff", "/nix/var/nix/profiles/system", "result"}); err != nil {
		return err
	}

	fmt.Println("Press Enter to confirm the switch or ctrl-c to abort...")
	reader := bufio.NewReader(os.Stdin)
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	return utils.ExecWithStdio(c, "sudo", append(args, "switch"))
}
