package config

import (
    "fmt"
    "log"
    "strings"
	"github.com/r363x/dbmanager/pkg/widgets/overlay"
	"github.com/r363x/dbmanager/pkg/widgets/button"
	"github.com/r363x/dbmanager/pkg/widgets/input"
	"github.com/r363x/dbmanager/pkg/widgets/dropdown"

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

        if m.Active() {

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


            case tea.KeyEnter:
                switch element := view.Elements[view.cur].(type) {
                case *dropdown.Model:
                    if element.Show {
                        *element, cmd = element.Update(msg)
                    } else {
                        m.BlurAll()
                        m.Deactivate()
                        element.Show = true
                    }

                case *button.Model:
                    *element, cmd = element.Update(msg)
                }

                cmds = append(cmds, cmd)

            default:
                cmd = m.updateElements(msg)
                cmds = append(cmds, cmd)
            }

        } else {

            switch element := view.Elements[view.cur].(type) {
            case *dropdown.Model:
                *element, cmd = element.Update(msg)
                cmds = append(cmds, cmd)
            }
        }

    case button.Msg:
        btn := *view.Elements[view.cur].(*button.Model)
        btn, cmd = btn.Update(msg)
        view.Elements[view.cur] = &btn
        cmds = append(cmds, cmd)

    case dropdown.Msg:
        switch msg.Type {
        case dropdown.SelectionData:
            m.setPlaceholders(msg.Data)
            m.Activate()

        default:
            dd := *view.Elements[view.cur].(*dropdown.Model)
            dd, cmd = dd.Update(msg)
            view.Elements[view.cur] = &dd
            cmds = append(cmds, cmd)
        }

    case Msg:
        switch msg.Type {
        case Close:
            log.Print("Close")
            m.Close()

        case Submit:

            data := make(map[string]string)

            for _, element := range m.views[m.cur].Elements {

                switch element := element.(type) {
                case *dropdown.Model:
                    data[element.Label] = element.Selection().Label
                case *input.Model:
                    value := element.Value()
                    if value == "" {
                        value = element.Placeholder
                    }
                    data[element.Label] = value
                }
            }

            // Deliver the map with form data
            cmds = append(cmds, DeliverFormData(data))
            log.Printf("%#v", data)
            m.Close()
        }
    }


    return m, tea.Batch(cmds...)
}

func (m *Model) BlurAll() {
    for _, element := range m.views[m.cur].Elements {
        element.Blur()
    }
}

func (m *Model) setPlaceholders(data map[string]string) {

    if data == nil {
        m.clearPlaceholders()
        return
    }

    elements := m.views[m.cur].Elements

    for i := range elements {

        switch element := elements[i].(type) {
        case *input.Model:

            for label, placeholder := range data {
                if strings.ToLower(element.Label) == strings.ToLower(label) {
                    element.Placeholder = placeholder
                }
            }
        }
    }
}

func (m *Model) clearPlaceholders() {
    elements := m.views[m.cur].Elements
    for i := range elements {
        switch element := elements[i].(type) {
        case *input.Model:
            element.Placeholder = ""
        }
    }
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
        case *dropdown.Model:
            *element, cmd[i] = element.Update(msg)
        }
    }

    return tea.Batch(cmd...)
}

func (m *Model) SetDimensions(width, height int) {
    m.ModelBase.SetDimensions(width, height)
    for _, element := range m.views[m.cur].Elements {
        if element, ok := element.(*dropdown.Model); ok {
            element.SetDimensions(width, height)
        }
    }
}

func (m Model) View() string {

    var (
        box         = gloss.NewStyle().Width(overlay.DefaultWidth-2).Align(gloss.Center)
        inputBorder = gloss.NewStyle().Border(gloss.NormalBorder()).Width(30)
        titleBox    = box.Height(3)
        inputsBox   = box.Align(gloss.Left).PaddingLeft(4).MarginBottom(3)
        buttonsBox  = box

        inputs      []string
        buttons     []string
        ptrDropdown *dropdown.Model
    )

    for _, element := range m.views[m.cur].Elements {

        switch element := element.(type) {
        case *dropdown.Model:
            // Save the state
            ptrDropdown = element
            state := ptrDropdown.Show

            // Hide when drawing the background
            ptrDropdown.Show = false

            inputs = append(inputs, gloss.JoinHorizontal(0,
                fmt.Sprintf("\n%-10s", element.Label + ": "),
                inputBorder.Render(element.View(),
            )))

            // Restore the state
            ptrDropdown.Show = state

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

    final := m.ModelBase.View(content)

    if ptrDropdown.Show {
        ptrDropdown.SetBackground(final)
        final = ptrDropdown.ModelBase.View(ptrDropdown.View())
    }

    return final
}

