package tui

import (
    "strings"

	"github.com/r363x/dbmanager/internal/tui/overlay/config"
)


func createConfigView() config.Model {

    btnClose := config.NewButton("Close")

    content := new(strings.Builder)

    content.WriteString("Connect to database\n\n\n")
    content.WriteString("INPUT 1: _____\n")
    content.WriteString("INPUT 2: _____\n")
    content.WriteString("INPUT 3: _____\n")
    content.WriteString(btnClose.View())
    content.WriteString("\n\n")

    view := config.View{
        Name: "Connect",
        Content: content.String(),
        Buttons: []config.Button{btnClose},
    }

    return config.New([]config.View{view})

}
