package tab

import (
    "fmt"
    "log"
    "strconv"
    cfg "github.com/r363x/dbmanager/internal/config"
    "github.com/r363x/dbmanager/pkg/widgets/config"
    "github.com/r363x/dbmanager/pkg/widgets/status"
    "github.com/r363x/dbmanager/pkg/widgets/results"
    "github.com/r363x/dbmanager/pkg/widgets/browser"

    tea "github.com/charmbracelet/bubbletea"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/textarea"
    "github.com/charmbracelet/lipgloss/tree"
    "github.com/r363x/dbmanager/internal/db"

)

type Element interface {
    Focus() tea.Cmd
    Blur()
}

// type MsgType int

// const (
//     Focus MsgType = iota
//     Blur
//     Refocus
// )

// type Msg struct {
//     Type MsgType
//     Data interface{}
// }

type dimensions struct {
    width  int
    height int
}

// TODO: unexport in favor of a proper New()
type Model struct {
	DbManager  db.Manager
    Elements   []Element
    cur        int

    StatusView status.Model
    dimensions dimensions
    lastFocus  int
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

    var (
        cmds       []tea.Cmd
        idxQuery   int
        idxResults int
        idxBrowser int
    )

    for i := range m.Elements {
        switch m.Elements[i].(type) {
        case *textarea.Model:
            idxQuery = i
        case *results.Model:
            idxResults = i
        case *browser.Model:
            idxBrowser = i
        }
    }

    switch msg := msg.(type) {
    case tea.KeyMsg:

        switch msg.Type {

        case tea.KeyCtrlQ:
            m.BlurAll()
            cmds = append(cmds, m.Elements[idxQuery].Focus())
            m.cur = idxQuery

        case tea.KeyCtrlR:
            m.BlurAll()
            cmds = append(cmds, m.Elements[idxResults].Focus())
            m.cur = idxResults

        case tea.KeyCtrlB:
            m.BlurAll()
            cmds = append(cmds, m.Elements[idxBrowser].Focus())
            m.cur = idxBrowser

        case tea.KeyF5:
            if element, ok := m.Elements[m.cur].(*textarea.Model); ok {
                if query := element.Value(); len(query) > 0 {

                    if m.DbManager != nil {
                        data, err := m.DbManager.ExecuteQuery(query); if err != nil {
                            return m, m.RefreshStatusCenter(err.Error())
                        }
                        return m, tea.Batch(
                            results.UpdateResults(data),
                            m.RefreshStatusCenter("Query: OK"),
                        )
                    } else {
                        return m, m.RefreshStatusCenter("Query: cannot execute query (no connection)")
                    }
                }
            }
            cmds = append(cmds, m.RefreshBrowser)

        default:
            switch element := m.Elements[m.cur].(type) {
            case *textarea.Model:
                e, cmd := element.Update(msg)
                m.Elements[m.cur] = &e
                cmds = append(cmds, cmd)

            case *results.Model:
                e, cmd := element.Update(msg)
                m.Elements[m.cur] = &e
                cmds = append(cmds, cmd)

            case *browser.Model:
                e, cmd := element.Update(msg)
                m.Elements[m.cur] = &e
                cmds = append(cmds, cmd)
            }
        }

    case config.Msg:

        switch msg.Type {
        case config.FormData:

            c := cfg.Config{}

            if data, ok := msg.Data.(map[string]string); ok {

                for key, value := range data {

                    switch key {
                    case "Type":
                        c.DatabaseConfig.Type = value
                    case "Host":
                        c.DatabaseConfig.Host = value
                    case "Port":
                        port, err := strconv.Atoi(value)
                        if err != nil {
                            log.Fatal("Error: ", err)
                        }
                        c.DatabaseConfig.Port = port
                    case "User":
                        c.DatabaseConfig.User = value
                    case "Password":
                        c.DatabaseConfig.Password = value
                    case "DB Name":
                        c.DatabaseConfig.DBName = value
                    }
                    log.Print(value)
                }


                dbManager, err := db.NewManager(c.DatabaseConfig)
                if err != nil {
                    cmds = append(cmds, m.RefreshStatusCenter(fmt.Sprintf("Error: %v", err)))
                }
                m.DbManager = dbManager
                err = m.DbManager.Connect()
                if err == nil {
                    m.RefreshDbView()
                }
            }



        // case Focus:
        //     cmd = append(cmd, msg.Element.Focus())

        // case Blur:
        //     msg.Element.Blur()

        // case Refocus:
        //     cmd = append(cmd, m.Elements[m.lastFocus].Focus())
        }

    case status.Msg:
        m.UpdateStatus(msg)

    case browser.Msg:
        updated, cmd := m.Elements[idxBrowser].(*browser.Model).Update(msg)
        m.Elements[idxBrowser] = &updated
        cmds = append(cmds, cmd)

    case results.Msg:
        updated, cmd := m.Elements[idxResults].(*results.Model).Update(msg)
        m.Elements[idxResults] = &updated
        cmds = append(cmds, cmd)
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
        ptrBrowserWidget *browser.Model
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

        case *browser.Model:
            ptrBrowserWidget = element
        }
    }

    return gloss.JoinVertical(0,
        gloss.JoinHorizontal(0,
            paneDBView.Render(ptrBrowserWidget.View()),
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

func (m *Model) RefreshBrowser() tea.Msg {

    data := browser.RefreshData{}

    dbs, cur, err := m.DbManager.GetDatabases()
    if err != nil {
        return m.RefreshStatusCenter(fmt.Sprintf("Error: %s", err.Error()))
    }

    data.CurDB = cur

    tables, err := m.DbManager.GetTables(); if err != nil {
        return m.RefreshStatusCenter(fmt.Sprintf("Error: %s", err.Error()))
    }

    for _, dbName := range dbs {
        db := browser.DBData{Name: dbName}

        for _, tableName := range tables {

            table := browser.TableData{Name: tableName}

            ts, err := m.DbManager.GetTableStructure(tableName, dbName)
            if err != nil {
                return m.RefreshStatusCenter(fmt.Sprintf("Error: %s", err.Error()))
            }

            for _, col := range ts.Columns {
                table.Columns = append(table.Columns, browser.ColumnData{
                    Name: col.Name,
                    DataType: col.DataType,
                })
            }
            db.Tables = append(db.Tables, table)
        }
        data.Databases = append(data.Databases, db)
    }

    return browser.Msg{
        Type: browser.RefreshResponse,
        Data: data,
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

    var (
        msgServer = "Server: "
        msgStatus = "Status: "
    )

    if m.DbManager == nil {
        msgServer += "N/A"
        msgStatus += "N/A"

    } else {
        version, err := m.DbManager.GetVersion()
        if err != nil {
            msgServer += "N/A"
        } else {
            msgServer += version
        }

        if err := m.DbManager.Status(); err != nil {
            msgStatus += err.Error()
        } else {
            msgStatus += "Connected"
        }
    }

    return status.Msg{
        Section: status.SecLeft,
        Message: fmt.Sprintf("%s   %s", msgServer, msgStatus),
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
    msg := "User: "

    if m.DbManager == nil {
        msg += "N/A"

    } else {

        if m.DbManager != nil {
            msg += m.DbManager.DbUser()
        } else {
            msg += "N/A"
        }
    }

    return status.Msg{
        Section: status.SecRight,
        Message: msg,
    }
}

func (m *Model) SetDimentions(width, height int) {
    m.dimensions.width = width
    m.dimensions.height = height
}

func (m *Model) BlurAll() {
    for _, element := range m.Elements {
        element.Blur()
    }
}

