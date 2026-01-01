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
	"strings"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var errToolMissing = errors.New("nvd tool must be installed for diffs")

var Command = &cli.Command{
	Name:            "rebuild",
	ArgsUsage:       "<[switch|boot|..]> - defaults to switch",
	Usage:           "Rebuild your nix configuration",
	SkipFlagParsing: true,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "fast",
			Usage: "Build in fast mode",
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
	var (
		flake           string
		fast            bool
		offline         bool
		rollback        bool
		update          bool
		confirm         bool
		extraFlags      []string
		explicitActions []string
		actionSeen      bool
	)

	// Manual argument parsing
	args := c.Args().Slice()
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--flake", "-f":
			if i+1 < len(args) {
				flake = args[i+1]
				i++
			}
		case "--fast":
			fast = true
		case "--offline", "-o":
			offline = true
		case "--rollback", "-r":
			rollback = true
		case "--update", "-u":
			update = true
		case "--confirm", "-c":
			confirm = true
		default:
			if strings.HasPrefix(arg, "-") {
				extraFlags = append(extraFlags, arg)
			} else if isKnownAction(arg) {
				explicitActions = append(explicitActions, arg)
				actionSeen = true
			} else {
				// profound heuristic: if it's not a flag and not a known action,
				// it's probably an argument to a flag (e.g. --build-host <host>)
				extraFlags = append(extraFlags, arg)
			}
		}
	}

	// Resolve flake
	if flake == "" {
		flake = utils.GetFlake(c)
	} else {
		flake = utils.ExpandHome(flake)
	}

	if update {
		if err := utils.ExecWithStdio(c, "git", []string{"-C", flake, "pull"}); err != nil {
			return err
		}
	}

	cmdArgs := []string{"nixos-rebuild"}
	if runtime.GOOS == "darwin" {
		cmdArgs = []string{"darwin-rebuild"}
	}
	cmdArgs = append(cmdArgs, "--flake", flake)

	if rollback {
		cmdArgs = append(cmdArgs, "--rollback")
	}
	if fast {
		cmdArgs = append(cmdArgs, "--fast")
	}
	if offline {
		cmdArgs = append(cmdArgs, "--option", "substitute", "false")
	}

	cmdArgs = append(cmdArgs, extraFlags...)

	if confirm {
		return buildConfirm(c, cmdArgs)
	}

	cmdArgs = append(cmdArgs, explicitActions...)

	if !actionSeen {
		cmdArgs = append(cmdArgs, "switch")
	}

	return utils.ExecWithStdio(c, "sudo", cmdArgs)
}

var knownActions = map[string]bool{
	"switch":                   true,
	"boot":                     true,
	"test":                     true,
	"build":                    true,
	"dry-build":                true,
	"dry-activate":             true,
	"edit":                     true,
	"repl":                     true,
	"list-generations":         true,
	"build-vm":                 true,
	"build-vm-with-bootloader": true,
	"changelog":                true, // darwin
	"check":                    true, // darwin
}

func isKnownAction(s string) bool {
	return knownActions[s]
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

	var nvd string
	if c.Bool("dry-run") {
		nvd = "nvd"
	} else {
		nvd, err = exec.LookPath("nvd")
		if err != nil {
			return fmt.Errorf("nvd tool must be installed for diffs%w ", errToolMissing)
		}
	}

	if err := utils.ExecWithStdio(c, "sudo", append(args, "build")); err != nil {
		return err
	}

	if err := utils.ExecWithStdio(c, nvd, []string{"diff", "/nix/var/nix/profiles/system", "result"}); err != nil {
		return err
	}

	fmt.Println("Press Enter to confirm the switch or ctrl-c to abort...")
	reader := bufio.NewReader(os.Stdin)
	if !c.Bool("dry-run") {
		_, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
	} else {
		fmt.Println("(dry-run mode, not waiting for confirmation)")
	}

	return utils.ExecWithStdio(c, "sudo", append(args, "switch"))
}
