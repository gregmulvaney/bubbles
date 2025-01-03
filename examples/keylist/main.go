package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gregmulvaney/bubbles/keylist"
)

type Model struct {
	basicList keylist.Model
	gridList  keylist.Model
}

func New() (m Model) {
	keyData := [][]string{
		{"Red", "Orange"},
		{"Yellow", "Green"},
		{"Blue", "Indigo"},
		{"Violet", "Purple"},
		{"Fuschia", "Magenta"},
	}

	m.basicList = keylist.New(
		keylist.WithItems(keyData),
		keylist.WithSeparator(":"),
	)

	m.gridList = keylist.New(
		keylist.WithItems(keyData),
		keylist.WithGrid(true),
		keylist.WithSeparator(":"),
		keylist.WithMaxRows(3),
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
	return lipgloss.JoinVertical(lipgloss.Top, m.basicList.View(), " ", m.gridList.View())
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
