package tui

import (
	"github.com/r363x/dbmanager/pkg/widgets/config"
	"github.com/r363x/dbmanager/pkg/widgets/button"
	"github.com/r363x/dbmanager/internal/tui/tab"
	"github.com/r363x/dbmanager/pkg/widgets/results"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/bubbles/textinput"
)


type model struct {
    tabs []tab.Model
    cur int
    overlay config.Model
}

func (m model) Init() tea.Cmd {

    tab := &m.tabs[m.cur]

    cmd := make([]tea.Cmd, len(tab.Elements))

    for i := range tab.Elements {
        switch element := tab.Elements[i].(type) {
        case *textarea.Model:
            *element, cmd[i] = element.Update(nil)
        case *results.Model:
            *element, cmd[i] = element.Update(nil)
        }
    }

    cmd = append(cmd,
        textarea.Blink,
        textinput.Blink,
        tab.RefreshStatusLeft,
        tab.RefreshStatusCenter(""),
        tab.RefreshStatusRight,
        tab.RefreshBrowser,
        doTick(3),
    )

	return tea.Batch(cmd...)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

    var cmd tea.Cmd
    curTab := m.tabs[m.cur]

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
                m.tabs[m.cur], cmd = curTab.Update(msg)
            }
        }

    case tea.WindowSizeMsg:
        m.SetDimensions(msg.Width, msg.Height)

        return m, nil

    case tickMsg:
        return m, tea.Batch(
            curTab.RefreshStatusLeft,
            curTab.RefreshStatusCenter(""),
            curTab.RefreshStatusRight,
            doTick(3),
        )

    // case status.Msg:
    //     curTab.UpdateStatus(msg)
    //     return m, nil

    case button.Msg:
        m.overlay, cmd = m.overlay.Update(msg)

    default:
        if m.overlay.Show {
            m.overlay, cmd = m.overlay.Update(msg)
        } else {
            m.tabs[m.cur], cmd = curTab.Update(msg)
        }
	}


    return m, cmd
}


func (m model) View() string {

    cur := m.tabs[m.cur]

	if m.overlay.Show {
        m.overlay.SetBackground(cur.View())
        return m.overlay.View()
	}

    return cur.View()
}

func (m *model) SetDimensions(width , height int) {
    m.tabs[m.cur].SetDimentions(width, height)
    m.overlay.SetDimensions(width, height)
}

