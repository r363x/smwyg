package dropdown

import (
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

func close() tea.Msg {
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
    Label    string
    Defaults map[string]string
    style           gloss.Style
    stylePressed    gloss.Style
    styleFocused    gloss.Style
    styleBlurred    gloss.Style
}

type Model struct {
    Items []Item
    cur   int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:

        switch msg.Type {
        case tea.KeyDown:
            //TODO
        case tea.KeyUp:
            //TODO
        case tea.KeyEnter:
            
        }

    case Msg:

        switch msg.Type {
        case Opened:
            //TODO
        case Closed:
            //TODO
        case Selected:
            cmd = DeliverData(msg.Data)
        }

    }

    return m, cmd
}

func (m Model) View() string {
    //TODO
}

func (m *Item) Focus() tea.Cmd {
    m.style = m.styleFocused
    return nil

}

func (m *Item) Blur() {
    m.style = m.styleBlurred
}

func New(items []Item) *Model {

    m := Model{cur: 0}

    empty := NewItem("---", map[string]string{})

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
        stylePressed: stylePressed,
        styleFocused: styleFocused,
        styleBlurred: styleBlurred,
    }
}

