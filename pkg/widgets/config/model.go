package config

import (
    "fmt"
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	"github.com/r363x/dbmanager/pkg/widgets/button"

	tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textinput"
)

type Element interface {
    Focus() tea.Cmd
    Blur()
}

type View struct {
    Name        string
    Elements    []Element
    InLabels    []string
    selected    int
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

    switch msg := msg.(type) {
	case tea.KeyMsg:

        view := &m.views[m.cur]

		switch msg.Type {
        case tea.KeyDown:
            if view.selected < len(view.Elements) - 1 {
                view.Elements[view.selected].Blur()
                view.selected++
                cmd = view.Elements[view.selected].Focus()
            }

        case tea.KeyUp:
            if view.selected >= 1 {
                view.Elements[view.selected].Blur()
                view.selected--
                cmd = view.Elements[view.selected].Focus()
            }

        default:
            cmd = m.updateElements(msg)
        }
    }

    return m, cmd
}

func (m *Model) updateElements(msg tea.Msg) tea.Cmd {

    elements := m.views[m.cur].Elements

    cmd := make([]tea.Cmd, len(elements))

    for i, element := range elements {

        switch element := element.(type) {
        case *textinput.Model:
            *element, cmd[i] = element.Update(msg)
        case *button.Model:
            *element, cmd[i] = element.Update(msg)
        }
    }

    return tea.Batch(cmd...)
}

func (m Model) View() string {

    var (
        box         = gloss.NewStyle().Width(overlay.DefaultWidth-2)
        inputBorder = gloss.NewStyle().Border(gloss.NormalBorder()).Width(30)
        titleBox    = box.Align(gloss.Center).Height(3)
        inputsBox   = box.Align(gloss.Left)
        buttonsBox  = box.Align(gloss.Center)
        inputs      []string
        buttons     []string
    )


    for i, element := range m.views[m.cur].Elements {

        switch element := element.(type) {
        case *textinput.Model:

            inputs = append(inputs, gloss.JoinHorizontal(0,
                fmt.Sprintf("\n%-10s", m.views[m.cur].InLabels[i] + ": "),
                inputBorder.Render(element.View(),
            )))

        case *button.Model:
            buttons = append(buttons, element.View())
        }
    }

    // Cosmetics
    inputs = append(inputs, "\n\n\n")

    content := gloss.JoinVertical(0,
        titleBox.Render("Connect to database"),
        inputsBox.Render(gloss.JoinVertical(0, inputs...)),
        buttonsBox.Render(gloss.JoinHorizontal(0.6, buttons...)),
    )

    return m.ModelBase.View(content)
}

func (m Model) Init() tea.Cmd {
    return textinput.Blink
}

