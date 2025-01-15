package browser

import (
	tea "github.com/charmbracelet/bubbletea"
)

type MsgType int

const (
    RefreshRequest MsgType = iota
    RefreshResponse
)

type ColumnData struct {
    Name string
    DataType string
}

type TableData struct {
    Name string
    Columns []ColumnData
}

type DBData struct {
    Name string
    Tables []TableData
}

type RefreshData struct {
    ServerType string
    ServerAddr string
    Databases  []DBData
    CurDB      string
}

type Msg struct {
    Type MsgType
    Data RefreshData
}

func requestRefresh() tea.Msg {
    return Msg{Type: RefreshRequest}
}

