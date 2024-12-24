package dropdown

import (

    tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
)

var (
    style = gloss.NewStyle().
        Align(gloss.Left).
        PaddingLeft(2).
        PaddingRight(2).
        MarginLeft(1).
        MarginRight(1)
    styleFocused = style.Background(gloss.Color("#8544b8"))
    styleBlurred = style.Background(gloss.Color("#7e7880"))
)

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

func NewItem(label string, defaults map[string]string) Item {

    return Item{
        Label: label,
        Defaults: defaults,
        style: style,
    }
}

