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
		Key:   lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
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

func (m *Model) SetItems(items [][]string) {
	m.data = items
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var keySets [][]string
	if m.maxRows > 0 {
		for i := 0; i < len(m.data); i += m.maxRows {
			end := i + m.maxRows
			if end > len(m.data) {
				end = len(m.data)
			}
			set := m.renderKeySet(m.data[i:end])
			keySets = append(keySets, set)
		}
	} else {
		keySets = append(keySets, m.renderKeySet(m.data))
	}
	var renderedKeyList []string
	for _, set := range keySets {
		list := lipgloss.NewStyle().PaddingRight(1).Render(lipgloss.JoinVertical(lipgloss.Top, set...))
		renderedKeyList = append(renderedKeyList, list)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, renderedKeyList...)
}

func (m *Model) renderKeySet(items [][]string) []string {
	var keySet []string
	sepWidth := 0
	if m.separator != "" {
		sepWidth = runewidth.StringWidth(m.separator)
	}
	keyWidth := 0
	for _, item := range items {
		kw := runewidth.StringWidth(item[0]) + sepWidth + 1
		if kw > keyWidth-1 {
			keyWidth = kw
		}
	}

	for _, item := range items {
		if !m.grid {
			keyWidth = runewidth.StringWidth(item[0]) + sepWidth + 1
		}

		if m.separator != "" {
			key := m.styles.Key.Width(keyWidth).Render(fmt.Sprintf("%s%s", item[0], m.separator))
			value := m.styles.Value.Render(item[1])
			keySet = append(keySet, lipgloss.JoinHorizontal(lipgloss.Left, key, value))
		} else {
			key := m.styles.Key.Width(keyWidth).Render(item[0])
			value := m.styles.Value.Render(item[1])
			keySet = append(keySet, lipgloss.JoinHorizontal(lipgloss.Left, key, value))
		}
	}

	return keySet
}
