package breadcrumbs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	styles Styles

	height int
	width  int

	crumbs []Crumb
}

type Crumb struct {
	Value  string
	Active bool
}

type Styles struct {
	Crumb lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Crumb: lipgloss.NewStyle().Padding(0, 1).Margin(0, 1).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("4")),
	}
}

func (m *Model) SetCrumbs(crumbs []Crumb) {
	m.crumbs = crumbs
}

func (m *Model) SetHeight(h int) {
	m.height = h
}

func (m *Model) SetWidth(w int) {
	m.width = w
}

type Option func(*Model)

func New(options ...Option) (m Model) {
	m.styles = DefaultStyles()
	m.height = 0
	m.width = 0

	for _, opt := range options {
		opt(&m)
	}

	return m
}

func WithCrumbs(crumbs []Crumb) Option {
	return func(m *Model) {
		m.crumbs = crumbs
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var crumbs []string
	for _, crumb := range m.crumbs {
		crumbs = append(crumbs, m.styles.Crumb.Render(crumb.Value))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, crumbs...)
}
