package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/urfave/cli/v3"
)

type Generation struct {
	Generation            int    `json:"generation"`
	Date                  string `json:"date"`
	NixosVersion          string `json:"nixosVersion"`
	KernelVersion         string `json:"kernelVersion"`
	ConfigurationRevision string `json:"configurationRevision"`
	Specializations       string `json:"specializations"`
	Current               bool   `json:"current"`
}

func ListGenerations(c *cli.Command) (*[]Generation, error) {
	var res *[]Generation
	out, err := ExecGetOutput(c, "nixos-rebuild", []string{"list-generations", "--json"})
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func PrintGenerations(gens []Generation) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Generation\tAge\tNixOS Version\tKernel Version\tConfiguration Revision\tCurrent\n")
	for g := len(gens) - 1; g >= 0; g-- {
		isCurrent := ""
		if gens[g].Current {
			isCurrent = "(*)"
		}
		fmt.Fprintf(w, "%d%s\t%s\t%s\t%s\n", gens[g].Generation, isCurrent, gens[g].RelativeAge(), gens[g].NixosVersion, gens[g].KernelVersion)
	}
	w.Flush()
}

func (g *Generation) IsOlderThan(ts time.Time) bool {
	if t, err := time.Parse("2006-01-02 15:04:05 MST", g.Date); err == nil {
		return t.Before(ts)
	}
	return false
}

func (g *Generation) RelativeAge() string {
	ts, err := time.ParseInLocation("2006-01-02 15:04:05", g.Date, time.Local)
	if err != nil {
		return "unknown"
	}
	now := time.Now()
	diff := now.Sub(ts)
	if diff < 0 {
		diff = -diff
	}
	switch {
	case diff < 2*time.Minute:
		return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
	case diff < 3*time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 36*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	default:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	}
}
