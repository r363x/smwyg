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


func (t *tab) refreshDbView() {

    dbs, cur, err := t.dbManager.GetDatabases()
    if err != nil {
        t.dbView.Root("N/A")
        return
    }

    t.dbView = tree.Root(styleServer.Render(
        fmt.Sprintf("  %s (%s)", t.dbManager.DbType(), t.dbManager.DbAddr()))).
            Enumerator(tree.RoundedEnumerator)

    tables, err := t.dbManager.GetTables(); if err != nil {
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

                ts, err := t.dbManager.GetTableStructure(table, db)
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
            t.dbView.Child(dbRoot)
            continue
        }

        // Database
        t.dbView = t.dbView.Child(db)
    }
}

