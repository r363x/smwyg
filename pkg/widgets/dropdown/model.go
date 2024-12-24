package dropdown

import (
	"github.com/r363x/dbmanager/pkg/widgets/overlay"

    tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
    styleFrame = gloss.NewStyle().Border(gloss.NormalBorder()).Background(gloss.Color("#7e7880"))
    styleTitle = styleFrame.Align(gloss.Center, gloss.Center).Border(gloss.NormalBorder(), false, false, true)
)

type Model struct {
    overlay.ModelBase
    Items       []Item
    Label       string
    Description string
    cur         int
}

func New(items []Item, desc string) *Model {

    base := overlay.NewBase()
    base.SetStyle(styleFrame)

    m := Model{
        ModelBase: base,
        Description: desc,
    }

    empty := NewItem("---", nil)
    m.Items = append(m.Items, empty)

    for i := range items {
        m.Items = append(m.Items, items[i])
    }

    m.Items[0].Focus()

    return &m
}


func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:

        switch msg.Type {
        case tea.KeyDown:
            if m.cur < len(m.Items) - 1 {
                m.Items[m.cur].Blur()
                m.cur++
                cmd = m.Items[m.cur].Focus()
            }

        case tea.KeyTab:
            m.Items[m.cur].Blur()

            switch n := len(m.Items) - 1; {
            case m.cur < n:
                m.cur++
            case m.cur == n:
                m.cur = 0
            }
            cmd = m.Items[m.cur].Focus()

        case tea.KeyUp:
            if m.cur >= 1 {
                m.Items[m.cur].Blur()
                m.cur--
                cmd = m.Items[m.cur].Focus()
            }

        case tea.KeyShiftTab:
            m.Items[m.cur].Blur()

            switch n := len(m.Items) - 1; {
            case m.cur >= 1:
                m.cur--
            case m.cur == 0:
                m.cur = n
            }
            cmd = m.Items[m.cur].Focus()

        case tea.KeyEnter:
            cmd = DeliverData(m.Selection().Defaults)
            m.Show = false
        }

    case Msg:
        switch msg.Type {
        case Opened:
            m.Show = true
            cmd = m.Items[0].Focus()

        case Closed:
            m.Show = false

        }

    }

    return m, cmd
}

func (m Model) View() string {
    if m.Show {
        m.SetWidth(len(m.Description) + 4)
        title := styleTitle.Width(m.GetWidth()-1).Render(m.Description)

        var items = []string{title}

        for _, item := range m.Items {
            item.style = item.style.Width(m.GetWidth()-1)
            items = append(items, item.style.Render(item.Label))
        }

        return gloss.JoinVertical(0, items...)
    } else {
        item := m.Selection()
        return item.style.Render(item.Label)
    }
}

func (m *Model) Selection() Item {
    return m.Items[m.cur]
}

func (m *Model) Focus() tea.Cmd {
    return nil
}

func (m *Model) Blur() {
}

