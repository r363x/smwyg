package button

import (
    "time"
    "log"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)


type Model struct {
    label           string
    style           gloss.Style
    stylePressed    gloss.Style
    styleSelected   gloss.Style
    styleUnselected gloss.Style
    focused         bool
    action          func()
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            return m, ButtonPress
        }
    case Msg:
        switch msg.Type {

        case ButtonPressed:
            m.style = m.stylePressed
            return m, ButtonRelease

        case ButtonReleased:
            m.action()
            time.Sleep(time.Millisecond * 100)
            m.style = m.styleSelected

        case ButtonSelected:
            m.style = m.styleSelected

        case ButtonUnselected:
            m.Blur()

        }
    }

    return m, cmd
}

func (m Model) View() string {
    return m.style.Render(m.label)
}

func (m *Model) Focus() tea.Cmd {
    m.focused = true
    return ButtonSelect
}

func (m *Model) Blur() {
    m.style = m.styleUnselected
}

func New(label string) *Model {

    s := gloss.NewStyle().Align(gloss.Center)

    sP := s.
        Border(gloss.RoundedBorder()).
        Background(gloss.Color("#c3ccdb"))

    sS := s.
        Border(gloss.NormalBorder()).
        Background(gloss.Color("#8544b8"))

    sU := s.
        Border(gloss.NormalBorder()).
        Background(gloss.Color("#5a3478"))

        return &Model{
            label: label,
            style: sU,
            stylePressed: sP,
            styleSelected: sS,
            styleUnselected: sU,
            action: func() {
                log.Print("NO ACTION ATTACHED!")
            },
        }
}

