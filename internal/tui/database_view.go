package tui

import (
    "github.com/charmbracelet/lipgloss/tree"
)


func (m *model) refreshDbView() {

    m.dbView = tree.New()

    dbs, cur, err := m.dbManager.GetDatabases()
    if err != nil {
        m.dbView.Root("N/A")
        return
    }

    m.dbView = m.dbView.Root(".")

    tables, err := m.dbManager.GetTables(); if err != nil {
        return
    }

    for _, db := range dbs {
        m.dbView = m.dbView.Child(db)
        if db == cur {
            m.dbView = m.dbView.Child(tree.New().Child(tables))
        }
    }
}

