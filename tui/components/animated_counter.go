package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// AnimTickMsg drives the animation loop at ~60fps.
// Exported so screens can intercept it in their Update methods.
type AnimTickMsg struct{}

func animTick() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return AnimTickMsg{}
	})
}

// AnimatedCounter smoothly counts from 0 to a target value using spring physics.
// Drop it into any ScreenModel that wants a spring-animated number display.
type AnimatedCounter struct {
	spring    harmonica.Spring
	pos       float64 // current interpolated value
	velocity  float64 // current velocity
	target    float64 // target value
	width     int
	label     string
	prefix    string // e.g. "$" for money
	suffix    string
	style     lipgloss.Style
	tickStyle lipgloss.Style
	done      bool
	animating bool
}

// NewAnimatedCounter creates a counter that springs from 0 to target.
// Uses critically-damped spring (dampingRatio=1.0) for smooth deceleration
// without overshoot. Angular frequency 6.0 gives ~0.5s settle time.
func NewAnimatedCounter(target int64, label string, width int) *AnimatedCounter {
	spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 1.0)

	return &AnimatedCounter{
		spring:    spring,
		pos:        0,
		velocity:   0,
		target:     float64(target),
		width:      width,
		label:      label,
		prefix:     "$",
		style:      lipgloss.NewStyle().Foreground(styles.Green).Bold(true),
		tickStyle:  lipgloss.NewStyle().Foreground(styles.Gray),
		done:       false,
		animating:  false,
	}
}

// SetPrefix changes the prefix (default "$").
func (c *AnimatedCounter) SetPrefix(p string) { c.prefix = p }

// SetSuffix changes the suffix.
func (c *AnimatedCounter) SetSuffix(s string) { c.suffix = s }

// SetStyle overrides the value render style.
func (c *AnimatedCounter) SetStyle(s lipgloss.Style) { c.style = s }

// SetTickStyle sets the style for the tick indicator shown while animating.
func (c *AnimatedCounter) SetTickStyle(s lipgloss.Style) { c.tickStyle = s }

// Done returns true when the spring has settled on the target.
func (c *AnimatedCounter) Done() bool { return c.done }

// CurrentValue returns the current interpolated value as int64.
func (c *AnimatedCounter) CurrentValue() int64 { return int64(c.pos) }

// Init starts the animation.
func (c *AnimatedCounter) Init() tea.Cmd {
	c.animating = true
	return animTick()
}

// Update advances the spring on each tick.
func (c *AnimatedCounter) Update(msg tea.Msg) (*AnimatedCounter, tea.Cmd) {
	switch msg.(type) {
	case AnimTickMsg:
		if c.done {
			return c, nil
		}
		// Advance the spring one frame toward the target
		c.pos, c.velocity = c.spring.Update(c.pos, c.velocity, c.target)

		// Check if spring has settled — use relative threshold for large numbers
		diff := c.pos - c.target
		if diff < 0 {
			diff = -diff
		}
		velAbs := c.velocity
		if velAbs < 0 {
			velAbs = -velAbs
		}
		// Threshold scales with target magnitude: 0.01% of target or 1.0, whichever is larger
		threshold := c.target * 0.0001
		if threshold < 1.0 {
			threshold = 1.0
		}
		if diff < threshold && velAbs < threshold {
			c.pos = c.target
			c.done = true
			c.animating = false
			return c, nil
		}
		return c, animTick()
	}
	return c, nil
}

// View renders the counter with label and animated value.
func (c *AnimatedCounter) View() string {
	var b strings.Builder

	if c.label != "" {
		labelStyle := lipgloss.NewStyle().
			Foreground(styles.Yellow).
			Width(c.width).
			Align(lipgloss.Left)
		b.WriteString(labelStyle.Render(c.label))
		b.WriteString("\n")
	}

	// Render the value with formatting
	valStr := fmt.Sprintf("%s%s%s", c.prefix, formatCounterMoney(int64(c.pos)), c.suffix)

	valueLine := c.style.Render(valStr)

	// Show a subtle tick while animating
	if c.animating && !c.done {
		tick := c.tickStyle.Render(" ▸")
		valueLine += tick
	}

	valueStyle := lipgloss.NewStyle().Width(c.width).Align(lipgloss.Left)
	b.WriteString(valueStyle.Render(valueLine))

	return b.String()
}

// formatCounterMoney formats int64 as currency with commas (matching styles.formatCurrency).
func formatCounterMoney(amount int64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	str := ""
	for amount > 0 {
		if str != "" {
			str = "," + str
		}
		if amount >= 1000 {
			str = fmt.Sprintf("%03d", amount%1000) + str
		} else {
			str = fmt.Sprintf("%d", amount%1000) + str
		}
		amount /= 1000
	}
	if str == "" {
		str = "0"
	}
	if negative {
		return "-" + str
	}
	return str
}