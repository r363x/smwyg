package tui

import (
    "fmt"
    "log"

    "github.com/charmbracelet/lipgloss/tree"
)


func (m *model) refreshDbView() {

    dbs, cur, err := m.dbManager.GetDatabases()
    if err != nil {
        m.dbView.Root("N/A")
        return
    }

    m.dbView = tree.Root(m.dbManager.DbType()).
        Enumerator(tree.RoundedEnumerator)

    tables, err := m.dbManager.GetTables(); if err != nil {
        return
    }

    for _, db := range dbs {
        if db == cur {

            // Database
            dbRoot := tree.Root(db)
            var tableRoot *tree.Tree

            for _, table := range tables {

                // Table
                tableRoot = tree.Root(table)

                ts, err := m.dbManager.GetTableStructure(table, db)
                if err != nil {
                    log.Printf("Error: %s", err)
                    return
                }
                for _, col := range ts.Columns {
                    tableRoot.Child(fmt.Sprintf("%-8s (%s)", col.Name, col.DataType))
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

