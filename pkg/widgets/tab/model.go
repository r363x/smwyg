package tab

import (
    "fmt"
    "log"
    "github.com/r363x/dbmanager/pkg/widgets/status"
    "github.com/r363x/dbmanager/pkg/widgets/results"

	tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/lipgloss/tree"
	"github.com/r363x/dbmanager/internal/db"

)

var (
    styleServer = gloss.NewStyle().Bold(true)
    styleCurDB = gloss.NewStyle().Bold(true)
)

type Element interface {
    Focus() tea.Cmd
    Blur()
}

type MsgType int

const (
    Focus MsgType = iota 
    Blur
    Refocus
)

type Msg struct {
    Type    MsgType
    Element Element
}

type dimensions struct {
    width int
    height int
}

// TODO: unexport in favor of a proper New()
type Model struct {
	DbManager  db.Manager
    Elements   []Element
    cur        int

    DbView     *tree.Tree
    StatusView status.Model
    dimensions dimensions
    lastFocus  int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var (
        cmd     []tea.Cmd
        ptrResults *results.Model
        ptrQuery   *textarea.Model
    )

    for i := range m.Elements {
        switch element := m.Elements[i].(type) {
        case *textarea.Model:
            ptrQuery = element
        case *results.Model:
            ptrResults = element
        }
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:

        switch msg.Type {

        case tea.KeyCtrlQ:
            for _, element := range m.Elements {
                element.Blur()
            }
            cmd = append(cmd, ptrQuery.Focus())

        case tea.KeyCtrlR:
            for _, element := range m.Elements {
                element.Blur()
            }
            cmd = append(cmd, ptrResults.Focus())


        case tea.KeyF5:
            if element, ok := m.Elements[m.cur].(*textarea.Model); ok {
                if query := element.Value(); len(query) > 0 {

                   data, err := m.DbManager.ExecuteQuery(query); if err != nil {
                       return m, m.RefreshStatusCenter(err.Error())
                   }
                   return m, tea.Batch(
                       results.UpdateResults(data),
                       m.RefreshStatusCenter("Query: OK"),
                   )
                }
            }
            m.RefreshDbView()
            return m, nil

        default:
            switch element := m.Elements[m.cur].(type) {
            case *textarea.Model:
                e, cmd := element.Update(msg)
                m.Elements[m.cur] = &e
                return m, cmd

            case *results.Model:
                e, cmd := element.Update(msg)
                m.Elements[m.cur] = &e
                return m, cmd
            }
        }

    case Msg:
        switch msg.Type {
        case Focus:
            cmd = append(cmd, msg.Element.Focus())

        case Blur:
            msg.Element.Blur()

        case Refocus:
            cmd = append(cmd, m.Elements[m.lastFocus].Focus())
        }
        return m, tea.Batch(cmd...)

    case status.Msg:
        m.UpdateStatus(msg)
        return m, nil

    case results.Msg:
        _ptr, _cmd := ptrResults.Update(msg)
        *ptrResults = _ptr
        return m, _cmd

    }


    return m, tea.Batch(cmd...)
}

func (m Model) View() string {

    var (
        paneDBView = gloss.NewStyle().
            Border(gloss.NormalBorder()).
            BorderForeground(gloss.Color("240")).
            Width(int(float64(m.dimensions.width) * 0.1)).
            Height(m.dimensions.height - 3)

        paneStatusView = gloss.NewStyle().
            Width(m.dimensions.width).
            Padding(0).
            Height(1)

        statusItemView = gloss.NewStyle().
            Padding(0).
            Width(paneStatusView.GetWidth() / 3)


        ptrQueryWidget   *textarea.Model
        ptrResultsWidget *results.Model
    )

    for i := range m.Elements {
        switch element := m.Elements[i].(type) {

        case *textarea.Model:
            element.SetWidth(m.dimensions.width - paneDBView.GetWidth() - 4)
            element.SetHeight(int(float64(m.dimensions.height) * 0.5) - 10)

            // Save the pointer
            ptrQueryWidget = element

        case *results.Model:
            newWidth := m.dimensions.width - paneDBView.GetWidth() 

            // Resize the table
            cols := element.Table.Columns()

            // Resize columns
            for i := range len(cols) {
                cols[i].Width = newWidth / len(cols)
            }

            element.Table.SetColumns(cols)
            element.Table.SetWidth(newWidth)
            element.Table.SetHeight(m.dimensions.height - ptrQueryWidget.Height() - 3)

            // Save the pointer
            ptrResultsWidget = element
        }
    }

    return gloss.JoinVertical(0,
        gloss.JoinHorizontal(0,
            paneDBView.Render(m.DbView.String()),
            gloss.JoinVertical(0,
                ptrQueryWidget.View(),
                ptrResultsWidget.View(),
        )),
        paneStatusView.Render(gloss.JoinHorizontal(0,
            statusItemView.AlignHorizontal(gloss.Left).Render(m.StatusView.Content.Left),
            statusItemView.AlignHorizontal(gloss.Center).Render(m.StatusView.Content.Center),
            statusItemView.AlignHorizontal(gloss.Right).Render(m.StatusView.Content.Right),
        )),
    )

}

func (m *Model) RefreshDbView() {

    dbs, cur, err := m.DbManager.GetDatabases()
    if err != nil {
        m.DbView.Root("N/A")
        return
    }

    m.DbView = tree.Root(styleServer.Render(
        fmt.Sprintf("  %s (%s)", m.DbManager.DbType(), m.DbManager.DbAddr()))).
            Enumerator(tree.RoundedEnumerator)

    tables, err := m.DbManager.GetTables(); if err != nil {
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

                ts, err := m.DbManager.GetTableStructure(table, db)
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
            m.DbView.Child(dbRoot)
            continue
        }

        // Database
        m.DbView = m.DbView.Child(db)
    }
}

func (m *Model) UpdateStatus(msg status.Msg) {
    switch msg.Section {
    case status.SecLeft:
        m.StatusView.Content.Left = msg.Message
    case status.SecCenter:
        m.StatusView.Content.Center = msg.Message
    case status.SecRight:
        m.StatusView.Content.Right = msg.Message
    }
}

func (m *Model) RefreshStatusLeft() tea.Msg {
    msg := "Server: "

    version, err := m.DbManager.GetVersion()
    if err != nil {
        msg += "N/A"
    } else {
        msg += version
    }
    msg += "   "

    msg += "Status: "

    if err := m.DbManager.Status(); err != nil {
        msg += fmt.Sprintf("Error: %s", err)
    } else {
        msg += "Connected"
    }

    return status.Msg{
        Section: status.SecLeft,
        Message: msg,
    }
}

func (m *Model) RefreshStatusCenter(text string) tea.Cmd {

    if text == "" {
        if m.StatusView.Content.Center == "" {
            text = "CENTER"
        } else {
            text = m.StatusView.Content.Center
        }
    }

    return func() tea.Msg {
        return status.Msg{
            Section: status.SecCenter,
            Message: text,
        }
    }
}

func (m *Model) RefreshStatusRight() tea.Msg {
    msg := "User: " + m.DbManager.DbUser()
    return status.Msg{
        Section: status.SecRight,
        Message: msg,
    }
}

func (m *Model) SetDimentions(width, height int) {
    m.dimensions.width = width
    m.dimensions.height = height
}

