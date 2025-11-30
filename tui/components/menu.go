package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// MenuItem represents a single menu item
type MenuItem struct {
	ID          string
	Title       string
	Description string
	Icon        string
	Disabled    bool
	Hidden      bool
}

// MenuSelectedMsg is sent when a menu item is selected
type MenuSelectedMsg struct {
	ID string
}

// Menu is a list-based menu component
type Menu struct {
	title     string
	items     []MenuItem
	cursor    int
	width     int
	height    int
	showIcons bool
	showDesc  bool
	hideHelp  bool
}

// NewMenu creates a new menu component
func NewMenu(title string, items []MenuItem) *Menu {
	return &Menu{
		title:     title,
		items:     items,
		cursor:    0,
		width:     50,
		height:    20,
		showIcons: true,
		showDesc:  true,
	}
}

// SetSize sets the menu dimensions
func (m *Menu) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetShowIcons toggles icon display
func (m *Menu) SetShowIcons(show bool) {
	m.showIcons = show
}

// SetShowDescriptions toggles description display
func (m *Menu) SetShowDescriptions(show bool) {
	m.showDesc = show
}

// SetHideHelp hides the built-in help text
func (m *Menu) SetHideHelp(hide bool) {
	m.hideHelp = hide
}

// Init initializes the menu
func (m *Menu) Init() tea.Cmd {
	return nil
}

// Update handles input for the menu
func (m *Menu) Update(msg tea.Msg) (*Menu, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Menu.Up):
			m.moveUp()
		case key.Matches(msg, keys.Menu.Down):
			m.moveDown()
		case key.Matches(msg, keys.Menu.Select):
			if !m.items[m.cursor].Disabled {
				return m, m.selectItem
			}
		}
	}
	return m, nil
}

func (m *Menu) moveUp() {
	startCursor := m.cursor
	for {
		m.cursor--
		if m.cursor < 0 {
			m.cursor = len(m.items) - 1
		}
		// Skip only hidden items (not disabled - we can still navigate to them)
		if !m.items[m.cursor].Hidden {
			break
		}
		// Prevent infinite loop if all items are hidden
		if m.cursor == startCursor {
			break
		}
	}
}

func (m *Menu) moveDown() {
	startCursor := m.cursor
	for {
		m.cursor++
		if m.cursor >= len(m.items) {
			m.cursor = 0
		}
		// Skip only hidden items (not disabled - we can still navigate to them)
		if !m.items[m.cursor].Hidden {
			break
		}
		// Prevent infinite loop if all items are hidden
		if m.cursor == startCursor {
			break
		}
	}
}

func (m *Menu) selectItem() tea.Msg {
	return MenuSelectedMsg{ID: m.items[m.cursor].ID}
}

// SelectedID returns the currently selected item's ID
func (m *Menu) SelectedID() string {
	return m.items[m.cursor].ID
}

// SelectedIndex returns the cursor position
func (m *Menu) SelectedIndex() int {
	return m.cursor
}

// View renders the menu
func (m *Menu) View() string {
	var b strings.Builder

	// Title
	if m.title != "" {
		titleStyle := styles.TitleStyle.Width(m.width).Align(lipgloss.Center)
		b.WriteString(titleStyle.Render(m.title))
		b.WriteString("\n\n")
	}

	// Menu items
	visibleCount := 0
	for i, item := range m.items {
		if item.Hidden {
			continue
		}

		// Build item text
		var itemText string
		if m.showIcons && item.Icon != "" {
			itemText = fmt.Sprintf("%s  %s", item.Icon, item.Title)
		} else {
			itemText = item.Title
		}

		// Style based on selection and state
		var itemStyle lipgloss.Style
		if i == m.cursor {
			itemStyle = styles.SelectedMenuItemStyle.Width(m.width - 4)
		} else if item.Disabled {
			itemStyle = lipgloss.NewStyle().
				Foreground(styles.Gray).
				Width(m.width-4).
				Padding(0, 2)
		} else {
			itemStyle = styles.MenuItemStyle.Width(m.width - 4)
		}

		b.WriteString(itemStyle.Render(itemText))
		b.WriteString("\n")

		// Description
		if m.showDesc && item.Description != "" && i == m.cursor {
			descStyle := styles.MenuDescriptionStyle.Width(m.width - 6)
			b.WriteString(descStyle.Render(item.Description))
			b.WriteString("\n")
		}

		visibleCount++
	}

	// Help text (optional)
	if !m.hideHelp {
		b.WriteString("\n")
		helpStyle := styles.HelpStyle
		b.WriteString(helpStyle.Render("↑/↓ navigate • enter select • q quit"))
	}

	return b.String()
}

// SimpleMenu creates a simple menu with just titles
func SimpleMenu(title string, options []string) *Menu {
	items := make([]MenuItem, len(options))
	for i, opt := range options {
		items[i] = MenuItem{
			ID:    fmt.Sprintf("%d", i),
			Title: opt,
		}
	}
	return NewMenu(title, items)
}
