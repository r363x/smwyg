package tui

import (
    "log"

	tea "github.com/charmbracelet/bubbletea"
)


func (m *model) SetDimensions(width int, height int) {
    m.dimensions.width = width
    m.dimensions.height = height
    m.overlay.SetDimensions(width, height)

}

func (t *tab) updateStatus(msg statusMsg) {
    switch msg.section {
    case secLeft:
        t.statusView.content.left = msg.message
    case secCenter:
        t.statusView.content.center = msg.message
    case secRight:
        t.statusView.content.right = msg.message
    }
}

func (t *tab) update(msg tea.KeyMsg) tea.Cmd {
    var cmd tea.Cmd

    switch msg.Type {

    case tea.KeyCtrlQ:
        t.resultView.Blur()
        t.queryView.Focus()

    case tea.KeyCtrlR:
        t.queryView.Blur()
        t.resultView.Focus()

    case tea.KeyF5:
        if t.queryView.Value() != "" && t.queryView.Focused() {
            if err := t.buildResultsTable(t.queryView.Value()); err != nil {
                log.Fatalln("Fatal: ", err)
            }
        }
        t.refreshDbView()
        return nil
    }

    switch {
    case t.queryView.Focused():
        t.queryView, cmd = t.queryView.Update(msg)
    case t.resultView.Focused():
        t.resultView, cmd = t.resultView.Update(msg)
    }

    return cmd
}
