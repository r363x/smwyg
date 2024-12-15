package tui

import (
	"github.com/r363x/dbmanager/pkg/widgets/config"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textarea"
)


type model struct {
    tabs []tab
    cur int
    overlay config.Model
    dimensions dimensions
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
        textarea.Blink,
        m.tabs[m.cur].refreshStatusLeft,
        m.tabs[m.cur].refreshStatusCenter,
        m.tabs[m.cur].refreshStatusRight,
        doTick(),
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
        case tea.KeyCtrlO:
            m.overlay.Show = !m.overlay.Show
        default:
            if m.overlay.Show {
                m.overlay, cmd = m.overlay.Update(msg)
            } else {
                cmd = m.tabs[m.cur].update(msg)
            }
        }

    case tea.WindowSizeMsg:
        m.SetDimensions(msg.Width, msg.Height)

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

    case config.ButtonMsg:
        m.overlay, cmd = m.overlay.Update(msg)

	}


    return m, cmd
}


func (m model) View() string {

    tabView := m.tabs[m.cur].populate(m.dimensions)

	if m.overlay.Show {
        m.overlay.SetBackground(tabView)
        return m.overlay.View()
	}

    return tabView
}
