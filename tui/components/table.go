package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// TableRowSelectedMsg is sent when a table row is selected
type TableRowSelectedMsg struct {
	Index int
	Row   table.Row
}

// GameTable wraps bubbles/table with game-specific styling
type GameTable struct {
	table    table.Model
	title    string
	width    int
	height   int
	focused  bool
}

// NewGameTable creates a new styled table
func NewGameTable(title string, columns []table.Column, rows []table.Row) *GameTable {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	// Apply custom styles
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Cyan).
		BorderBottom(true).
		Bold(true).
		Foreground(styles.Cyan)
	s.Selected = s.Selected.
		Foreground(styles.Black).
		Background(styles.Cyan).
		Bold(true)
	s.Cell = s.Cell.
		Padding(0, 1)

	t.SetStyles(s)

	return &GameTable{
		table:   t,
		title:   title,
		width:   60,
		height:  15,
		focused: true,
	}
}

// SetSize sets the table dimensions
func (t *GameTable) SetSize(width, height int) {
	t.width = width
	t.height = height
	t.table.SetWidth(width - 4)
	t.table.SetHeight(height - 4)
}

// SetFocused sets the focus state
func (t *GameTable) SetFocused(focused bool) {
	t.focused = focused
	t.table.Focus()
	if !focused {
		t.table.Blur()
	}
}

// SetRows updates the table rows
func (t *GameTable) SetRows(rows []table.Row) {
	t.table.SetRows(rows)
}

// SetColumns updates the table columns
func (t *GameTable) SetColumns(columns []table.Column) {
	t.table.SetColumns(columns)
}

// SelectedRow returns the currently selected row
func (t *GameTable) SelectedRow() table.Row {
	return t.table.SelectedRow()
}

// Cursor returns the cursor position
func (t *GameTable) Cursor() int {
	return t.table.Cursor()
}

// Init initializes the table
func (t *GameTable) Init() tea.Cmd {
	return nil
}

// Update handles input for the table
func (t *GameTable) Update(msg tea.Msg) (*GameTable, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Enter):
			if len(t.table.Rows()) > 0 {
				return t, func() tea.Msg {
					return TableRowSelectedMsg{
						Index: t.table.Cursor(),
						Row:   t.table.SelectedRow(),
					}
				}
			}
		}
	}

	t.table, cmd = t.table.Update(msg)
	return t, cmd
}

// View renders the table
func (t *GameTable) View() string {
	var b strings.Builder

	// Title
	if t.title != "" {
		titleStyle := styles.TitleStyle.Width(t.width).Align(lipgloss.Center)
		b.WriteString(titleStyle.Render(t.title))
		b.WriteString("\n\n")
	}

	// Box style based on focus
	boxStyle := styles.BoxStyle.Width(t.width)
	if t.focused {
		boxStyle = styles.FocusedBoxStyle.Width(t.width)
	}

	b.WriteString(boxStyle.Render(t.table.View()))

	return b.String()
}

// SimpleTable creates a simple table from string data
func SimpleTable(title string, headers []string, data [][]string, widths []int) *GameTable {
	columns := make([]table.Column, len(headers))
	for i, h := range headers {
		width := 15
		if i < len(widths) {
			width = widths[i]
		}
		columns[i] = table.Column{Title: h, Width: width}
	}

	rows := make([]table.Row, len(data))
	for i, d := range data {
		rows[i] = d
	}

	return NewGameTable(title, columns, rows)
}

// InvestmentTable creates a table specifically for displaying investments
func InvestmentTable(title string) *GameTable {
	columns := []table.Column{
		{Title: "Company", Width: 20},
		{Title: "Invested", Width: 12},
		{Title: "Value", Width: 12},
		{Title: "Equity", Width: 8},
		{Title: "P/L", Width: 12},
	}

	return NewGameTable(title, columns, []table.Row{})
}

// LeaderboardTable creates a table for displaying leaderboards
func LeaderboardTable(title string) *GameTable {
	columns := []table.Column{
		{Title: "#", Width: 4},
		{Title: "Investor", Width: 20},
		{Title: "Firm", Width: 20},
		{Title: "Net Worth", Width: 15},
		{Title: "ROI", Width: 10},
	}

	return NewGameTable(title, columns, []table.Row{})
}

// StartupTable creates a table for displaying available startups
func StartupTable(title string) *GameTable {
	columns := []table.Column{
		{Title: "#", Width: 3},
		{Title: "Name", Width: 18},
		{Title: "Category", Width: 12},
		{Title: "Valuation", Width: 12},
		{Title: "Risk", Width: 8},
		{Title: "Growth", Width: 8},
	}

	return NewGameTable(title, columns, []table.Row{})
}
