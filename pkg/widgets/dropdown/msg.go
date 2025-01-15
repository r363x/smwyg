package dropdown

import (
    tea "github.com/charmbracelet/bubbletea"
)

type MsgType int

const (
    Opened MsgType = iota
    Closed
    Selected
    SelectionData
)

type Msg struct {
    Type MsgType
    Data map[string]string
}

func Open() tea.Msg {
    return Msg{Type: Opened}
}

func Close() tea.Msg {
    return Msg{Type: Closed}
}

func Select() tea.Msg {
    return Msg{Type: Selected}
}

func DeliverData(data map[string]string) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: SelectionData,
            Data: data,
        }
    }
}

