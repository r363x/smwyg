package tui

import (
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/button"

    "github.com/charmbracelet/bubbles/textinput"
    gloss "github.com/charmbracelet/lipgloss"
)


func createConfigView() config.Model {

    var (
        focusedStyle = gloss.NewStyle().Foreground(gloss.Color("205"))
        noStyle = gloss.NewStyle()
    )

    var views = []config.View{{
        Name: "Connect",
        InLabels: []string{"Type", "Host", "Port", "User", "Password", "DB Name"},
    }}

    elements := make([]config.Element, 8)

    for i := range 6 {

        element := textinput.New()

        switch i {
        case 0:
            element.PromptStyle = focusedStyle
            element.TextStyle = focusedStyle
            element.Cursor.Style = focusedStyle

        default:
            element.PromptStyle = noStyle
            element.TextStyle = noStyle
            element.Cursor.Style = noStyle
        }

        elements[i] = &element
    }

    elements[6] = button.New("Close")
    elements[7] = button.New("Save")

    views[0].Elements = elements
    m := config.New(views)

    return m
}

