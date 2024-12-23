package dropdown

import (
	"github.com/r363x/dbmanager/pkg/widgets/overlay"

    tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

type MsgType int

const (
    Opened MsgType = iota
    Closed
    Selected
    SelectionData
)

var (
    style = gloss.NewStyle().
        Align(gloss.Left).
        PaddingLeft(2).
        PaddingRight(2).
        MarginLeft(1).
        MarginRight(1)

    stylePressed = style.Background(gloss.Color("#c3ccdb"))

    styleFocused = style.Background(gloss.Color("#8544b8"))

    styleBlurred = style.Background(gloss.Color("#5a3478"))
)

type Msg struct {
    Type MsgType
    Data map[string]string
}

func Open() tea.Msg {
    return Msg{Type: Opened}
}

func Close() tea.Msg {
    return Msg{Type: Closed}
}

func Select() tea.Msg {
    return Msg{Type: Selected}
}

func DeliverData(data map[string]string) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: SelectionData,
            Data: data,
        }
    }
}

type Item struct {
    Label           string
    Defaults        map[string]string
    style           gloss.Style
}

func (m *Item) Focus() tea.Cmd {
    m.style = styleFocused
    return nil
}

func (m *Item) Blur() {
    m.style = styleBlurred
}

type Model struct {
    overlay.ModelBase
    Items []Item
    Label string
    cur   int
}

func (m *Model) Selection() Item {
    return m.Items[m.cur]
}

func (m *Model) Focus() tea.Cmd {
    return nil
}

func (m *Model) Blur() {
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
        var items []string

        for _, item := range m.Items {
            items = append(items, item.style.Render(item.Label))
        }

        return gloss.JoinVertical(0, items...)
    } else {
        item := m.Selection()
        return item.style.Render(item.Label)
    }
}

func New(items []Item) *Model {

    base := overlay.NewBase()

    m := Model{ModelBase: base, cur: 0}

    empty := NewItem("---", nil)

    m.Items = append(m.Items, empty)

    for i := range items {
        m.Items = append(m.Items, items[i])
    }

    return &m
}


func NewItem(label string, defaults map[string]string) Item {

    return Item{
        Label: label,
        Defaults: defaults,
        style: style,
    }
}

