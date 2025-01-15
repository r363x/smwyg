package tui

import (
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/button"
    "github.com/r363x/dbmanager/pkg/widgets/input"
    "github.com/r363x/dbmanager/pkg/widgets/dropdown"

    tea "github.com/charmbracelet/bubbletea"
)

var (
    defaultsMysql = map[string]string{
        "host": "localhost",
        "port": "3306",
    }
    defaultsPostgres = map[string]string{
        "host": "localhost",
        "port": "5432",
    }
)

func createConfigView() config.Model {

    var (
        views     = []config.View{{Name: "Connect"}}
        inLabels  = []string{"Host", "Port", "User", "Password", "DB Name"}
        btnLabels = []string{"Connect", "Close"}
        elements    []config.Element
        items       []dropdown.Item
    )


    items = append(items, dropdown.NewItem("mysql", defaultsMysql))
    items = append(items, dropdown.NewItem("postgres", defaultsPostgres))

    dbTypes := dropdown.New(items, "Available drivers")
    dbTypes.Label = "Type"
    dbTypes.SetWidth(20)
    dbTypes.SetHeight(8)

    elements = append(elements, dbTypes)

    for i, label := range inLabels {

        element := input.New(label)
        switch i {
        case 0:
            element.Focus()
        default:
            element.Blur()
        }

        elements = append(elements, &element)
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

