package simpletext

import (
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
    overlay.ModelBase
}

func New() Model {
    return Model{overlay.NewBase()}
}

func (o Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    return o, nil
}

func (o Model) View() string {
    return o.BaseView()
}

func (o Model) Init() tea.Cmd {
    return nil
}

