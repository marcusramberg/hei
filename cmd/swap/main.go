package swap

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

var (
	errArgMissing    = errors.New("missing required arguments")
	errInvalidBackup = errors.New("invalid backup")
)

const backupSuffix = ".nix-store-backup"

var Command = &cli.Command{
	Name:      "swap",
	ArgsUsage: "[targets]",
	Usage:     "Recursively swap nix-store symlinks with copies (or back)",
	Action:    testAction,
}

func testAction(ctx context.Context, c *cli.Command) error {
	if !c.Args().Present() {
		return fmt.Errorf("%w: must specify targets to swap", errArgMissing)
	}
	for _, t := range c.Args().Slice() {
		err := filepath.Walk(t, recurseAndSwap)
		if err != nil {
			return err
		}
	}
	return nil
}

func recurseAndSwap(path string, info os.FileInfo, err error) error {
	if strings.Contains(path, backupSuffix) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("invalid target %s: %w", path, err)
	}
	if info.IsDir() {
		return nil
	}
	backup := fmt.Sprintf("%s%s", path, backupSuffix)
	if info.Mode()&fs.ModeSymlink == 0 {
		return restoreBackup(path, backup)
	}
	return swapFile(path, backup)
}

func swapFile(path, backup string) error {
	target, err := filepath.EvalSymlinks(path)
	if err != nil {
		return fmt.Errorf("failed to resolve symlink for %s: %w", path, err)
	}
	contents, err := os.ReadFile(target)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", target, err)
	}
	err = os.Rename(path, backup)
	if err != nil {
		return fmt.Errorf("failed to make backup of %s: %w", path, err)
	}
	err = os.WriteFile(path, contents, 0o600)
	if err != nil {
		return fmt.Errorf("failed to write %s from symlink: %w", path, err)
	}
	slog.Info("Swapped file and made backup", "file", path)
	return nil
}

func restoreBackup(path, backup string) error {
	b, err := os.Lstat(backup)
	if err != nil {
		return fmt.Errorf("file %s found, and no backup to restore. Bailing out: %w", path, err)
	}
	if b.Mode()&fs.ModeSymlink == 0 {
		return fmt.Errorf("%w: backup %s isn't a symlink, bailing out", errInvalidBackup, backup)
	}
	err = os.Rename(backup, path)
	if err != nil {
		return fmt.Errorf("failed to restore backup of %s: %w", path, err)
	}
	slog.Info("Restored from backup", "file", path)
	return nil
}
