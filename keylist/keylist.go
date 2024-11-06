package keylist

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	items     []Item
	maxHeight int
	keyBlock  bool

	styles Styles
}

type Item struct {
	Key   string
	Value string
}

type Styles struct {
	Key   lipgloss.Style
	Value lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Key:   lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("2")),
		Value: lipgloss.NewStyle().PaddingLeft(1),
	}
}

func (m *Model) SetValues(items []Item) {
	m.items = items
}

func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

type Option func(*Model)

func New(opts ...Option) (m Model) {
	m.styles = DefaultStyles()
	m.keyBlock = false

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func WithItems(items []Item) Option {
	return func(m *Model) {
		m.items = items
	}
}

func WithMaxHeight(h int) Option {
	return func(m *Model) {
		m.maxHeight = h
	}
}

func WithKeyWidth(k bool) Option {
	return func(m *Model) {
		m.keyBlock = k
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	items := m.renderItems()
	return items
}

func (m *Model) renderItems() string {
	renderedItems := make([]string, 0, len(m.items))
	keyWidth := 0
	for _, item := range m.items {
		keyWidth = max(keyWidth, len(item.Key))
	}

	for _, item := range m.items {
		key := m.styles.Key.Width(keyWidth).Render(item.Key)
		value := m.styles.Value.Render(item.Value)
		renderedItem := lipgloss.JoinHorizontal(lipgloss.Left, key, value)
		renderedItems = append(renderedItems, renderedItem)
	}

	return lipgloss.JoinVertical(lipgloss.Top, renderedItems...)
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}
