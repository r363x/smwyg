package config

import (
    tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
)

var (
    FocusedStyle = gloss.NewStyle().Foreground(gloss.Color("205"))
    NoStyle = gloss.NewStyle()
)


type Input struct {
    textinput.Model
    label string
}

func (m Input) Update(msg tea.Msg) (Input, tea.Cmd) {
    model, cmd := m.Model.Update(msg)
    m.Model = model

    return m, cmd
}

func (m *Input) Focus() tea.Cmd {
    m.Model.PromptStyle = FocusedStyle
    m.Model.TextStyle = FocusedStyle
    m.Model.Cursor.Style = FocusedStyle
    
    return m.Model.Focus()
}

func (m *Input) Blur() {
    m.Model.PromptStyle = NoStyle
    m.Model.TextStyle = NoStyle
    m.Model.Cursor.Style = NoStyle
    m.Model.Blur()
}

func NewInput(label string) Input {
    return Input{textinput.New(), label}
}

