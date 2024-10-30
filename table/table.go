package table

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cols []Column
	rows []Row

	dimensions
}

type Row struct {
	Id   string
	Body []string
}

type Column struct {
	Title string
	Width *int
	Flex  *bool
}

type Styles struct {
	Header   lipgloss.Style
	Cell     lipgloss.Style
	Selected lipgloss.Style
}

type dimensions struct {
	width  int
	height int
}

type Option func(*Model)

func New(opts ...Option) Model {
	m := Model{}

	for _, opt := range opts {
		opt(&m)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	return m, nil
}

func (m Model) View() string {
	return fmt.Sprintf("H:%d W%d", m.dimensions.height, m.dimensions.width)
}

func (m *Model) UpdateDimensions(width int, height int) {

}

// func (m *Model) renderHeaderColumns() []string {
// 	renderedColumns := make([]string, len(m.cols))
// 	// Width of table consumed by static columns
// 	populatedWidth := 0
// 	// Count of columns with Flex enabled
// 	flexColumnCount := 0

// 	for i, column := range m.cols {
// 		if column.Flex {
// 			flexColumnCount += 1
// 			continue
// 		}

// 		if column.Width != 0 {
// 			renderedColumns[i] = lipgloss.NewStyle().Width(column.Width).MaxWidth(column.Width).Render(column.Title)
// 			populatedWidth += column.Width
// 		}
// 	}

// 	availableWidth := m.dimensions.width - populatedWidth
// 	flexCellWidth := availableWidth / flexColumnCount
// 	for i, column := range m.cols {
// 		if column.Flex == false {
// 			continue
// 		}

// 		renderedColumns[i] += lipgloss.NewStyle().Width(flexCellWidth).MaxWidth(flexCellWidth).Render(column.Title)
// 	}

// 	return renderedColumns
// }

// func (m *Model) renderHeader() string {
// 	headerColumns := m.renderHeaderColumns()
// 	header := lipgloss.JoinHorizontal(lipgloss.Left, headerColumns...)
// 	return lipgloss.NewStyle().Width(m.dimensions.width).MaxWidth(m.dimensions.width).Render(header)
// }

// func WithColumns(cols []Column) Option {
// 	return func(m *Model) {
// 		m.cols = cols
// 	}
// }

// func WithRows(rows []Row) Option {
// 	return func(m *Model) {
// 		m.rows = rows
// 	}
// }
