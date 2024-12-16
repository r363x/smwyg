package button

import (
	tea "github.com/charmbracelet/bubbletea"
)


func ButtonPress() tea.Msg {
    return Msg{Type: ButtonPressed}
}

func ButtonRelease() tea.Msg {
    return Msg{Type: ButtonReleased}
}

func ButtonSelect() tea.Msg {
    return Msg{Type: ButtonSelected}
}

func ButtonUnselect() tea.Msg {
    return Msg{Type: ButtonUnselected}
}

