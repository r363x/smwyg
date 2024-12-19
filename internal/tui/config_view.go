package tui

import (
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/button"
    "github.com/r363x/dbmanager/pkg/widgets/input"

    tea "github.com/charmbracelet/bubbletea"
)

func createConfigView() config.Model {

    var (
        views = []config.View{{Name: "Connect"}}
        inLabels = []string{"Type", "Host", "Port", "User", "Password", "DB Name"}
        btnLabels = []string{"Connect", "Close"}
    )

    elements := make([]config.Element, len(inLabels))

    for i, label := range inLabels {

        element := input.New(label)
        switch i {
        case 0:
            element.Focus()
        default:
            element.Blur()
        }

        elements[i] = &element
    }

    for _, label := range btnLabels {

        btn := button.New(label)

        switch label {
        case "Connect":
            btn.SetAction(func() tea.Msg {
                return config.Msg{Type: config.Submit}
            })

        case "Close":
            btn.SetAction(func() tea.Msg {
                return config.Msg{Type: config.Close}
            })
        }

        elements = append(elements, btn)
    }

    views[0].Elements = elements
    m := config.New(views)

    return m
}

