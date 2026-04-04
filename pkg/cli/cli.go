package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"atomicptr.dev/peon/pkg/cli/commands/hook"
	"atomicptr.dev/peon/pkg/cli/commands/list"
	"atomicptr.dev/peon/pkg/cli/commands/remove"
	switchcmd "atomicptr.dev/peon/pkg/cli/commands/switch"
	"atomicptr.dev/peon/pkg/config"
	"atomicptr.dev/peon/pkg/git"
	"atomicptr.dev/peon/pkg/meta"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
)

func Run() error {
	h := tint.NewHandler(os.Stderr, &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: time.Kitchen,
		AddSource:  true,
	})
	logger := slog.New(h)
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
			hook.Command(),
		},
		Version: meta.VersionStringWith(g.Version()),
	}

	return cmd.Run(ctx, os.Args)
}
