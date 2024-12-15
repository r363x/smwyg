package button

import (
    "time"
    "log"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)


type MsgType int

const (
    ButtonPress MsgType = iota
    ButtonRelease
    ButtonSelect
    ButtonUnselect
)

type Msg struct {
    Type MsgType
    Button *Model
}

type Model struct {
    label           string
    style           gloss.Style
    stylePressed    gloss.Style
    styleSelected   gloss.Style
    styleUnselected gloss.Style
    action          func()
}

func ButtonPressed(b *Model) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: ButtonPress,
            Button: b,
        }
    }
}

func ButtonReleased(b *Model) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: ButtonRelease,
            Button: b,
        }
    }
}

func ButtonSelected(b *Model) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: ButtonSelect,
            Button: b,
        }
    }
}

func ButtonUnselected(b *Model) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: ButtonUnselect,
            Button: b,
        }
    }
}

func (b Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            cmd = ButtonPressed(&b)
        }
    case Msg:
        switch msg.Type {

        case ButtonPress:
            msg.Button.style = msg.Button.stylePressed
            cmd = ButtonReleased(&b)

        case ButtonRelease:
            msg.Button.action()
            time.Sleep(time.Millisecond * 100)
            msg.Button.style = msg.Button.styleSelected

        case ButtonSelect:
            msg.Button.style = msg.Button.styleSelected

        case ButtonUnselect:
            msg.Button.style = msg.Button.styleUnselected

        }
    }

    return b, cmd
}

func (b Model) Init() tea.Cmd {
    return nil
}

func (b Model) View() string {
    return b.style.Render(b.label)
}

func (b *Model) Focus() tea.Cmd {
    return ButtonSelected(b)
}

func (b *Model) Blur() {
    b.style = b.styleUnselected
}

func New(label string) Model {

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

        return Model{
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
