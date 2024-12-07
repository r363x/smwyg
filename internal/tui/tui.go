package tui

import (
    "time"

	"github.com/charmbracelet/bubbles/textinput"
    gloss "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/bubbles/table"
    "github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss/tree"
	"github.com/r363x/dbmanager/internal/db"
)

type dimensions struct {
    width  int
    height int
}

type model struct {
	dbManager  db.Manager
    dbView     *tree.Tree
    queryView  textarea.Model
    resultView table.Model
    statusView statusView
    dimensions dimensions
}

type tickMsg time.Time

func doTick() tea.Cmd {
    return tea.Tick(time.Second * 3, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func New(dbManager db.Manager) (*tea.Program, error) {
	m := model{
		dbManager: dbManager,
        dbView: tree.New(),
        resultView: table.New(),
		queryView: textarea.New(),
	}

    m.queryView.FocusedStyle.Base = m.queryView.FocusedStyle.Base.
        BorderStyle(gloss.NormalBorder())
    m.queryView.BlurredStyle.Base = m.queryView.BlurredStyle.Base.
        BorderStyle(gloss.NormalBorder()).
        BorderForeground(gloss.Color("240"))

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
    if err == nil {
        m.refreshDbView()
    }

	m.queryView.Placeholder = "Enter your SQL query here"
    m.queryView.Focus()

	return tea.NewProgram(m), nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
        textinput.Blink,
        m.refreshStatusLeft,
        m.refreshStatusCenter,
        m.refreshStatusRight,
        doTick(),
    )
}

