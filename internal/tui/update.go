package tui

import (
    "log"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

        case tea.KeyCtrlQ:
            m.resultView.Blur()
            m.queryView.Focus()

        case tea.KeyCtrlR:
            m.queryView.Blur()
            m.resultView.Focus()

        case tea.KeyF5:
            if m.queryView.Value() != "" && m.queryView.Focused() {
                if err := m.buildResultsTable(m.queryView.Value()); err != nil {
                    log.Fatalln("Fatal: ", err)
                }
            }
            m.refreshDbView()
            return m, nil
        }

    case tea.WindowSizeMsg:
        m.dimensions.width = msg.Width
        m.dimensions.height = msg.Height

        return m, nil

    case tickMsg:
        return m, tea.Batch(
            m.refreshStatusLeft,
            m.refreshStatusCenter,
            m.refreshStatusRight,
            doTick(),
        )

    case statusMsg:
        switch msg.section {
        case secLeft:
            m.statusView.content.left = msg.message
        case secCenter:
            m.statusView.content.center = msg.message
        case secRight:
            m.statusView.content.right = msg.message
        }
        return m, nil
	}


    switch {
    case m.queryView.Focused():
        m.queryView, cmd = m.queryView.Update(msg)
    case m.resultView.Focused():
        m.resultView, cmd = m.resultView.Update(msg)
    }

    return m, cmd
}
