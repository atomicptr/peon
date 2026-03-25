package switchcmd

import "github.com/urfave/cli/v3"

func Command() *cli.Command {
	return &cli.Command{
		Name:      "switch",
		Usage:     "Switch to a worktree, create a new one if needed",
		Aliases:   []string{"sw"},
		ArgsUsage: "[name]",
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name:      "name",
				UsageText: "Worktree name",
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "create",
				Aliases: []string{"c"},
				Usage:   "Create a new worktree",
			},
		},
	}
}
