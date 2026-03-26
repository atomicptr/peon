package shell

import "github.com/urfave/cli/v3"

func Command() *cli.Command {
	return &cli.Command{
		Name:  "shell",
		Usage: "Shell Integrations",
		Commands: []*cli.Command{
			bashCommand(),
			fishCommand(),
		},
	}
}
