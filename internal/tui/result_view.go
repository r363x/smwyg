package tui

import (
    "fmt"
    "strings"

    "github.com/charmbracelet/bubbles/table"
)

func (m *model) buildResultsTable(query string) error {
    data, err := m.dbManager.ExecuteQuery(query); if err != nil {
        return err
    }

    columns := make([]table.Column, 0)

    _columns := make([]string, 0)

    // Set headers
    for col := range data[0] {
        _columns = append(_columns, col)
        columns = append(columns, table.Column{Title: strings.ToUpper(col), Width: m.resultView.Width() / len(data[0])})
    }
    m.resultView.SetRows(nil)
    m.resultView.SetColumns(columns)

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
            }

        }
        rows = append(rows, row)
    }

    m.resultView.SetRows(rows)

    return nil
}
