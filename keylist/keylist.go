package keylist

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	styles Styles

	maxRows   int
	grid      bool
	data      [][]string
	separator string
}

type Styles struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

// TODO: Should this be shared
func DefaultStyles() Styles {
	return Styles{
		Key:   lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).PaddingRight(0),
		Value: lipgloss.NewStyle(),
	}
}

type Option func(*Model)

func New(opts ...Option) (m Model) {
	m.maxRows = 0
	m.grid = false
	m.data = make([][]string, 0)
	m.separator = ""

	m.styles = DefaultStyles()

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithItems(items [][]string) Option {
	return func(m *Model) {
		m.data = items
	}
}

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

func WithMaxRows(i int) Option {
	return func(m *Model) {
		m.maxRows = i
	}
}

func WithGrid(b bool) Option {
	return func(m *Model) {
		m.grid = b
	}
}

func WithSeparator(s string) Option {
	return func(m *Model) {
		m.separator = s
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var items []string
	keyWidth := 0
	var keylist string
	separatorWidth := runewidth.StringWidth(m.separator) + 1

	for _, item := range m.data {
		kw := runewidth.StringWidth(item[0])
		if kw > keyWidth-1 {
			keyWidth = kw + 1 + separatorWidth
		}
	}

	if m.maxRows > 0 {
	} else {
		for _, item := range m.data {
			if !m.grid {
				keyWidth = runewidth.StringWidth(item[0]) + 1 + separatorWidth
			}

			key := m.styles.Key.Width(keyWidth).Render(fmt.Sprintf("%s%s", item[0], m.separator))
			value := m.styles.Value.Render(item[1])
			items = append(items, lipgloss.JoinHorizontal(lipgloss.Left, key, value))
		}
		keylist = lipgloss.JoinVertical(lipgloss.Top, items...)
	}

	return keylist
}
