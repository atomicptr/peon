package remove

import "github.com/urfave/cli/v3"

func Command() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "Remove worktree",
		Aliases:   []string{"rm"},
		ArgsUsage: "[name]",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "name",
				UsageText: "Worktree name",
			},
		},
	}
}
