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
			zshCommand(),
		},
	}
}

func EvalCommand() string {
	p := os.Getenv("SHELL")
	n := filepath.Base(p)

	switch n {
	case "bash":
		return bashEvalCommand
	case "fish":
		return fishEvalCommand
	case "zsh":
		return zshEvalCommand
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
	case "zsh":
		return filepath.Join(xdg.Home, ".zshrc")
	default:
		return filepath.Join(xdg.Home, ".bashrc")
	}
}
