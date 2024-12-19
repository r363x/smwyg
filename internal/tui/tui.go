package tui

import (
    "time"
    "github.com/r363x/dbmanager/internal/tui/tab"
    "github.com/r363x/dbmanager/pkg/widgets/results"
    "github.com/r363x/dbmanager/pkg/widgets/browser"

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

type tickMsg time.Time

func doTick(sec time.Duration) tea.Cmd {
    return tea.Tick(time.Second * sec, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func New(DbManager db.Manager) (*tea.Program, error) {
    var m model
	tab := tab.Model{DbManager: DbManager}

    queryWidget := textarea.New()
    resultsWidget := results.New()
    browserWidget := browser.New()

    // Set query area styles
    queryWidget.FocusedStyle.Base = queryWidget.FocusedStyle.Base.
        BorderStyle(gloss.NormalBorder())

    queryWidget.BlurredStyle.Base = queryWidget.BlurredStyle.Base.
        BorderStyle(gloss.NormalBorder()).
        BorderForeground(gloss.Color("240"))

    // Set results area styles
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

    resultsWidget.Table.SetStyles(s)

    // Connect to the database TODO: handle err
    tab.DbManager.Connect()

	queryWidget.Placeholder = "Enter your SQL query here"
    queryWidget.Focus()

    tab.Elements = append(tab.Elements, &queryWidget)
    tab.Elements = append(tab.Elements, &resultsWidget)
    tab.Elements = append(tab.Elements, &browserWidget)

    m.tabs = append(m.tabs, tab)
    m.cur = 0
    m.overlay = createConfigView()

	return tea.NewProgram(m), nil
}


