package tui

import (
    gloss "github.com/charmbracelet/lipgloss"
)

func (m model) View() string {

    // Left side, narrow
    paneDBView := gloss.NewStyle().
        Border(gloss.NormalBorder()).
        BorderForeground(gloss.Color("240")).
        Width(int(float64(m.dimensions.width) * 0.08)).
        Height(m.dimensions.height - 3)

    m.queryView.SetWidth(m.dimensions.width - paneDBView.GetWidth() - 4)
    m.queryView.SetHeight(int(float64(m.dimensions.height) * 0.5) - 10)

    newWidth := m.dimensions.width - paneDBView.GetWidth() 

    // Resize the table
    cols := m.resultView.Columns()
    // Resize columns
    for i := range len(cols) {
        cols[i].Width = newWidth / len(cols)
    }
    m.resultView.SetColumns(cols)

    m.resultView.SetWidth(newWidth)
    m.resultView.SetHeight(m.dimensions.height - m.queryView.Height() - 3)

    // Bottom, narrow
    paneStatusView := gloss.NewStyle().
        Width(m.dimensions.width).
        Padding(0).
        Height(1)
    statusItemView := gloss.NewStyle().
        Padding(0).
        Width(paneStatusView.GetWidth() / 3)

    return gloss.JoinVertical(0,
        gloss.JoinHorizontal(0,
            paneDBView.Render(m.dbView.String()),
            gloss.JoinVertical(0,
                m.queryView.View(),
                m.resultView.View(),
        )),
        paneStatusView.Render(gloss.JoinHorizontal(0,
            statusItemView.AlignHorizontal(gloss.Left).Render(m.statusView.content.left),
            statusItemView.AlignHorizontal(gloss.Center).Render(m.statusView.content.center),
            statusItemView.AlignHorizontal(gloss.Right).Render(m.statusView.content.right),
        )),
    )

}
