package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
)

// Inline styles for splash screen
var (
	splashMagenta = lipgloss.Color("#FF00FF")
	splashCyan    = lipgloss.Color("#00FFFF")
	splashYellow  = lipgloss.Color("#FFFF00")
	splashGray    = lipgloss.Color("#808080")
)

// ASCII art unicorn logo
const unicornLogo = `
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣤⣄⡀⠀⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⢀⣾⣿⣿⣿⣿⣿⣷⡀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⢰⣿⣿⣿⣿⣿⣿⣿⣿⣿⡆⠀
    ⠀⠀⠀⠀⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⡇⠀
    ⠀⠀⠀⠀⠀⠀⠀⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿⠁⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⢻⣿⣿⣿⣿⣿⣿⣿⡟⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⢿⣿⣿⣿⡿⠋⠀⠀⠀
    ⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠉⠁⠀⠀⠀⠀⠀
`

const unicornTitle = `
██╗   ██╗███╗   ██╗██╗ ██████╗ ██████╗ ██████╗ ███╗   ██╗
██║   ██║████╗  ██║██║██╔════╝██╔═══██╗██╔══██╗████╗  ██║
██║   ██║██╔██╗ ██║██║██║     ██║   ██║██████╔╝██╔██╗ ██║
██║   ██║██║╚██╗██║██║██║     ██║   ██║██╔══██╗██║╚██╗██║
╚██████╔╝██║ ╚████║██║╚██████╗╚██████╔╝██║  ██║██║ ╚████║
 ╚═════╝ ╚═╝  ╚═══╝╚═╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝
`

// Splash taglines cycle through these
var splashTaglines = []string{
	"The Ultimate Startup Adventure",
	"Invest as a VC, or Build as a Founder",
	"Hunt for Unicorns, or Become One",
	"Where Legends Are Made",
}

// SplashScreen is the animated splash screen with spring physics.
type SplashScreen struct {
	width     int
	height    int
	spinner   spinner.Model
	startTime time.Time

	// Spring-animated values for smooth transitions
	fadeSpring    harmonica.Spring // title fade-in (0→1)
	infoSpring    harmonica.Spring // info box slide-in (0→1)
	promptSpring  harmonica.Spring // "press any key" pulse (0→1→0→...)
	fadePos       float64
	fadeVel       float64
	infoPos       float64
	infoVel       float64
	promptPos     float64
	promptVel     float64

	// Timing
	phase     int   // which tagline to show
	ready     bool  // can accept keypress
	infoShown bool  // info box started appearing
}

// NewSplashScreen creates a new splash screen
func NewSplashScreen(width, height int) *SplashScreen {
	s := spinner.New()
	s.Spinner = spinner.Spinner{
		Frames: []string{"🦄", "✨🦄", "✨✨🦄", "✨🦄", "🦄"},
		FPS:    time.Second / 6,
	}
	s.Style = lipgloss.NewStyle()

	// Critically-damped springs: smooth, no overshoot.
	// Angular freq 3.0 = ~0.8s settle, 4.0 = ~0.6s, 6.0 = ~0.4s
	dt := harmonica.FPS(60)
	return &SplashScreen{
		width:        width,
		height:       height,
		spinner:      s,
		startTime:    time.Now(),
		fadeSpring:   harmonica.NewSpring(dt, 3.0, 1.0), // slow title fade
		infoSpring:   harmonica.NewSpring(dt, 4.0, 1.0),  // medium info box
		promptSpring: harmonica.NewSpring(dt, 2.0, 1.0),  // gentle prompt
		fadePos:      0,
		fadeVel:      0,
		infoPos:      0,
		infoVel:      0,
		promptPos:    0,
		promptVel:    0,
	}
}

// splashTickMsg is sent to advance the splash animation
type splashTickMsg struct{}

func splashTick() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

// Init initializes the splash screen
func (s *SplashScreen) Init() tea.Cmd {
	// Kick off the title fade-in spring
	s.fadePos, s.fadeVel = s.fadeSpring.Update(s.fadePos, s.fadeVel, 1.0)
	return tea.Batch(s.spinner.Tick, splashTick())
}

// Update handles splash screen updates
func (s *SplashScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Any key press advances to main menu
		return s, SwitchTo(ScreenMainMenu)

	case splashTickMsg:
		elapsed := time.Since(s.startTime)

		// Advance title fade spring toward 1.0
		s.fadePos, s.fadeVel = s.fadeSpring.Update(s.fadePos, s.fadeVel, 1.0)

		// After ~1.2s, start info box slide-in
		if elapsed > 1200*time.Millisecond && !s.infoShown {
			s.infoShown = true
		}
		if s.infoShown {
			s.infoPos, s.infoVel = s.infoSpring.Update(s.infoPos, s.infoVel, 1.0)
		}

		// Cycle taglines every ~800ms for the first 3.2s
		taglineIdx := int(elapsed / (800 * time.Millisecond))
		if taglineIdx >= len(splashTaglines) {
			taglineIdx = len(splashTaglines) - 1
		}
		s.phase = taglineIdx

		// After ~2.5s, start prompt pulse and mark ready
		if elapsed > 2500*time.Millisecond {
			s.ready = true
			// Pulse the prompt: target alternates between 1.0 and 0.3
			target := 1.0
			if int(elapsed/(600*time.Millisecond))%2 == 1 {
				target = 0.3
			}
			s.promptPos, s.promptVel = s.promptSpring.Update(s.promptPos, s.promptVel, target)
		}

		return s, splashTick()

	case spinner.TickMsg:
		var cmd tea.Cmd
		s.spinner, cmd = s.spinner.Update(msg)
		return s, cmd
	}

	return s, nil
}

// View renders the splash screen
func (s *SplashScreen) View() string {
	var b strings.Builder

	// Title with spring-based opacity (simulated via color blending)
	// We can't do real alpha in terminal, so we interpolate between gray→magenta
	titleColor := blendColor(splashGray, splashMagenta, s.fadePos)
	titleStyle := lipgloss.NewStyle().
		Foreground(titleColor).
		Bold(true)

	b.WriteString(titleStyle.Render(unicornTitle))
	b.WriteString("\n")

	// Spinner
	spinnerStyle := lipgloss.NewStyle().
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(spinnerStyle.Render(s.spinner.View()))
	b.WriteString("\n\n")

	// Tagline — also fades in with the title spring
	taglineStyle := lipgloss.NewStyle().
		Foreground(splashCyan).
		Italic(true).
		Width(s.width).
		Align(lipgloss.Center)

	tagline := splashTaglines[s.phase]
	if s.phase >= len(splashTaglines)-1 && s.ready {
		tagline = "Where Dreams Become Unicorns... Or Don't"
	}
	b.WriteString(taglineStyle.Render(tagline))
	b.WriteString("\n\n")

	// Info box slides in with spring (after delay)
	if s.infoShown {
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(blendColor(splashGray, splashCyan, s.infoPos)).
			Padding(1, 2).
			Width(50).
			Align(lipgloss.Center)

		info := "🚀 Invest in startups\n💰 Compete against AI VCs\n🏆 Build your reputation\n🦄 Find the next unicorn"

		// As the spring approaches 1.0, the info becomes fully visible
		infoOpacity := s.infoPos
		if infoOpacity > 1 {
			infoOpacity = 1
		}
		infoColor := blendColor(splashGray, lipgloss.Color("#FFFFFF"), infoOpacity)
		infoText := lipgloss.NewStyle().Foreground(infoColor).Render(info)
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(infoBox.Render(infoText)))
		b.WriteString("\n\n")
	}

	// Press any key prompt with pulsing spring
	if s.ready {
		promptOpacity := s.promptPos
		if promptOpacity > 1 {
			promptOpacity = 1
		}
		if promptOpacity < 0 {
			promptOpacity = 0
		}
		promptColor := blendColor(lipgloss.Color("#404040"), splashYellow, promptOpacity)
		promptStyle := lipgloss.NewStyle().
			Foreground(promptColor).
			Bold(true).
			Width(s.width).
			Align(lipgloss.Center)
		b.WriteString(promptStyle.Render("Press any key to continue..."))
	}

	// Copyright
	b.WriteString("\n\n")
	copyrightStyle := lipgloss.NewStyle().
		Foreground(splashGray).
		Width(s.width).
		Align(lipgloss.Center)
	b.WriteString(copyrightStyle.Render("© 2025 Unicorn Game"))

	return b.String()
}

// blendColor interpolates between two lipgloss colors based on t (0.0→1.0).
// Colors are hex strings like "#FF00FF".
func blendColor(a, b lipgloss.Color, t float64) lipgloss.Color {
	if t <= 0 {
		return a
	}
	if t >= 1 {
		return b
	}
	ar, ag, ab := hexToRGB(string(a))
	br, bg, bb := hexToRGB(string(b))
	r := ar + int(float64(br-ar)*t)
	g := ag + int(float64(bg-ag)*t)
	bl := ab + int(float64(bb-ab)*t)
	return lipgloss.Color(rgbToHex(r, g, bl))
}

func hexToRGB(hex string) (int, int, int) {
	// Handle "#RRGGBB" format
	if len(hex) < 7 {
		return 128, 128, 128
	}
	hex = hex[1:] // strip #
	r := hexCharToInt(hex[0])*16 + hexCharToInt(hex[1])
	g := hexCharToInt(hex[2])*16 + hexCharToInt(hex[3])
	b := hexCharToInt(hex[4])*16 + hexCharToInt(hex[5])
	return r, g, b
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02X%02X%02X", clamp255(r), clamp255(g), clamp255(b))
}

func clamp255(v int) int {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return v
}

func hexCharToInt(c byte) int {
	switch {
	case c >= '0' && c <= '9':
		return int(c - '0')
	case c >= 'a' && c <= 'f':
		return int(c-'a') + 10
	case c >= 'A' && c <= 'F':
		return int(c-'A') + 10
	}
	return 0
}