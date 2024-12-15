package config

import (
    "time"
    "log"
	"github.com/r363x/dbmanager/pkg/widgets/overlay"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
)

type ButtonMsg Msg

type MsgType int

const (
    ButtonPress MsgType = iota
    ButtonRelease
    ButtonSelect
    ButtonUnselect
)

type Msg struct {
    Type MsgType
    button *Button
}

type Selectable interface {
    Focus() tea.Cmd
    Blur()
}

type Button struct {
    label           string
    style           gloss.Style
    stylePressed    gloss.Style
    styleSelected   gloss.Style
    styleUnselected gloss.Style
    action          func()
}

func ButtonPressed(b *Button) tea.Cmd {
    return func() tea.Msg {
        return ButtonMsg{
            Type: ButtonPress,
            button: b,
        }
    }
}

func ButtonReleased(b *Button) tea.Cmd {
    return func() tea.Msg {
        return ButtonMsg{
            Type: ButtonRelease,
            button: b,
        }
    }
}

func ButtonSelected(b *Button) tea.Cmd {
    return func() tea.Msg {
        return ButtonMsg{
            Type: ButtonSelect,
            button: b,
        }
    }
}

func ButtonUnselected(b *Button) tea.Cmd {
    return func() tea.Msg {
        return ButtonMsg{
            Type: ButtonUnselect,
            button: b,
        }
    }
}

func (b Button) Update(msg tea.Msg) (Button, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyEnter:
            cmd = ButtonPressed(&b)
        }
    case ButtonMsg:
        switch msg.Type {

        case ButtonPress:
            msg.button.style = msg.button.stylePressed
            cmd = ButtonReleased(&b)

        case ButtonRelease:
            msg.button.action()
            time.Sleep(time.Millisecond * 100)
            msg.button.style = msg.button.styleSelected

        case ButtonSelect:
            msg.button.style = msg.button.styleSelected

        case ButtonUnselect:
            msg.button.style = msg.button.styleUnselected

        }
    }

    return b, cmd
}

func (b Button) Init() tea.Cmd {
    return nil
}

func (b Button) View() string {
    return b.style.Render(b.label)
}

func (b *Button) Focus() tea.Cmd {
    return ButtonSelected(b)
}

func (b *Button) Blur() {
    b.style = b.styleUnselected
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
            case *Button:
                updated, _cmd := sel.Update(msg)
                m.views[m.cur].Selectables[cur] = &updated
                cmd = _cmd
            }
        }
    case ButtonMsg:
        btn, _cmd := msg.button.Update(msg)
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

