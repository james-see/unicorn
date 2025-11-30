package components

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// FormSubmitMsg is sent when a form is submitted
type FormSubmitMsg struct {
	Values map[string]string
}

// FormCancelMsg is sent when a form is cancelled
type FormCancelMsg struct{}

// FormField represents a single form field
type FormField struct {
	ID          string
	Label       string
	Placeholder string
	Default     string
	Required    bool
	Type        string // "text", "number", "password"
	MinValue    int64
	MaxValue    int64
}

// Form is a multi-field input form
type Form struct {
	title      string
	fields     []FormField
	inputs     []textinput.Model
	focusIndex int
	width      int
	height     int
	submitted  bool
	err        string
}

// NewForm creates a new form with the given fields
func NewForm(title string, fields []FormField) *Form {
	inputs := make([]textinput.Model, len(fields))
	
	for i, field := range fields {
		ti := textinput.New()
		ti.Placeholder = field.Placeholder
		ti.SetValue(field.Default)
		ti.Width = 30
		ti.CharLimit = 50
		
		if field.Type == "password" {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		
		if field.Type == "number" {
			ti.CharLimit = 15
		}
		
		if i == 0 {
			ti.Focus()
			ti.PromptStyle = styles.InputPromptStyle
			ti.TextStyle = lipgloss.NewStyle().Foreground(styles.White)
		}
		
		inputs[i] = ti
	}
	
	return &Form{
		title:      title,
		fields:     fields,
		inputs:     inputs,
		focusIndex: 0,
		width:      50,
		height:     20,
	}
}

// SetSize sets the form dimensions
func (f *Form) SetSize(width, height int) {
	f.width = width
	f.height = height
	for i := range f.inputs {
		f.inputs[i].Width = width - 10
	}
}

// Init initializes the form
func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input for the form
func (f *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Global.Tab):
			f.nextField()
		case msg.String() == "shift+tab":
			f.prevField()
		case key.Matches(msg, keys.Global.Enter):
			if f.focusIndex == len(f.fields)-1 {
				// On last field, submit
				if f.validate() {
					return f, f.submit
				}
			} else {
				f.nextField()
			}
		case key.Matches(msg, keys.Global.Back):
			return f, func() tea.Msg { return FormCancelMsg{} }
		}
	}
	
	// Update focused input
	cmd := f.updateInputs(msg)
	cmds = append(cmds, cmd)
	
	return f, tea.Batch(cmds...)
}

func (f *Form) nextField() {
	f.focusIndex++
	if f.focusIndex >= len(f.fields) {
		f.focusIndex = 0
	}
	f.updateFocus()
}

func (f *Form) prevField() {
	f.focusIndex--
	if f.focusIndex < 0 {
		f.focusIndex = len(f.fields) - 1
	}
	f.updateFocus()
}

func (f *Form) updateFocus() {
	for i := range f.inputs {
		if i == f.focusIndex {
			f.inputs[i].Focus()
			f.inputs[i].PromptStyle = styles.InputPromptStyle
			f.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(styles.White)
		} else {
			f.inputs[i].Blur()
			f.inputs[i].PromptStyle = lipgloss.NewStyle().Foreground(styles.Gray)
			f.inputs[i].TextStyle = lipgloss.NewStyle().Foreground(styles.Gray)
		}
	}
}

func (f *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(f.inputs))
	for i := range f.inputs {
		f.inputs[i], cmds[i] = f.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (f *Form) validate() bool {
	for i, field := range f.fields {
		value := f.inputs[i].Value()
		
		// Required check
		if field.Required && strings.TrimSpace(value) == "" {
			f.err = field.Label + " is required"
			f.focusIndex = i
			f.updateFocus()
			return false
		}
		
		// Number validation
		if field.Type == "number" && value != "" {
			num, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				f.err = field.Label + " must be a valid number"
				f.focusIndex = i
				f.updateFocus()
				return false
			}
			
			if field.MinValue > 0 && num < field.MinValue {
				f.err = field.Label + " must be at least " + strconv.FormatInt(field.MinValue, 10)
				f.focusIndex = i
				f.updateFocus()
				return false
			}
			
			if field.MaxValue > 0 && num > field.MaxValue {
				f.err = field.Label + " must be at most " + strconv.FormatInt(field.MaxValue, 10)
				f.focusIndex = i
				f.updateFocus()
				return false
			}
		}
	}
	
	f.err = ""
	return true
}

func (f *Form) submit() tea.Msg {
	values := make(map[string]string)
	for i, field := range f.fields {
		values[field.ID] = f.inputs[i].Value()
	}
	return FormSubmitMsg{Values: values}
}

// GetValue returns the value of a field by ID
func (f *Form) GetValue(id string) string {
	for i, field := range f.fields {
		if field.ID == id {
			return f.inputs[i].Value()
		}
	}
	return ""
}

// GetInt64Value returns the numeric value of a field
func (f *Form) GetInt64Value(id string) int64 {
	val := f.GetValue(id)
	num, _ := strconv.ParseInt(val, 10, 64)
	return num
}

// SetValue sets a field value by ID
func (f *Form) SetValue(id, value string) {
	for i, field := range f.fields {
		if field.ID == id {
			f.inputs[i].SetValue(value)
			return
		}
	}
}

// View renders the form
func (f *Form) View() string {
	var b strings.Builder
	
	// Title
	if f.title != "" {
		titleStyle := styles.TitleStyle.Width(f.width).Align(lipgloss.Center)
		b.WriteString(titleStyle.Render(f.title))
		b.WriteString("\n\n")
	}
	
	// Fields
	for i, field := range f.fields {
		// Label
		labelStyle := styles.InputLabelStyle
		if i == f.focusIndex {
			labelStyle = labelStyle.Foreground(styles.Cyan)
		}
		b.WriteString(labelStyle.Render(field.Label))
		if field.Required {
			b.WriteString(lipgloss.NewStyle().Foreground(styles.Red).Render(" *"))
		}
		b.WriteString("\n")
		
		// Input
		inputStyle := styles.InputStyle.Width(f.width - 4)
		if i == f.focusIndex {
			inputStyle = styles.InputFocusedStyle.Width(f.width - 4)
		}
		b.WriteString(inputStyle.Render(f.inputs[i].View()))
		b.WriteString("\n\n")
	}
	
	// Error message
	if f.err != "" {
		errStyle := lipgloss.NewStyle().Foreground(styles.Red)
		b.WriteString(errStyle.Render("⚠ " + f.err))
		b.WriteString("\n\n")
	}
	
	// Help
	helpStyle := styles.HelpStyle
	b.WriteString(helpStyle.Render("tab next • enter submit • esc cancel"))
	
	return b.String()
}

// SingleInput creates a form with a single text input
func SingleInput(title, label, placeholder, defaultVal string) *Form {
	return NewForm(title, []FormField{
		{
			ID:          "value",
			Label:       label,
			Placeholder: placeholder,
			Default:     defaultVal,
		},
	})
}

// NumberInput creates a form with a single number input
func NumberInput(title, label string, min, max int64, defaultVal string) *Form {
	return NewForm(title, []FormField{
		{
			ID:          "value",
			Label:       label,
			Placeholder: "Enter amount",
			Default:     defaultVal,
			Type:        "number",
			MinValue:    min,
			MaxValue:    max,
		},
	})
}

// NameAndFirmInput creates a form for entering player name and firm name
func NameAndFirmInput() *Form {
	return NewForm("🦄 PLAYER SETUP", []FormField{
		{
			ID:          "name",
			Label:       "Your Name",
			Placeholder: "Enter your name",
			Required:    true,
		},
		{
			ID:          "firm",
			Label:       "Firm Name",
			Placeholder: "e.g., Sequoia Capital",
		},
	})
}
