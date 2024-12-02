package tui

import (
	"fmt"
    "math"
    "strings"
    "log"

	"github.com/charmbracelet/bubbles/textinput"
    gloss "github.com/charmbracelet/lipgloss"
    // table "github.com/charmbracelet/lipgloss/table"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r363x/dbmanager/internal/db"
)

var (
    paneStyle = gloss.NewStyle().Border(gloss.NormalBorder())
    // HeaderStyle = gloss.NewStyle().Align(gloss.Center, gloss.Center).Bold(true)
    // EvenRowStyle = gloss.NewStyle().Align(gloss.Left, gloss.Center).Faint(true)
    // OddRowStyle = gloss.NewStyle().Align(gloss.Left, gloss.Center)
)

type dimensions struct {
    width  int
    height int
}

type dbView struct {
    content []string
    cursor  int
    focused bool
}

type model struct {
	dbManager  db.Manager
    dbView     dbView
    queryView  textarea.Model
    resultView table.Model
    dimensions dimensions
}

func (m *model) setResults(data interface{}) {
    m.resultView.SetRows([]table.Row{})
    switch data.(type) {

    case []map[string]interface{}:
        items := data.([]map[string]interface{})
        columns := make([]table.Column, 0)

        // Set headers
        for header := range items[0] {
            columns = append(columns, table.Column{Title: strings.ToUpper(header)})
        }
        m.resultView.SetColumns(columns)

        // Set rows
        rows := make([]table.Row, 0)
        for _, item := range items {
            row := make([]string, 0)
            for _, val := range item {
                switch v := val.(type) {
                case int64:
                    row = append(row, fmt.Sprintf("%d", v))
                case []byte:
                    row = append(row, string(v))
                }
            }
            rows = append(rows, row)
        }

        m.resultView.SetRows(rows)
    }
}


func (m *model) refreshDbView() {

    tables, err := m.dbManager.GetTables(); if err != nil {
        m.dbView.content[0] = fmt.Sprintf("Error: %v", err)
        return
    }

    m.dbView.content = tables
}

func New(dbManager db.Manager) (*tea.Program, error) {
	m := model{
		dbManager: dbManager,
        dbView: dbView{
            cursor: 0,
            focused: false,
        },
        resultView: table.New(),
		queryView: textarea.New(),
	}

    // Set table styles
    s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(gloss.NormalBorder()).
		BorderForeground(gloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(gloss.Color("229")).
		Background(gloss.Color("57")).
		Bold(false)
	m.resultView.SetStyles(s)

    // Connect to the database
    err := m.dbManager.Connect()
    if err != nil {
        fmt.Println(err)
    }

    m.refreshDbView()
	m.queryView.Placeholder = "Enter your SQL query here"
	m.queryView.Focus()

	return tea.NewProgram(m), nil
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

        case tea.KeyCtrlV:
            m.resultView.Blur()
            m.queryView.Blur()
            m.dbView.focused = true

        case tea.KeyCtrlQ:
            m.resultView.Blur()
            m.dbView.focused = false
            m.queryView.Focus()

        case tea.KeyCtrlR:
            m.queryView.Blur()
            m.dbView.focused = false
            m.resultView.Focus()

		case tea.KeyEnter:
            if m.queryView.Value() != "" && m.queryView.Focused() {
                results, err := m.dbManager.ExecuteQuery(m.queryView.Value())
                if err != nil {
                    log.Fatalln("Fatal: ", err)
                } else {
                    m.setResults(results)
                }
            }

            m.refreshDbView()
			return m, nil
		}
    case tea.WindowSizeMsg:
        m.dimensions.width = msg.Width
        m.dimensions.height = msg.Height
        return m, nil
	}

	m.queryView, cmd = m.queryView.Update(msg)
	return m, cmd
}

func (m model) View() string {

    // Left side, narrow
    paneDBView := paneStyle.
        Width(int(math.Ceil(float64(m.dimensions.width) * 0.1))).
        Height(int(math.Ceil(float64(m.dimensions.height) * 0.95)))

    paneQueryEditor := paneStyle.
        Width(int(math.Ceil(float64(m.dimensions.width) * 0.88))).
        Height(int(math.Ceil(float64(m.dimensions.height) * 0.61)))

    paneQueryResults := paneStyle.
        Width(int(math.Ceil(float64(m.dimensions.width) * 0.88))).
        Height(int(math.Ceil(float64(m.dimensions.height) * 0.31)))
        

    return gloss.JoinHorizontal(0.1, paneDBView.Render(m.dbView.content...), gloss.JoinVertical(0.1, paneQueryEditor.Render(m.queryView.View()), paneQueryResults.Render(m.resultView.View())))

}
