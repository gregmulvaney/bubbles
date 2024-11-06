package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	Keymap Keymap

	cols   []Column
	rows   []Row
	cursor int
	focus  bool
	styles Styles

	flexCellWidth int
	dimensions    dimensions

	viewport viewport.Model
	start    int
	end      int
}

type Column struct {
	Title  string
	Width  int
	Flex   bool
	Hidden bool
}

type Row []string

type dimensions struct {
	width  int
	height int
}

type Keymap struct {
	LineUp   key.Binding
	LineDown key.Binding
}

func DefaultKeymap() Keymap {
	return Keymap{
		LineUp: key.NewBinding(
			key.WithKeys("up", "k"),
		),
		LineDown: key.NewBinding(
			key.WithKeys("down", "j"),
		),
	}
}

type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
		Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("0")).Background(lipgloss.Color("4")),
	}
}

func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

func (m *Model) SetRows(r []Row) {
	m.rows = r
}

type Option func(*Model)

func New(opts ...Option) (m Model) {
	m.Keymap = DefaultKeymap()
	m.styles = DefaultStyles()
	m.cursor = 0

	m.dimensions = dimensions{
		width:  0,
		height: 0,
	}
	m.flexCellWidth = 0

	m.viewport = viewport.New(10, 20)

	for _, opt := range opts {
		opt(&m)
	}

	m.UpdateViewport()

	return m
}

func WithColumns(cols []Column) Option {
	return func(m *Model) {
		m.cols = cols
	}
}

func WithRows(rows []Row) Option {
	return func(m *Model) {
		m.rows = rows
	}
}

func WithFocus(f bool) Option {
	return func(m *Model) {
		m.focus = f
	}
}

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

func WithKeymap(k Keymap) Option {
	return func(m *Model) {
		m.Keymap = k
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keymap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.Keymap.LineUp):
			m.MoveUp(1)
		}
	}

	return m, nil
}

func (m *Model) Focused() bool {
	return m.focus
}

func (m *Model) Focus() {
	m.focus = true
	m.UpdateViewport()
}

func (m *Model) Blur() {
	m.focus = false
	m.UpdateViewport()
}

func (m Model) View() string {
	header := m.renderHeader()
	return lipgloss.JoinVertical(lipgloss.Top, header, m.viewport.View())
}

func (m *Model) UpdateDimensions(width int, height int) {
	m.dimensions.width = width
	m.dimensions.height = height
	m.UpdateViewport()
}

func (m *Model) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.rows))

	if m.cursor >= 0 {
		m.start = clamp(m.cursor-m.viewport.Height, 0, m.cursor)
	} else {
		m.cursor = 0
	}
	m.end = clamp(m.cursor+m.viewport.Height, 0, len(m.rows))

	m.viewport.Width = m.dimensions.width
	m.viewport.Height = min(len(m.rows), m.dimensions.height)

	for i := range m.rows {
		renderedRows = append(renderedRows, m.renderRow(i))
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, renderedRows...))
}

func (m Model) SelectedRow() Row {
	if m.cursor < 0 || m.cursor > len(m.rows) {
		return nil
	}

	return m.rows[m.cursor]
}

func (m *Model) MoveUp(n int) {
	m.cursor = clamp(m.cursor-n, 0, len(m.rows)-1)
	switch {
	case m.start == 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset, 0, m.cursor))
	case m.start < m.viewport.Height:
		m.viewport.YOffset = (clamp(clamp(m.viewport.YOffset+n, 0, m.cursor), 0, m.viewport.Height))
	case m.viewport.YOffset >= 1:
		m.viewport.YOffset = clamp(m.viewport.YOffset+n, 1, m.viewport.Height)
	}
	m.UpdateViewport()
}

func (m *Model) MoveDown(n int) {
	m.cursor = clamp(m.cursor+n, 0, len(m.rows)-1)
	m.UpdateViewport()

	switch {
	case m.end == len(m.rows) && m.viewport.YOffset > 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.viewport.Height))
	case m.cursor > (m.end-m.start)/2 && m.viewport.YOffset > 0:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset-n, 1, m.cursor))
	case m.viewport.YOffset > 1:
	case m.cursor > m.viewport.YOffset+m.viewport.Height-1:
		m.viewport.SetYOffset(clamp(m.viewport.YOffset+1, 0, 1))
	}
}

func (m *Model) renderHeaderColumns() []string {
	renderedColumns := make([]string, len(m.cols))
	// Width consumed by staticly sized elements
	populatedWidth := 0
	flexColCount := 0

	// Render all staticly sized columns first
	for i, col := range m.cols {
		if col.Hidden {
			continue
		}
		if col.Flex {
			flexColCount += 1
			continue
		}
		if col.Width != 0 {
			renderedColumns[i] = m.styles.Header.Width(col.Width).MaxWidth(col.Width).Render(col.Title)
		}
	}

	availableWidth := m.dimensions.width - populatedWidth

	if flexColCount > 0 {
		flexCellWidth := availableWidth / flexColCount
		m.flexCellWidth = flexCellWidth
		for i, col := range m.cols {
			if !col.Flex {
				continue
			}
			renderedColumns[i] = m.styles.Header.Width(flexCellWidth).MaxWidth(flexCellWidth).Render(col.Title)
		}
	}
	return renderedColumns
}

func (m *Model) renderHeader() string {
	columns := m.renderHeaderColumns()
	header := lipgloss.JoinHorizontal(lipgloss.Left, columns...)
	return lipgloss.NewStyle().Width(m.dimensions.width).MaxWidth(m.dimensions.width).Render(header)
}

func (m *Model) renderRow(r int) string {
	m.renderHeaderColumns()
	cells := make([]string, 0, len(m.cols))

	for i, value := range m.rows[r] {
		var width int
		if m.cols[i].Hidden {
			continue
		}
		if m.cols[i].Width != 0 {
			width = m.cols[i].Width
		}
		if m.cols[i].Flex {
			width = m.flexCellWidth
		}
		renderedCell := m.styles.Cell.Width(width).MaxWidth(width).Render(runewidth.Truncate(value, width, "..."))
		cells = append(cells, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)

	if r == m.cursor {
		return m.styles.Selected.Width(m.dimensions.width).MaxWidth(m.dimensions.width).Render(row)
	}

	return lipgloss.NewStyle().Width(m.dimensions.width).MaxWidth(m.dimensions.width).Render(row)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func clamp(v, low, high int) int {
	return min(max(v, low), high)
}
