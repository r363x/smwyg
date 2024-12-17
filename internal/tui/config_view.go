package tui

import (
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/button"
)

func createConfigView() config.Model {

    var (
        views = []config.View{{Name: "Connect"}}
        inLabels = []string{"Type", "Host", "Port", "User", "Password", "DB Name"}
        btnLabels = []string{"Close", "Save"}
    )

    elements := make([]config.Element, len(inLabels))

    for i, label := range inLabels {

        element := config.NewInput(label)
        switch i {
        case 0:
            element.Focus()
        default:
            element.Blur()
        }

        elements[i] = &element
    }

    for _, label := range btnLabels {
        elements = append(elements, button.New(label))
    }

    views[0].Elements = elements
    m := config.New(views)

    return m
}

