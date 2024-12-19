package config

import (
    "fmt"
    "log"
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	"github.com/r363x/dbmanager/pkg/widgets/button"
	"github.com/r363x/dbmanager/pkg/widgets/input"

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

type MsgType int

type Msg struct {
    Type MsgType
    Data interface{}
}

const (
    Close MsgType = iota
    Submit
    FormData
)

func DeliverFormData(data map[string]string) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: FormData,
            Data: data,
        }
    }
}

func (m *Model) Close() {
    view := m.views[m.cur]

    view.Elements[view.cur].Blur()
    view.cur = 0
    view.Elements[view.cur].Focus()
    m.Show = false
}

func New(views []View) Model {
    base := overlay.NewBase()
    return Model{ base, views, 0 }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var (
        cmds []tea.Cmd
        cmd  tea.Cmd
    )

    view := &m.views[m.cur]

    switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
        case tea.KeyDown:
            if view.cur < len(view.Elements) - 1 {
                view.Elements[view.cur].Blur()
                view.cur++
                cmd = view.Elements[view.cur].Focus()
                cmds = append(cmds, cmd)
            }

        case tea.KeyTab:

            view.Elements[view.cur].Blur()

            switch n := len(view.Elements) - 1; {
            case view.cur < n:
                view.cur++
            case view.cur == n:
                view.cur = 0
            }

            cmd = view.Elements[view.cur].Focus()
            cmds = append(cmds, cmd)


        case tea.KeyUp:
            if view.cur >= 1 {
                view.Elements[view.cur].Blur()
                view.cur--
                cmd = view.Elements[view.cur].Focus()
                cmds = append(cmds, cmd)
            }

        case tea.KeyShiftTab:

            view.Elements[view.cur].Blur()

            switch n := len(view.Elements) - 1; {
            case view.cur >= 1:
                view.cur--
            case view.cur == 0:
                view.cur = n
            }

            cmd = view.Elements[view.cur].Focus()
            cmds = append(cmds, cmd)

        default:
            cmd = m.updateElements(msg)
            cmds = append(cmds, cmd)
        }

    case button.Msg:
        btn := *view.Elements[view.cur].(*button.Model)
        btn, cmd = btn.Update(msg)
        view.Elements[view.cur] = &btn
        cmds = append(cmds, cmd)

    case Msg:
        switch msg.Type {
        case Close:
            log.Print("Close")
            m.Close()

        case Submit:

            data := make(map[string]string)

            for _, element := range m.views[m.cur].Elements {

                switch element := element.(type) {
                case *input.Model:
                    data[element.Label] = element.Value()
                }
            }

            // Deliver the map with form data
            cmds = append(cmds, DeliverFormData(data))
            m.Close()
        }
    }


    return m, tea.Batch(cmds...)
}

func (m *Model) updateElements(msg tea.Msg) tea.Cmd {

    elements := m.views[m.cur].Elements

    cmd := make([]tea.Cmd, len(elements))

    for i, element := range elements {

        switch element := element.(type) {
        case *input.Model:
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
        case *input.Model:
            inputs = append(inputs, gloss.JoinHorizontal(0,
                fmt.Sprintf("\n%-10s", element.Label + ": "),
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

