package config

import (
    "fmt"
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	"github.com/r363x/dbmanager/pkg/widgets/button"

	tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
)

type Element interface {
    Focus() tea.Cmd
    Blur()
}

type View struct {
    Name        string
    Elements    []Element
    cur    int
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

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd

    view := &m.views[m.cur]

    switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
        case tea.KeyDown:
            if view.cur < len(view.Elements) - 1 {
                view.Elements[view.cur].Blur()
                view.cur++
                cmd = view.Elements[view.cur].Focus()
            }

        case tea.KeyUp:
            if view.cur >= 1 {
                view.Elements[view.cur].Blur()
                view.cur--
                cmd = view.Elements[view.cur].Focus()
            }

        default:
            cmd = m.updateElements(msg)
        }

    case button.Msg:
        btn := *view.Elements[view.cur].(*button.Model)
        btn, cmd = btn.Update(msg)
        view.Elements[view.cur] = &btn
    }

    return m, cmd
}

func (m *Model) updateElements(msg tea.Msg) tea.Cmd {

    elements := m.views[m.cur].Elements

    cmd := make([]tea.Cmd, len(elements))

    for i, element := range elements {

        switch element := element.(type) {
        case *Input:
            *element, cmd[i] = element.Update(msg)
        case *button.Model:
            *element, cmd[i] = element.Update(msg)
        }
    }

    return tea.Batch(cmd...)
}

func (m Model) View() string {

    var (
        box         = gloss.NewStyle().Width(overlay.DefaultWidth-2).Align(gloss.Center)
        inputBorder = gloss.NewStyle().Border(gloss.NormalBorder()).Width(30)
        titleBox    = box.Height(3)
        inputsBox   = box.Align(gloss.Left).PaddingLeft(4).MarginBottom(3)
        buttonsBox  = box

        inputs  []string
        buttons []string
    )

    for _, element := range m.views[m.cur].Elements {

        switch element := element.(type) {
        case *Input:
            inputs = append(inputs, gloss.JoinHorizontal(0,
                fmt.Sprintf("\n%-10s", element.label + ": "),
                inputBorder.Render(element.View(),
            )))

        case *button.Model:
            buttons = append(buttons, element.View())
        }
    }

    content := gloss.JoinVertical(0,
        titleBox.Render("Connect to database"),
        inputsBox.Render(gloss.JoinVertical(0, inputs...)),
        buttonsBox.Render(gloss.JoinHorizontal(0, buttons...)),
    )

    return m.ModelBase.View(content)
}

