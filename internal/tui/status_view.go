package tui

import (
    "fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type statusDetails struct {
    left   string
    center string
    right  string
}

type statusView struct {
    content statusDetails
}

type statusMsg struct {
    section int
    message string
}

const (
    secLeft  = iota
    secCenter
    secRight
)


func (m *model) refreshStatusLeft() tea.Msg {
    msg := "Server: "

    version, err := m.dbManager.GetVersion()
    if err != nil {
        msg += "N/A"
    } else {
        msg += version
    }
    msg += "   "

    msg += "Status: "

    if err := m.dbManager.Status(); err != nil {
        msg += fmt.Sprintf("Error: %s", err)
    } else {
        msg += "Connected"
    }

    return statusMsg{
        section: secLeft,
        message: msg,
    }
}

func (m *model) refreshStatusCenter() tea.Msg {
    return statusMsg{
        section: secCenter,
        message: "CENTER",
    }
}

func (m *model) refreshStatusRight() tea.Msg {
    msg := "User: " + m.dbManager.DbUser()
    return statusMsg{
        section: secRight,
        message: msg,
    }
}

