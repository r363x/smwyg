package config

import (
    "time"
    "log"

	"github.com/r363x/dbmanager/internal/tui/overlay"
	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/textinput"
)

type ButtonMsg interface{}

type ButtonPressMsg struct {
    button *Button
}

type ButtonReleaseMsg struct {
    button *Button
}

type ButtonSelectMsg struct {
    button *Button
}

type ButtonUnselectMsg struct {
    button *Button
}

type Button struct {
    label           string
    style           gloss.Style
    stylePressed    gloss.Style
    styleSelected   gloss.Style
    styleUnselected gloss.Style
    action          func()
}

func (b *Button) ButtonPressed() tea.Msg {
    return ButtonMsg(ButtonPressMsg{button: b})
}

func (b *Button) ButtonReleased() tea.Msg {
    return ButtonMsg(ButtonReleaseMsg{button: b})
}

func (b *Button) View() string {
    return b.style.Render(b.label)
}

func NewButton(label string) Button {

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

        return Button{
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

type View struct {
    Name    string
    Content string
    Buttons []Button
    Inputs  []textinput.Model
}

type Model struct {
    overlay.ModelBase
    views []View
    cur   int
}

func New(views []View) Model {
    base := overlay.NewBase()
    return Model{ base, views, 0 }
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case ButtonPressMsg:
        msg.button.style = msg.button.stylePressed
        cmd = msg.button.ButtonReleased
    case ButtonReleaseMsg:
        msg.button.action()
        time.Sleep(time.Millisecond * 100)
        msg.button.style = msg.button.styleSelected
    case ButtonSelectMsg:
        msg.button.style = msg.button.styleSelected
    case ButtonUnselectMsg:
        msg.button.style = msg.button.styleUnselected
    }

    m.SetContents(m.View())

    return m, cmd
}

func (m Model) Init() tea.Cmd {
    return nil
}

