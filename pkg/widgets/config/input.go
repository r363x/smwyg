package input

import (
    tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
)

var (
    FocusedStyle = gloss.NewStyle().Foreground(gloss.Color("205"))
    NoStyle = gloss.NewStyle()
)


type Model struct {
    textinput.Model
    Label string
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    model, cmd := m.Model.Update(msg)
    m.Model = model

    return m, cmd
}

func (m *Model) Focus() tea.Cmd {
    m.Model.PromptStyle = FocusedStyle
    m.Model.TextStyle = FocusedStyle
    m.Model.Cursor.Style = FocusedStyle
    
    return m.Model.Focus()
}

func (m *Model) Blur() {
    m.Model.PromptStyle = NoStyle
    m.Model.TextStyle = NoStyle
    m.Model.Cursor.Style = NoStyle
    m.Model.Blur()
}

func New(label string) Model {
    return Model{textinput.New(), label}
}

