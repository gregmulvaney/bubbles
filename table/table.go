package table

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type Model struct {
	keymap Keymap
	styles Styles

	// data
	cols []Column
	rows []Row

	// state
	focused bool
	width   int
	height  int

	// table
	cursor        int
	flexCellWidth int
	start         int
	end           int

	viewport viewport.Model
}

type Column struct {
	Title    string
	Width    int
	MinWidth int
	Flex     bool
	Hidden   bool
	Auto     bool
	Padding  int
}

type Row []string

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

// Update styles
func (m *Model) SetStyles(s Styles) {
	m.styles = s
}

// Update rows
func (m *Model) SetRows(r []Row) {
	m.rows = r
}

// Update width
func (m *Model) SetWidth(w int) {
	m.width = w
	m.viewport.Width = w
}

// Update height
func (m *Model) SetHeight(h int) {
	m.height = h
	m.viewport.Height = h - 4
}

type Option func(*Model)

func New(opts ...Option) (m Model) {
	// Set defaults
	m.keymap = DefaultKeymap()
	m.styles = DefaultStyles()
	m.cursor = 0
	m.flexCellWidth = 0
	m.height = 0
	m.width = 0
	m.viewport = viewport.New(0, 0)

	// Run all options
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
		m.focused = f
	}
}

func WithStyles(s Styles) Option {
	return func(m *Model) {
		m.styles = s
	}
}

func WithKeymap(k Keymap) Option {
	return func(m *Model) {
		m.keymap = k
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.LineDown):
			m.MoveDown(1)
		case key.Matches(msg, m.keymap.LineUp):
			m.MoveUp(1)
		}
	case tea.WindowSizeMsg:
		m.viewport.Height = msg.Height - 4
		m.viewport.Width = msg.Width
		m.UpdateViewport()
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) Focused() bool {
	return m.focused
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) View() string {
	header := m.renderHeader()
	return lipgloss.JoinVertical(lipgloss.Top, header, m.viewport.View())
}

func (m *Model) UpdateViewport() {
	renderedRows := make([]string, 0, len(m.rows))

	if m.cursor >= 0 {
		m.start = clamp(m.cursor-m.viewport.Height, 0, m.cursor)
	} else {
		m.cursor = 0
	}
	m.end = clamp(m.cursor+m.viewport.Height, 0, len(m.rows))

	m.viewport.Width = m.width
	m.viewport.Height = min(len(m.rows), m.height)

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
	consumedWidth := 0
	flexColumnCount := 0
	// autoCellWidth := 0

	for i, col := range m.cols {
		if col.Hidden {
			continue
		} else if col.Flex {
			flexColumnCount += 1
			continue
		} else if col.Width > 0 {
			consumedWidth += col.Width
			renderedColumns[i] = m.styles.Header.Width(col.Width).MaxWidth(col.Width).Render(col.Title)
		}
	}

	availableWidth := m.width - consumedWidth

	if flexColumnCount > 0 {
		flexCellWidth := availableWidth / flexColumnCount
		m.flexCellWidth = flexCellWidth
		for i, col := range m.cols {
			if !col.Flex || col.Hidden || col.Width > 0 {
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
	return lipgloss.NewStyle().Width(m.width).MaxWidth(m.width).Render(header)
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
		renderedCell := m.styles.Cell.Width(width).MaxWidth(width).Render(runewidth.Truncate(value, width-3, "..."))
		cells = append(cells, renderedCell)
	}

	row := lipgloss.JoinHorizontal(lipgloss.Left, cells...)

	if r == m.cursor {
		return m.styles.Selected.Width(m.width).MaxWidth(m.width).Render(row)
	}

	return lipgloss.NewStyle().Width(m.width).MaxWidth(m.width).Render(row)
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
