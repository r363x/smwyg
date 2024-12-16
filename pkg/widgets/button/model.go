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
    styleFocused    gloss.Style
    styleBlurred    gloss.Style
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
            cmd = ButtonRelease

        case ButtonReleased:
            m.action()
            time.Sleep(time.Millisecond * 500)
            m.Focus()
        }
    }

    return m, cmd
}

func (m Model) View() string {
    return m.style.Render(m.label)
}

func (m *Model) Focus() tea.Cmd {
    m.focused = true
    m.style = m.styleFocused
    return nil
}

func (m *Model) Blur() {
    m.focused = false
    m.style = m.styleBlurred
}

func New(label string) *Model {

    s := gloss.NewStyle().
        Align(gloss.Center).
        PaddingLeft(2).
        PaddingRight(2).
        MarginLeft(1).
        MarginRight(1)

    sP := s.Background(gloss.Color("#c3ccdb"))

    sS := s.Background(gloss.Color("#8544b8"))

    sU := s.Background(gloss.Color("#5a3478"))

        return &Model{
            label: label,
            style: sU,
            stylePressed: sP,
            styleFocused: sS,
            styleBlurred: sU,
            action: func() {
                log.Print("NO ACTION ATTACHED!")
            },
        }
}

