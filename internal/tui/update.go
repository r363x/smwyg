package tui

import (
    "log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
        case tea.KeyCtrlO:
            m.overlay.focused = !m.overlay.focused
        default:
            return m, m.tabs[m.cur].update(msg)
        }

    case tea.WindowSizeMsg:
        m.dimensions.width = msg.Width
        m.dimensions.height = msg.Height

        return m, nil

    case tickMsg:
        return m, tea.Batch(
            m.tabs[m.cur].refreshStatusLeft,
            m.tabs[m.cur].refreshStatusCenter,
            m.tabs[m.cur].refreshStatusRight,
            doTick(),
        )

    case statusMsg:
        m.tabs[m.cur].updateStatus(msg)
        return m, nil
	}

    return m, nil
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
