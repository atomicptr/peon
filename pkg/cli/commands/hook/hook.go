package hook

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "hook",
		Usage: "Hooks for Shell Integrations",
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
	case bashName:
		return bashEvalCommand
	case fishName:
		return fishEvalCommand
	case zshName:
		return zshEvalCommand
	default:
		return bashEvalCommand
	}
}

func ConfigFile() string {
	p := os.Getenv("SHELL")
	n := filepath.Base(p)

	switch n {
	case bashName:
		return filepath.Join(xdg.Home, ".bashrc")
	case fishName:
		return filepath.Join(xdg.ConfigHome, "fish", "config.fish")
	case zshName:
		return filepath.Join(xdg.Home, ".zshrc")
	default:
		return filepath.Join(xdg.Home, ".bashrc")
	}
}
