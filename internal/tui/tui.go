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

type tab struct {
	dbManager  db.Manager
    dbView     *tree.Tree
    queryView  textarea.Model
    resultView table.Model
    statusView statusView
}

type overlayMenu struct {
    content string
    focused bool
}

type model struct {
    tabs []tab
    cur int
    overlay overlayMenu
    dimensions dimensions
}

type tickMsg time.Time

func doTick() tea.Cmd {
    return tea.Tick(time.Second * 3, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func New(dbManager db.Manager) (*tea.Program, error) {
    var m model
	tab := tab {
        dbManager: dbManager,
        dbView: tree.New(),
        resultView: table.New(),
        queryView: textarea.New(),
    }

    tab.queryView.FocusedStyle.Base = tab.queryView.FocusedStyle.Base.
        BorderStyle(gloss.NormalBorder())
    tab.queryView.BlurredStyle.Base = tab.queryView.BlurredStyle.Base.
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
    tab.resultView.SetStyles(s)

    // Connect to the database
    err := tab.dbManager.Connect()
    if err == nil {
        tab.refreshDbView()
    }

	tab.queryView.Placeholder = "Enter your SQL query here"
    tab.queryView.Focus()

    m.tabs = append(m.tabs, tab)
    m.cur = 0
    m.overlay.content = "Hello there" 
    m.overlay.focused = false

	return tea.NewProgram(m), nil
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
        textinput.Blink,
        m.tabs[m.cur].refreshStatusLeft,
        m.tabs[m.cur].refreshStatusCenter,
        m.tabs[m.cur].refreshStatusRight,
        doTick(),
    )
}

