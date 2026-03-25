package list

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/atomicptr/peon/pkg/git"
	"github.com/atomicptr/peon/pkg/util"
	"github.com/charmbracelet/x/term"
	"github.com/urfave/cli/v3"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "List all available worktrees",
		Aliases: []string{"ls"},
		Action:  listCommand,
	}
}

func listCommand(ctx context.Context, cmd *cli.Command) error {
	git, err := git.FromContext(ctx)
	if err != nil {
		return err
	}

	if !git.IsGitDir() {
		return fmt.Errorf("%s is not part of a git tree", git.WorkingDir)
	}

	worktrees, err := git.GetWorktrees()
	if err != nil {
		return err
	}

	var rows []wtRow

	for _, wt := range worktrees {
		if wt.IsBare {
			continue
		}

		c, err := git.FetchCommitInfo(wt.Head)
		if err != nil {
			slog.Error("could not fetch commit info", "err", err, "hash", wt.Head)
			continue
		}

		rows = append(rows, wtRow{
			Worktree: wt,
			Commit:   c,
		})
	}

	return printTable(rows)
}

func makeTableRow(wt git.Worktree, c git.CommitInfo) []string {
	path := wt.Path

	if cwd, err := os.Getwd(); err == nil {
		p, err := filepath.Rel(cwd, wt.Path)
		if err == nil {
			path = p
		}
	}

	branch := wt.Branch

	branch = strings.TrimPrefix(branch, "refs/heads/")
	branch = strings.TrimPrefix(branch, "refs/remotes/")
	branch = strings.TrimPrefix(branch, "refs/tags/")

	if branch == "" && wt.IsDetached {
		branch = "[detached]"
	}

	return []string{
		branch,
		path,
		wt.Head[0:7],
		c.Subject,
		util.RelativeTime(c.Timestamp),
	}
}

const (
	indexBranch = iota
	indexPath
	indexCommit
	indexMessage
	indexAge
)

var listHeaders = []string{"Branch", "Path", "Commit", "Message", "Age"}

var listHeaderStyle = lipgloss.
	NewStyle().
	Foreground(lipgloss.Color("4")).
	Bold(true).
	Align(lipgloss.Left)

var listRowStyle = lipgloss.
	NewStyle().
	Padding(0, 0)

type wtRow struct {
	Worktree git.Worktree
	Commit   git.CommitInfo
}

func printTable(wtRows []wtRow) error {
	var rows [][]string

	colWidths := make([]int, len(listHeaders))

	for _, wtRow := range wtRows {
		row := makeTableRow(wtRow.Worktree, wtRow.Commit)

		for col, s := range row {
			w := lipgloss.Width(s)

			if w > colWidths[col] {
				colWidths[col] = w
			}
		}

		rows = append(rows, row)
	}

	w, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		return fmt.Errorf("could not get terminal size: %w", err)
	}

	const (
		marginWidth  = 2
		borderWidth  = 0
		paddingWidth = 0
	)

	chrome := (marginWidth * 2) + (borderWidth * 2) + (paddingWidth * 2) + (len(listHeaders) - 1)

	remainingWidth := max(w - colWidths[indexCommit] - colWidths[indexAge] - chrome)

	branchWidth := max(min(colWidths[indexBranch], remainingWidth/4), len(listHeaders[indexBranch]))
	pathWidth := max(min(colWidths[indexPath], remainingWidth/3), len(listHeaders[indexPath]))

	messageWidth := remainingWidth - branchWidth - pathWidth

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		BorderRow(false).
		Headers(listHeaders...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			s := listRowStyle
			if row == table.HeaderRow {
				s = listHeaderStyle
			}

			switch col {
			case indexBranch:
				s = s.Width(branchWidth)
			case indexPath:
				s = s.Width(pathWidth)
			case indexMessage:
				s = s.Width(messageWidth)
			default:
				s = s.Width(max(colWidths[col], len(listHeaders[col])))
			}

			return s
		})

	if _, err := lipgloss.Println(t); err != nil {
		return err
	}

	return nil
}
