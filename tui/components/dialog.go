package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// DialogResult represents the result of a dialog
type DialogResult int

const (
	DialogConfirm DialogResult = iota
	DialogCancel
	DialogOptionA
	DialogOptionB
)

// DialogResultMsg is sent when a dialog is closed
type DialogResultMsg struct {
	ID     string
	Result DialogResult
}

// Dialog types
type DialogType int

const (
	DialogInfo DialogType = iota
	DialogConfirmation
	DialogChoice
	DialogError
	DialogSuccess
)

// Dialog is a modal dialog component
type Dialog struct {
	id           string
	dialogType   DialogType
	title        string
	message      string
	optionA      string
	optionB      string
	selectedIdx  int
	width        int
}

// NewDialog creates a new dialog
func NewDialog(id string, dialogType DialogType, title, message string) *Dialog {
	return &Dialog{
		id:         id,
		dialogType: dialogType,
		title:      title,
		message:    message,
		optionA:    "OK",
		optionB:    "Cancel",
		width:      50,
	}
}

// SetOptions sets custom button labels
func (d *Dialog) SetOptions(optionA, optionB string) {
	d.optionA = optionA
	d.optionB = optionB
}

// SetWidth sets the dialog width
func (d *Dialog) SetWidth(width int) {
	d.width = width
}

// Init initializes the dialog
func (d *Dialog) Init() tea.Cmd {
	return nil
}

// Update handles input for the dialog
func (d *Dialog) Update(msg tea.Msg) (*Dialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Left), msg.String() == "h":
			if d.dialogType == DialogConfirmation || d.dialogType == DialogChoice {
				d.selectedIdx = 0
			}
		case key.Matches(msg, keys.Global.Right), msg.String() == "l":
			if d.dialogType == DialogConfirmation || d.dialogType == DialogChoice {
				d.selectedIdx = 1
			}
		case key.Matches(msg, keys.Global.Tab):
			if d.dialogType == DialogConfirmation || d.dialogType == DialogChoice {
				d.selectedIdx = (d.selectedIdx + 1) % 2
			}
		case key.Matches(msg, keys.Global.Enter):
			return d, d.confirm
		case key.Matches(msg, keys.Global.Back):
			if d.dialogType == DialogConfirmation || d.dialogType == DialogChoice {
				return d, d.cancel
			}
			return d, d.confirm // For info dialogs, esc also closes
		case msg.String() == "y", msg.String() == "Y":
			if d.dialogType == DialogConfirmation {
				d.selectedIdx = 0
				return d, d.confirm
			}
		case msg.String() == "n", msg.String() == "N":
			if d.dialogType == DialogConfirmation {
				d.selectedIdx = 1
				return d, d.cancel
			}
		}
	}
	return d, nil
}

func (d *Dialog) confirm() tea.Msg {
	result := DialogConfirm
	if d.dialogType == DialogChoice {
		if d.selectedIdx == 0 {
			result = DialogOptionA
		} else {
			result = DialogOptionB
		}
	}
	return DialogResultMsg{ID: d.id, Result: result}
}

func (d *Dialog) cancel() tea.Msg {
	return DialogResultMsg{ID: d.id, Result: DialogCancel}
}

// View renders the dialog
func (d *Dialog) View() string {
	var b strings.Builder

	// Determine border color based on type
	borderColor := styles.Cyan
	icon := ""
	switch d.dialogType {
	case DialogError:
		borderColor = styles.Red
		icon = "❌ "
	case DialogSuccess:
		borderColor = styles.Green
		icon = "✅ "
	case DialogConfirmation:
		borderColor = styles.Yellow
		icon = "❓ "
	case DialogChoice:
		borderColor = styles.Magenta
		icon = "🔀 "
	default:
		icon = "ℹ️  "
	}

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(borderColor).
		Bold(true).
		Width(d.width - 4).
		Align(lipgloss.Center)
	b.WriteString(titleStyle.Render(icon + d.title))
	b.WriteString("\n\n")

	// Message
	messageStyle := lipgloss.NewStyle().
		Foreground(styles.White).
		Width(d.width - 4).
		Align(lipgloss.Center)
	b.WriteString(messageStyle.Render(d.message))
	b.WriteString("\n\n")

	// Buttons
	if d.dialogType == DialogConfirmation || d.dialogType == DialogChoice {
		// Two buttons
		btnA := d.optionA
		btnB := d.optionB

		btnAStyle := styles.DialogButtonStyle
		btnBStyle := styles.DialogButtonStyle

		if d.selectedIdx == 0 {
			btnAStyle = styles.DialogButtonActiveStyle
		} else {
			btnBStyle = styles.DialogButtonActiveStyle
		}

		buttons := lipgloss.JoinHorizontal(
			lipgloss.Center,
			btnAStyle.Render(btnA),
			btnBStyle.Render(btnB),
		)

		buttonContainer := lipgloss.NewStyle().
			Width(d.width - 4).
			Align(lipgloss.Center)
		b.WriteString(buttonContainer.Render(buttons))
	} else {
		// Single OK button
		btnStyle := styles.DialogButtonActiveStyle
		buttonContainer := lipgloss.NewStyle().
			Width(d.width - 4).
			Align(lipgloss.Center)
		b.WriteString(buttonContainer.Render(btnStyle.Render(d.optionA)))
	}

	// Wrap in box
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(d.width)

	return boxStyle.Render(b.String())
}

// Helper functions to create common dialogs

// InfoDialog creates an informational dialog
func InfoDialog(id, title, message string) *Dialog {
	return NewDialog(id, DialogInfo, title, message)
}

// ConfirmDialog creates a yes/no confirmation dialog
func ConfirmDialog(id, title, message string) *Dialog {
	d := NewDialog(id, DialogConfirmation, title, message)
	d.SetOptions("Yes", "No")
	return d
}

// ErrorDialog creates an error dialog
func ErrorDialog(id, title, message string) *Dialog {
	return NewDialog(id, DialogError, title, message)
}

// SuccessDialog creates a success dialog
func SuccessDialog(id, title, message string) *Dialog {
	return NewDialog(id, DialogSuccess, title, message)
}

// ChoiceDialog creates a choice dialog with custom options
func ChoiceDialog(id, title, message, optionA, optionB string) *Dialog {
	d := NewDialog(id, DialogChoice, title, message)
	d.SetOptions(optionA, optionB)
	return d
}

// Overlay renders a dialog centered over content
func Overlay(content string, dialog *Dialog, width, height int) string {
	// Dim the background content (for future use in layered rendering)
	_ = content // Background content would be rendered behind dialog in full implementation

	// Center the dialog
	dialogView := dialog.View()

	// Place dialog in center
	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		dialogView,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(styles.DarkGray),
	)
}
