package tui

import (
    "strings"

    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/x/ansi"
)


func (m model) View() string {

    tabView := m.tabs[m.cur].populate(m.dimensions)

	if m.overlay.focused {
        return m.overlayView(tabView)
	}

    return tabView
}

func (t *tab) populate(dim dimensions) string {

    // Left side, narrow
    paneDBView := gloss.NewStyle().
        Border(gloss.NormalBorder()).
        BorderForeground(gloss.Color("240")).
        Width(int(float64(dim.width) * 0.1)).
        Height(dim.height - 3)

    t.queryView.SetWidth(dim.width - paneDBView.GetWidth() - 4)
    t.queryView.SetHeight(int(float64(dim.height) * 0.5) - 10)

    newWidth := dim.width - paneDBView.GetWidth() 

    // Resize the table
    cols := t.resultView.Columns()
    // Resize columns
    for i := range len(cols) {
        cols[i].Width = newWidth / len(cols)
    }
    t.resultView.SetColumns(cols)

    t.resultView.SetWidth(newWidth)
    t.resultView.SetHeight(dim.height - t.queryView.Height() - 3)

    // Bottom, narrow
    paneStatusView := gloss.NewStyle().
        Width(dim.width).
        Padding(0).
        Height(1)
    statusItemView := gloss.NewStyle().
        Padding(0).
        Width(paneStatusView.GetWidth() / 3)

    return gloss.JoinVertical(0,
        gloss.JoinHorizontal(0,
            paneDBView.Render(t.dbView.String()),
            gloss.JoinVertical(0,
                t.queryView.View(),
                t.resultView.View(),
        )),
        paneStatusView.Render(gloss.JoinHorizontal(0,
            statusItemView.AlignHorizontal(gloss.Left).Render(t.statusView.content.left),
            statusItemView.AlignHorizontal(gloss.Center).Render(t.statusView.content.center),
            statusItemView.AlignHorizontal(gloss.Right).Render(t.statusView.content.right),
        )),
    )

}

func truncateLeft(line string, padding int) string {
	if strings.Contains(line, "\n") {
		panic("line must not contain newline")
	}

	// NOTE: line has no newline, so [strings.Join] after [strings.Split] is safe.
	wrapped := strings.Split(ansi.Hardwrap(line, padding, true), "\n")
	if len(wrapped) == 1 {
		return ""
	}

	var ansiStyle string
	ansiStyles := ansiStyleRegexp.FindAllString(wrapped[0], -1)
	if l := len(ansiStyles); l > 0 {
		ansiStyle = ansiStyles[l-1]
	}

	return ansiStyle + strings.Join(wrapped[1:], "")
}
