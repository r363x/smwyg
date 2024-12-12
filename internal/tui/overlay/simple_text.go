package overlay

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SimpleTextOverlay struct {
    ModelBase
}

func New() SimpleTextOverlay {
    return SimpleTextOverlay{NewBase()}
}

func (o SimpleTextOverlay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return o, nil
}

func (o SimpleTextOverlay) Init() tea.Cmd {
    return nil
}

