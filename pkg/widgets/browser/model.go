package browser

import (
    "fmt"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/tree"
)

type Model struct {
    Tree    *tree.Tree
    Data    RefreshData
    cur     int
    focused bool
}

var (
    styleServer = gloss.NewStyle().Bold(true)
    styleCurDB = gloss.NewStyle().Bold(true)
    styleFocused = gloss.NewStyle().Foreground(gloss.Color("#8544b8"))
    styleBlurred = gloss.NewStyle()
)

func New() Model {

    return Model{
        Tree: tree.New(),
        focused: false,
    }
}

type NodeType int

const (
    Server NodeType = iota
    Database
    Table
    Column
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    var cmd tea.Cmd

    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.Type {
        case tea.KeyDown:
            //TODO

        case tea.KeyUp:
            //TODO

        case tea.KeyEnter:
            //TODO
        }

    case Msg:
        switch msg.Type {
        case RefreshResponse:
            m.Data = msg.Data
            m.RefreshTree()
        }
    }

    return m, cmd
}

func (m Model) View() string {
    return m.Tree.String()
}

func (m *Model) Focus() tea.Cmd {
    m.focused = true
    return nil
}

func (m *Model) Blur() {
    m.focused = false
}

func findDb(data *RefreshData, name string) int {
    var idx int

    for i := range data.Databases {
        if data.Databases[i].Name == name {
            idx = i
        }
    }
    return idx
}

func findLongestCol(table *TableData) int {
    longest := 0
    for _, col := range table.Columns {
        if len(col.Name) > longest {
            longest = len(col.Name)
        }
    }
    return longest
}

func (m *Model) RefreshTree() {

	m.Tree = tree.New().
        Root(styleServer.Render(
            fmt.Sprintf("  %s (%s)", m.Data.ServerType, m.Data.ServerAddr))).
                Enumerator(tree.DefaultEnumerator)

    for _, db := range m.Data.Databases {
        if db.Name == m.Data.CurDB {

            // Database
            dbRoot := tree.Root(styleCurDB.Render("* " + db.Name + "  ←"))
            var tableRoot *tree.Tree

            idx := findDb(&m.Data, m.Data.CurDB)

            for _, table := range m.Data.Databases[idx].Tables {

                // Table
                tableRoot = tree.Root(table.Name)

                longestCol := findLongestCol(&table)
                for _, col := range table.Columns {
                    tableRoot.Child(fmt.Sprintf("%-*s (%s)", longestCol + 1, col.Name, col.DataType))
                }
                dbRoot.Child(tableRoot)
            }
            m.Tree.Child(dbRoot)
            continue
        }

        // Database
        m.Tree = m.Tree.Child(db.Name)
    }
}

