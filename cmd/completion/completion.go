// Package completion provides dynamic fish shell completion support
package completion

import (
	"fmt"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"
)

// FishCompletion generates fish completions with dynamic completion support
func FishCompletion(root *cli.Command) error {
	// Generate the standard static completions
	static, err := root.ToFishCompletion()
	if err != nil {
		return err
	}
	fmt.Print(static)

	fmt.Println("\n# Dynamic completions")
	addDynamicFishCompletions(root, root.Name, nil)

	return nil
}

func addDynamicFishCompletions(cmd *cli.Command, binary string, ancestors []string) {
	for _, sub := range cmd.Commands {
		if sub.Hidden {
			continue
		}

		subPath := slices.Concat(ancestors, []string{sub.Name})

		if sub.ArgsUsage != "" && sub.ShellComplete != nil {
			condition := buildFishCondition(subPath, sub.Commands)
			fmt.Printf("complete -c %s -n '%s' -xa '(%s %s --generate-shell-completion 2>/dev/null)'\n",
				binary, condition, binary, strings.Join(subPath, " "))
		}

		addDynamicFishCompletions(sub, binary, subPath)
	}
}

func buildFishCondition(path []string, siblings []*cli.Command) string {
	var parts []string

	// Must have seen all ancestors
	for _, ancestor := range path {
		parts = append(parts, fmt.Sprintf("__fish_seen_subcommand_from %s", ancestor))
	}

	siblingCount := 0
	for _, sib := range siblings {
		siblingCount += 1 + len(sib.Aliases)
	}
	siblingNames := make([]string, 0, siblingCount)
	for _, sib := range siblings {
		siblingNames = append(siblingNames, sib.Name)
		siblingNames = append(siblingNames, sib.Aliases...)
	}
	if len(siblingNames) > 0 {
		parts = append(parts, fmt.Sprintf("not __fish_seen_subcommand_from %s", strings.Join(siblingNames, " ")))
	}

	return strings.Join(parts, "; and ")
}
