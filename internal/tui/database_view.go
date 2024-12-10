package tui

import (
    "fmt"
    "log"

    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/lipgloss/tree"
)

var (
    styleServer = gloss.NewStyle().Bold(true)
    styleCurDB = gloss.NewStyle().Bold(true)
)


func (m *model) refreshDbView() {

    dbs, cur, err := m.dbManager.GetDatabases()
    if err != nil {
        m.dbView.Root("N/A")
        return
    }

    m.dbView = tree.Root(styleServer.Render(
        fmt.Sprintf("  %s (%s)", m.dbManager.DbType(), m.dbManager.DbAddr()))).
            Enumerator(tree.RoundedEnumerator)

    tables, err := m.dbManager.GetTables(); if err != nil {
        return
    }

    for _, db := range dbs {
        if db == cur {

            // Database
            dbRoot := tree.Root(styleCurDB.Render("* " + db + "  ←"))
            var tableRoot *tree.Tree

            for _, table := range tables {

                // Table
                tableRoot = tree.Root(table)

                ts, err := m.dbManager.GetTableStructure(table, db)
                if err != nil {
                    log.Printf("Error: %s", err)
                    return
                }
                longestCol := 0
                for _, col := range ts.Columns {
                    if len(col.Name) > longestCol {
                        longestCol = len(col.Name)
                    }
                }
                for _, col := range ts.Columns {
                    tableRoot.Child(fmt.Sprintf("%-*s (%s)", longestCol + 1, col.Name, col.DataType))
                }
                dbRoot.Child(tableRoot)
            }
            m.dbView.Child(dbRoot)
            continue
        }

        // Database
        m.dbView = m.dbView.Child(db)
    }
}

