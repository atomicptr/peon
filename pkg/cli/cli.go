package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/atomicptr/peon/pkg/cli/commands/list"
	"github.com/atomicptr/peon/pkg/cli/commands/remove"
	switchcmd "github.com/atomicptr/peon/pkg/cli/commands/switch"
	"github.com/atomicptr/peon/pkg/config"
	"github.com/atomicptr/peon/pkg/git"
	"github.com/atomicptr/peon/pkg/meta"
	"github.com/urfave/cli/v3"
)

func Run() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)

	ctx := context.Background()

	cfg, err := config.FromEnv()
	if err != nil {
		return fmt.Errorf("could not load config: %w", err)
	}

	ctx = context.WithValue(ctx, config.ContextKey, cfg)

	g, err := git.FromEnv()
	if err != nil {
		return fmt.Errorf("could not find git: %w", err)
	}

	slog.Debug("git found", "path", g.Executable, "version", g.Version())

	ctx = context.WithValue(ctx, git.ContextKey, g)

	cmd := &cli.Command{
		Name:  "peon",
		Usage: "Peon is a Git Worktree management tool designed for humans",
		Commands: []*cli.Command{
			switchcmd.Command(),
			list.Command(),
			remove.Command(),
		},
		Version: meta.VersionStringWith(g.Version()),
	}

	return cmd.Run(ctx, os.Args)
}
