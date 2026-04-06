package errormsg

import (
	"fmt"
	"regexp"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/samber/lo"
)

const (
	colorError = "#f38ba8"
	colorText  = "#cdd6f4"
	colorBase  = "#1e1e2e"
)

var backtickRegex = regexp.MustCompile("`([^`]+)`")

func formatHighlights(input string) string {
	return backtickRegex.ReplaceAllStringFunc(input, func(m string) string {
		inner := strings.Trim(m, "`")
		return fmt.Sprintf("`%s`", lipgloss.NewStyle().Bold(true).Render(inner))
	})
}

func Render(title string, lines ...string) string {
	container := lipgloss.NewStyle().
		Padding(1, 1).
		Margin(1, 0, 1, 1).
		Border(lipgloss.RoundedBorder(), false, false, false, true).
		BorderForeground(lipgloss.Color(colorError))

	badge := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorError)).
		Padding(0, 1).
		Bold(true).
		MarginRight(1).
		Render("ERROR")

	headerText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Bold(true).
		Render(title)

	bodyLines := lipgloss.JoinVertical(
		lipgloss.Left,
		lo.Map(lines, func(line string, _ int) string {
			return lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorText)).
				Render(formatHighlights(line))
		})...,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, badge, headerText),
		"",
		bodyLines,
	)

	return container.Render(content)
}
