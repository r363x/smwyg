package tui

import (
    "strings"
    "regexp"

    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/x/ansi"
)

var (
    width = 80
    height = 40
    ansiStyleRegexp = regexp.MustCompile(`\x1b[[\d;]*m`)
    overlayStyle = gloss.NewStyle().
        Bold(true).
        Border(gloss.NormalBorder()).
        Foreground(gloss.Color("#FAFAFA")).
        Background(gloss.Color("#7D56F4")).
        Padding(2).
        Width(width).
        Height(height)
)

func (m *model) overlayView(bg string) string {
	// wrappedBG := ansi.Hardwrap(bg, m.dimensions.width, true)

    row := (m.dimensions.height - height) / 2
    col := (m.dimensions.width - width) / 2

	overlay := overlayStyle.Render("Hello there")

	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, overlayLine := range overlayLines {
		bgLine := bgLines[i+row] // TODO: index handling
		if len(bgLine) < col {
			bgLine += strings.Repeat(" ", col-len(bgLine)) // add padding
		}

		bgLeft := ansi.Truncate(bgLine, col, "")
		bgRight := truncateLeft(bgLine, col+ansi.StringWidth(overlayLine))

		bgLines[i+row] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(bgLines, "\n")
}

