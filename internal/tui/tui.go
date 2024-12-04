package tui

import (
	"fmt"
    "log"

	"github.com/charmbracelet/bubbles/textinput"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r363x/dbmanager/internal/db"
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

type statusDetails struct {
    left   string
    center string
    right  string
}

type statusView struct {
    content statusDetails
}

type model struct {
	dbManager  db.Manager
    dbView     dbView
    queryView  textarea.Model
    resultView table.Model
    statusView statusView
    dimensions dimensions
}

func (m *model) refreshDbView() {

    tables, err := m.dbManager.GetTables(); if err != nil {
        m.dbView.content[0] = fmt.Sprintf("Error: %v", err)
        return
    }

    m.dbView.content = tables
}

func (m *model) refreshStatusView() {
    m.statusView.content.left = "LEFT"
    m.statusView.content.center = "CENTER"
    m.statusView.content.right = "RIGHT"
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
        Bold(true).
        AlignHorizontal(gloss.Center)
    s.Selected = s.Selected.
        Foreground(gloss.Color("229")).
        Background(gloss.Color("57")).
        Bold(true)
    m.resultView.SetStyles(s)

    // Connect to the database
    err := m.dbManager.Connect()
    if err != nil {
        fmt.Println(err)
    }

    m.refreshDbView()
	m.queryView.Placeholder = "Enter your SQL query here"
    m.queryView.Focus()
	// m.resultView.Focus()

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

        case tea.KeyF5:
            if m.queryView.Value() != "" && m.queryView.Focused() {
                err := m.dbManager.ExecuteQuery(m.queryView.Value(), &m.resultView, m.dimensions.width)
                if err != nil {
                    log.Fatalln("Fatal: ", err)
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

    switch {
    case m.queryView.Focused():
        m.queryView, cmd = m.queryView.Update(msg)
    case m.resultView.Focused():
        m.resultView, cmd = m.resultView.Update(msg)
    }
    m.refreshStatusView()

    return m, cmd
}

func (m model) View() string {

    // Left side, narrow
    paneDBView := gloss.NewStyle().
        Border(gloss.NormalBorder()).
        BorderForeground(gloss.Color("240")).
        Width(int(float64(m.dimensions.width) * 0.1)).
        Height(m.dimensions.height - 3)

    m.queryView.SetWidth(m.dimensions.width - paneDBView.GetWidth())
    m.queryView.SetHeight(int(float64(m.dimensions.height) * 0.5) - 3)

    m.resultView.SetWidth(m.dimensions.width - paneDBView.GetWidth())
    m.resultView.SetHeight(m.dimensions.height - m.queryView.Height() - 1)

    // Bottom, narrow
    paneStatusView := gloss.NewStyle().
        Width(m.dimensions.width).
        Height(1)
    statusItemView := gloss.NewStyle().
        Width(paneStatusView.GetWidth() / 3)

    return gloss.JoinVertical(0,
        gloss.JoinHorizontal(0,
            paneDBView.Render(m.dbView.content...),
            gloss.JoinVertical(0,
                m.queryView.View(),
                m.resultView.View(),
        )),
        paneStatusView.Render(gloss.JoinHorizontal(0,
            statusItemView.AlignHorizontal(gloss.Left).Render(m.statusView.content.left),
            statusItemView.AlignHorizontal(gloss.Center).Render(m.statusView.content.center),
            statusItemView.AlignHorizontal(gloss.Right).Render(m.statusView.content.right),
        )),
    )

}
