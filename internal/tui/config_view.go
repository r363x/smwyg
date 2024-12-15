package tui

import (
    "github.com/r363x/dbmanager/pkg/widgets/overlay"
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/button"

    "github.com/charmbracelet/bubbles/textinput"
    gloss "github.com/charmbracelet/lipgloss"
)


func createConfigView() config.Model {

    // inStyle := gloss.NewStyle().
    //     Border(gloss.NormalBorder()).
    //     Width(30)
    // inStyle := struct{
    //     prompt gloss.Style
    //     text   gloss.Style
    // }{
    //     prompt: gloss.NewStyle().Border(gloss.NormalBorder()).Width(30),
    //     text: gloss.NewStyle().Foreground(gloss.Color("#b8a625")),
    // }

    focusedStyle := gloss.NewStyle().
        Foreground(gloss.Color("205")).
        Border(gloss.NormalBorder()).
        Width(30)

	noStyle := gloss.NewStyle().
        Border(gloss.NormalBorder()).
        Width(30)

    inType := textinput.New()
    inType.PromptStyle = focusedStyle
    inType.TextStyle = focusedStyle
    inType.Cursor.Style = focusedStyle

    inHost := textinput.New()
    inHost.PromptStyle = noStyle
    inHost.TextStyle = noStyle

    inPort := inHost
    inUser := inHost
    inPassword := inHost
    inDBName := inHost

    title := gloss.NewStyle().
        Align(gloss.Center).
        Height(3).
        Width(overlay.DefaultWidth-2)
    inputs := gloss.NewStyle().
        Align(gloss.Left).
        Width(100).
        Width(overlay.DefaultWidth-2)
    buttons := gloss.NewStyle().
        Align(gloss.Center).
        Width(overlay.DefaultWidth-2)

    btnClose := button.New("Close")

    content := gloss.JoinVertical(0,
        title.Render("Connect to database"),
        inputs.Render(
            gloss.JoinVertical(0,
                gloss.JoinHorizontal(0, "\nType:     ", inType.View()),
                gloss.JoinHorizontal(0, "\nHost:     ", inHost.View()),
                gloss.JoinHorizontal(0, "\nPort:     ", inPort.View()),
                gloss.JoinHorizontal(0, "\nUser:     ", inUser.View()),
                gloss.JoinHorizontal(0, "\nPassword: ", inPassword.View()),
                gloss.JoinHorizontal(0, "\nDB Name:  ", inDBName.View()),
                "\n\n\n",
            ),
        ),
        buttons.Render(btnClose.View()),
    )

    view := config.View{
        Name: "Connect",
        Content: content,
    }

    view.Selectables = append(view.Selectables,
        &inType,
        &inHost,
        &inPort,
        &inUser,
        &inPassword,
        &inDBName,
        &btnClose,
    )

    m := config.New([]config.View{view})
    m.FocusOn(0)
    return m
}
