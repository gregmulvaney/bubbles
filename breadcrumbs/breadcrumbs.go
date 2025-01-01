package breadcrumbs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	styles Styles

	crumbs []string
}

type Styles struct {
	Crumb  lipgloss.Style
	Active lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Crumb:  lipgloss.NewStyle().Padding(0, 1).MarginRight(1).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("4")).Bold(true),
		Active: lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("214")).Bold(true),
	}
}

func (m *Model) SetCrumbs(crumbs []string) {
	m.crumbs = crumbs
}

type Option func(*Model)

func New(options ...Option) (m Model) {
	for _, opt := range options {
		opt(&m)
	}

	return m
}

func WithCrumbs(crumbs []string) Option {
	return func(m *Model) {
		m.crumbs = crumbs
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var crumbs []string
	for i, crumb := range m.crumbs {
		if i == len(m.crumbs)-1 {
			crumbs = append(crumbs, m.styles.Active.Render(crumb))
			continue
		}
		crumbs = append(crumbs, m.styles.Crumb.Render(crumb))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, crumbs...)
}
