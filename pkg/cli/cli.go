package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"atomicptr.dev/bits"
	"atomicptr.dev/deeperr"
	"atomicptr.dev/peon/pkg/cli/commands/hook"
	"atomicptr.dev/peon/pkg/cli/commands/list"
	"atomicptr.dev/peon/pkg/cli/commands/remove"
	switchcmd "atomicptr.dev/peon/pkg/cli/commands/switch"
	"atomicptr.dev/peon/pkg/cli/common/errormsg"
	"atomicptr.dev/peon/pkg/config"
	"atomicptr.dev/peon/pkg/git"
	"atomicptr.dev/peon/pkg/meta"
	"atomicptr.dev/peon/pkg/util"
	"charm.land/lipgloss/v2"
	"github.com/lmittmann/tint"
	"github.com/urfave/cli/v3"
)

func Run() error {
	debugMode := bits.GetEnvBool("PEON_DEBUG", false)

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

	err := cmd.Run(context.Background(), os.Args)

	// if the error was one of our annotated errors, display a nicer error message
	if e, ok := errors.AsType[deeperr.Error](err); ok {
		b := lipgloss.NewStyle().Bold(true)

		messages := []string{
			fmt.Sprintf("%s: %s", b.Render("    Code"), fmt.Sprintf("E%d", e.Code())),
		}

		chunks := util.SplitWordsEvery(e.Message(), 80)
		first := true

		for _, chunk := range chunks {
			prefix := strings.Repeat(" ", 9)

			if first {
				prefix = b.Render(" Message:")
				first = false
			}

			messages = append(messages, fmt.Sprintf("%s %s", prefix, chunk))
		}

		if err := e.Unwrap(); err != nil {
			messages = append(messages, fmt.Sprintf("%s: %s", b.Render("   Error"), err.Error()))
		}

		if debugMode {
			file, line := e.Location()
			messages = append(messages, fmt.Sprintf("%s: %s:%d", b.Render("Location"), util.CutPathAfter(file, "/pkg/"), line))
		}

		fmt.Fprintln(
			os.Stderr,
			errormsg.Render(
				"Oops! Something went wrong!",
				messages...,
			))
		os.Exit(1)
	}

	return err
}
