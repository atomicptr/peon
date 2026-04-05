package cli

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"atomicptr.dev/bits"
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
	cmd := &cli.Command{
		Name:    "peon",
		Usage:   "Peon is a Git Worktree management tool designed for humans",
		Version: meta.VersionStringWith(""),
		Commands: []*cli.Command{
			switchcmd.Command(),
			list.Command(),
			remove.Command(),
			hook.Command(),
		},
		Before: func(ctx context.Context, cmd *cli.Command) (context.Context, error) {
			debugMode := bits.GetEnvBool("PEON_DEBUG", false)

			level := slog.LevelInfo

			if debugMode {
				level = slog.LevelDebug
			}

			h := tint.NewHandler(os.Stderr, &tint.Options{
				Level:      level,
				TimeFormat: time.Kitchen,
				AddSource:  true,
			})
			logger := slog.New(h)
			slog.SetDefault(logger)

			cfg, err := config.FromEnv()
			if err != nil {
				return ctx, fmt.Errorf("could not load config: %w", err)
			}

			ctx = context.WithValue(ctx, config.ContextKey, cfg)

			g, err := git.FromEnv()
			if err != nil {
				return ctx, fmt.Errorf("could not find git: %w", err)
			}

			cmd.Version = meta.VersionStringWith(g.Version())

			slog.Debug("git found", "path", g.Executable, "version", g.Version())

			ctx = context.WithValue(ctx, git.ContextKey, g)

			return ctx, nil
		},
	}

	return cmd.Run(context.Background(), os.Args)
}
