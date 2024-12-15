package config

import (
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	"github.com/r363x/dbmanager/pkg/widgets/button"

	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/bubbles/textinput"
)


type Selectable interface {
    Focus() tea.Cmd
    Blur()
}


type View struct {
    Name        string
    Content     string
    Selectables []Selectable
    curFocus    int
}

type Model struct {
    overlay.ModelBase
    views []View
    cur   int
}

func (m *Model) FocusOn(itemIndex int) {
    m.views[m.cur].curFocus = itemIndex
}

func New(views []View) Model {
    base := overlay.NewBase()
    return Model{ base, views, 0 }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
	case tea.KeyMsg:

        cur := m.views[m.cur].curFocus
        sel := m.views[m.cur].Selectables[cur]

		switch msg.Type {
        case tea.KeyDown:
            next := cur + 1

            if next < len(m.views[m.cur].Selectables) {
                sel.Blur()
                cmd = m.views[m.cur].Selectables[next].Focus()
                m.views[m.cur].curFocus = next
            }

        case tea.KeyUp:
            prev := cur - 1

            if prev >= 0 {
                sel.Blur()
                cmd = m.views[m.cur].Selectables[prev].Focus()
                m.views[m.cur].curFocus = prev
            }

        default:
            switch sel := sel.(type) {
            case *textinput.Model:
                updated, _cmd := sel.Update(msg)
                m.views[m.cur].Selectables[cur] = &updated
                cmd = _cmd
            case *button.Model:
                updated, _cmd := sel.Update(msg)
                m.views[m.cur].Selectables[cur] = &updated
                cmd = _cmd
            }
        }
    case button.Msg:
        btn, _cmd := msg.Button.Update(msg)
        m.views[m.cur].Selectables[m.views[m.cur].curFocus] = &btn
        cmd = _cmd
    }

    return m, cmd
}

func (m Model) View() string {

    m.SetContents(m.views[m.cur].Content)
    return m.BaseView()
}

func (m Model) Init() tea.Cmd {
    return nil
}

