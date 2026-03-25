package list

import "github.com/urfave/cli/v3"

func Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all available worktrees",
		Aliases: []string{"ls"},
	}
}
