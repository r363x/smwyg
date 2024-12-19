package button

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MsgType int

const (
    ButtonPressed MsgType = iota
    ButtonReleased
)

type Msg struct {
    Type MsgType
}

func ButtonPress() tea.Msg {
    return Msg{Type: ButtonPressed}
}

func ButtonRelease() tea.Msg {
    return Msg{Type: ButtonReleased}
}

