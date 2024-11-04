package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/keylist"
)

type Model struct {
	keylist keylist.Model
}

func New() (m Model) {
	items := []keylist.Item{
		{Key: "City:", Value: "Tokyo"},
		{Key: "Country:", Value: "Japan"},
		{Key: "Continent:", Value: "Asia"},
	}

	m.keylist = keylist.New(
		keylist.WithItems(items),
	)
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	return m.keylist.View()
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
