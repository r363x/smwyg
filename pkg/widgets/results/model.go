package results

import (
    "fmt"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
    // gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/table"
)

type Model struct {
    Table table.Model
}

type MsgType int

const (
    DataChange MsgType = iota
)

type Msg struct {
    Type MsgType
    Data []map[string]interface{}
}

func UpdateResults(data []map[string]interface{}) tea.Cmd {
    return func() tea.Msg {
        return Msg{
            Type: DataChange,
            Data: data,
        }
    }
}

func New() Model {

    placeholderData := make([]map[string]interface{}, 1)

    // for i := range len(placeholderData) {
    //     row := make(map[string]interface{})
    //     for range 4 {
    //         row[strings.Repeat(" ", 10)] = strings.Repeat(" ", 10)
    //     }
    //     placeholderData[i] = row
    // }

    m := Model{Table: table.New()}
    m.fillTable(placeholderData)

    return m
}

func (m *Model) fillTable(data []map[string]interface{}) {

    columns := make([]table.Column, 0)

    _columns := make([]string, 0)

    // Set headers
    for col := range data[0] {
        _columns = append(_columns, col)
        columns = append(columns, table.Column{Title: strings.ToUpper(col), Width: m.Table.Width() / len(data[0])})
    }
    m.Table.SetRows(nil)
    m.Table.SetColumns(columns)

    // Set rows
    rows := make([]table.Row, 0)
    for _, item := range data {
        row := make([]string, 0)
        for _, key := range _columns {

            switch val := item[key].(type) {
            case int64:
                row = append(row, fmt.Sprintf("%d", val))
            case []byte:
                row = append(row, string(val))
            case string:
                row = append(row, val)
            }

        }
        rows = append(rows, row)
    }

    m.Table.SetRows(rows)
}

func (m Model) View() string {
    return m.Table.View()
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        m.Table, cmd = m.Table.Update(msg)

    case Msg:
        switch msg.Type {
        case DataChange:
            m.fillTable(msg.Data)
        }
    }

    return m, cmd
}

func (m *Model) Blur() {
    m.Table.Blur()
}

func (m *Model) Focus() tea.Cmd {
    m.Table.Focus()
    return nil
}
