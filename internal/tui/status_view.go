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


func (t *tab) refreshStatusLeft() tea.Msg {
    msg := "Server: "

    version, err := t.dbManager.GetVersion()
    if err != nil {
        msg += "N/A"
    } else {
        msg += version
    }
    msg += "   "

    msg += "Status: "

    if err := t.dbManager.Status(); err != nil {
        msg += fmt.Sprintf("Error: %s", err)
    } else {
        msg += "Connected"
    }

    return statusMsg{
        section: secLeft,
        message: msg,
    }
}

func (t *tab) refreshStatusCenter() tea.Msg {
    return statusMsg{
        section: secCenter,
        message: "CENTER",
    }
}

func (t *tab) refreshStatusRight() tea.Msg {
    msg := "User: " + t.dbManager.DbUser()
    return statusMsg{
        section: secRight,
        message: msg,
    }
}

