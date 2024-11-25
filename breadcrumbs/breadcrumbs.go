package breadcrumbs

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	crumbs []string

	dimensions dimensions
	styles     Styles
}

type dimensions struct {
	width  int
	height int
}

type Styles struct {
	Crumb  lipgloss.Style
	Active lipgloss.Style
}

func DefaultStyles() Styles {
	baseStyle := lipgloss.NewStyle().Padding(0, 1).Background(lipgloss.Color("4")).Foreground(lipgloss.Color("0")).Bold(true).MarginRight(1)
	return Styles{
		Crumb:  baseStyle,
		Active: baseStyle.Background(lipgloss.Color("214")),
	}
}

func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

type Option func(*Model)

func New(opts ...Option) (m Model) {
	m.styles = DefaultStyles()
	m.dimensions = dimensions{
		width:  10,
		height: 0,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithCrumbs(c []string) Option {
	return func(m *Model) {
		m.crumbs = c
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var crumbs = make([]string, len(m.crumbs))
	for i, crumb := range m.crumbs {
		if i < len(m.crumbs)-1 {
			renderedCrumb := m.styles.Crumb.Render(crumb)
			crumbs[i] = renderedCrumb
		} else {
			renderedCrumb := m.styles.Active.Render(crumb)
			crumbs[i] = renderedCrumb
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, crumbs...)
}
