package shell

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

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

// TODO: add support for zsh (and more?)
func EvalCommand() string {
	p := os.Getenv("SHELL")
	n := filepath.Base(p)

	switch n {
	case "bash":
		return bashEvalCommand
	case "fish":
		return fishEvalCommand
	default:
		return bashEvalCommand
	}
}

func ConfigFile() string {
	p := os.Getenv("SHELL")
	n := filepath.Base(p)

	switch n {
	case "bash":
		return filepath.Join(xdg.Home, ".bashrc")
	case "fish":
		return filepath.Join(xdg.ConfigHome, "fish", "config.fish")
	default:
		return filepath.Join(xdg.Home, ".bashrc")
	}
}
