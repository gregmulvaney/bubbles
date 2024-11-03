package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gregmulvaney/bubbles/table"
)

type Model struct {
	width, height int
	table         table.Model
}

func New() (m Model) {
	cols := []table.Column{
		{Title: "id", Hidden: true},
		{Title: "#", Width: 3},
		{Title: "City", Flex: true},
		{Title: "Country", Flex: true},
	}

	rows := []table.Row{
		{"1", "1", "Tokyo", "Japan"},
		{"2", "2", "Los Angeles", "USA"},
		{"3", "3", "London", "Great Britain"},
		{"4", "4", "Warsaw", "Poland"},
	}

	m.table = table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocus(true),
	)

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyCtrlC.String():
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	m.table.UpdateDimensions(m.width, m.height)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return m.table.View()
}

func main() {
	p := tea.NewProgram(New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
