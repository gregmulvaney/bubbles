package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/breadcrumbs"
)

type Model struct {
	width, height int
	breadcrumbs   breadcrumbs.Model
}

func New() (m Model) {
	crumbs := []string{
		"Home",
		"Page 2",
		"Page 3",
	}

	m.breadcrumbs = breadcrumbs.New(
		breadcrumbs.WithCrumbs(crumbs),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.breadcrumbs.View()
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
