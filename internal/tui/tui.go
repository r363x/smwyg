package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/r363x/dbmanager/internal/db"
)

type model struct {
	dbManager db.Manager
	query     textinput.Model
	results   string
}

func New(dbManager db.Manager) (*tea.Program, error) {
	m := model{
		dbManager: dbManager,
		query:     textinput.New(),
		results:   "",
	}
	m.query.Placeholder = "Enter your SQL query here"
	m.query.Focus()

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
		case tea.KeyEnter:
			results, err := m.dbManager.ExecuteQuery(m.query.Value())
			if err != nil {
				m.results = fmt.Sprintf("Error: %v", err)
			} else {
				m.results = fmt.Sprintf("%v", results)
			}
			return m, nil
		}
	}

	m.query, cmd = m.query.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"Enter your query:\n\n%s\n\nResults:\n\n%s\n\n(press ctrl+c to quit)",
		m.query.View(),
		m.results,
	)
}
