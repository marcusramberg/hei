// Package deletegen implements the 'delete' command to delete Nix generations
package deletegen

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"code.bas.es/marcus/hei/utils"
	"github.com/urfave/cli/v3"
)

var errArgMissing = errors.New("required argument missing")
var errInvalidRange = errors.New("invalid range format")

var Command = &cli.Command{
	Name:      "delete",
	Aliases:   []string{"rm"},
	ArgsUsage: "[gen] or [start-end]",
	Usage:     "Delete the specified generation(s). Supports single gen (42) or range (10-15)",
	Action:    delAction,
}

func delAction(ctx context.Context, c *cli.Command) error {
	if c.Args().Len() != 1 {
		return fmt.Errorf("you must provide 1 argument, the generation to delete (e.g., 42 or 10-15): %w", errArgMissing)
	}

	arg := c.Args().Get(0)
	gens, err := parseGenerations(arg)
	if err != nil {
		return err
	}

	args := append([]string{"nix-env", "--delete-generations", "--profile", "/nix/var/nix/profiles/system"}, gens...)
	return utils.ExecWithStdio(ctx, c, "sudo", args)
}

// parseGenerations parses a generation argument which can be either a single
// generation number or a range in the format "start-end".
func parseGenerations(arg string) ([]string, error) {
	if start, end, ok := strings.Cut(arg, "-"); ok {
		startNum, err := strconv.Atoi(start)
		if err != nil {
			return nil, fmt.Errorf("invalid start of range %q: %w", start, errInvalidRange)
		}
		endNum, err := strconv.Atoi(end)
		if err != nil {
			return nil, fmt.Errorf("invalid end of range %q: %w", end, errInvalidRange)
		}
		if startNum > endNum {
			return nil, fmt.Errorf("start %d must be less than or equal to end %d: %w", startNum, endNum, errInvalidRange)
		}

		gens := make([]string, 0, endNum-startNum+1)
		for i := startNum; i <= endNum; i++ {
			gens = append(gens, strconv.Itoa(i))
		}
		return gens, nil
	}

	// Single generation - validate it's a number
	if _, err := strconv.Atoi(arg); err != nil {
		return nil, fmt.Errorf("invalid generation %q: must be a number or range (e.g., 42 or 10-15): %w", arg, errInvalidRange)
	}
	return []string{arg}, nil
}
