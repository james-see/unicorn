package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jamesacampbell/unicorn/database"
	"github.com/jamesacampbell/unicorn/tui/keys"
	"github.com/jamesacampbell/unicorn/tui/styles"
)

// getLevelTitleLocal returns the title for a given level
func getLevelTitleLocal(level int) string {
	titles := map[int]string{
		1:  "Aspiring VC",
		2:  "Associate",
		3:  "Senior Associate",
		4:  "Principal",
		5:  "Junior Partner",
		6:  "Partner",
		7:  "Senior Partner",
		8:  "Managing Partner",
		9:  "Founding Partner",
		10: "Legend",
	}
	if title, ok := titles[level]; ok {
		return title
	}
	if level > 10 {
		return "VC Legend"
	}
	return "Newcomer"
}

// ProgressionScreen shows player progression
type ProgressionScreen struct {
	width      int
	height     int
	playerName string
	nameInput  textinput.Model
	profile    *database.PlayerProfile
	inputMode  bool
}

// NewProgressionScreen creates a new progression screen
func NewProgressionScreen(width, height int, playerName string) *ProgressionScreen {
	nameInput := textinput.New()
	nameInput.Placeholder = "Enter player name"
	nameInput.CharLimit = 30
	nameInput.Width = 25
	
	s := &ProgressionScreen{
		width:      width,
		height:     height,
		nameInput:  nameInput,
		playerName: playerName,
	}
	
	if playerName != "" {
		profile, err := database.GetPlayerProfile(playerName)
		if err == nil {
			s.profile = profile
			s.inputMode = false
		} else {
			s.inputMode = true
			s.nameInput.Focus()
		}
	} else {
		s.inputMode = true
		s.nameInput.Focus()
	}
	
	return s
}

// Init initializes the progression screen
func (s *ProgressionScreen) Init() tea.Cmd {
	if s.inputMode {
		return textinput.Blink
	}
	return nil
}

// Update handles progression input
func (s *ProgressionScreen) Update(msg tea.Msg) (ScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if key.Matches(msg, keys.Global.Back) {
			if s.inputMode && s.playerName != "" {
				s.inputMode = false
				return s, nil
			}
			return s, PopScreen()
		}
		
		if s.inputMode && msg.Type == tea.KeyEnter {
			name := strings.TrimSpace(s.nameInput.Value())
			if name != "" {
				s.playerName = name
				profile, err := database.GetPlayerProfile(name)
				if err == nil {
					s.profile = profile
				}
				s.inputMode = false
			}
			return s, nil
		}
	}
	
	if s.inputMode {
		var cmd tea.Cmd
		s.nameInput, cmd = s.nameInput.Update(msg)
		return s, cmd
	}
	
	return s, nil
}

// View renders progression
func (s *ProgressionScreen) View() string {
	var b strings.Builder
	
	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(styles.Black).
		Background(styles.Magenta).
		Bold(true).
		Width(60).
		Align(lipgloss.Center)
	
	b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(headerStyle.Render("📈 PROGRESSION & LEVELS 📈")))
	b.WriteString("\n\n")
	
	if s.inputMode {
		// Name input
		inputBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(1, 2).
			Width(40)
		
		var input strings.Builder
		input.WriteString("Enter player name:\n\n")
		input.WriteString(s.nameInput.View())
		
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(inputBox.Render(input.String())))
	} else {
		// Progression display
		progBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Magenta).
			Padding(1, 2).
			Width(55)
		
		var prog strings.Builder
		titleStyle := lipgloss.NewStyle().Foreground(styles.Magenta).Bold(true)
		prog.WriteString(titleStyle.Render(fmt.Sprintf("Progress for %s", s.playerName)))
		prog.WriteString("\n\n")
		
		if s.profile == nil {
			prog.WriteString("No progression data yet - play some games!")
		} else {
			// Level info
			levelStyle := lipgloss.NewStyle().Foreground(styles.Gold).Bold(true)
			title := getLevelTitleLocal(s.profile.Level)
			prog.WriteString(levelStyle.Render(fmt.Sprintf("Level %d - %s", s.profile.Level, title)))
			prog.WriteString("\n\n")
			
			// XP bar
			xpProgress := s.profile.ProgressPercent / 100.0
			if xpProgress > 1 {
				xpProgress = 1
			}
			
			barWidth := 30
			filledWidth := int(xpProgress * float64(barWidth))
			
			labelStyle := lipgloss.NewStyle().Foreground(styles.Yellow)
			prog.WriteString(labelStyle.Render("Experience: "))
			
			filledStyle := lipgloss.NewStyle().Foreground(styles.Green)
			emptyStyle := lipgloss.NewStyle().Foreground(styles.DarkGray)
			prog.WriteString(filledStyle.Render(strings.Repeat("█", filledWidth)))
			prog.WriteString(emptyStyle.Render(strings.Repeat("░", barWidth-filledWidth)))
			prog.WriteString(fmt.Sprintf(" %d/%d XP\n\n", s.profile.ExperiencePoints, s.profile.NextLevelXP))
			
			// Stats
			prog.WriteString(labelStyle.Render("Total Points Earned: "))
			prog.WriteString(fmt.Sprintf("%d\n", s.profile.TotalPointsEarned))
			
			prog.WriteString(labelStyle.Render("Level-up Points: "))
			prog.WriteString(fmt.Sprintf("%d\n", s.profile.LevelUpPoints))
			
			// Unlocks
			prog.WriteString("\n")
			prog.WriteString(titleStyle.Render("Unlocks:\n"))
			
			unlockStyle := lipgloss.NewStyle().Foreground(styles.Green)
			lockedStyle := lipgloss.NewStyle().Foreground(styles.Gray)
			
			if s.profile.Level >= 2 {
				prog.WriteString(unlockStyle.Render("✓ Syndicate Investing (Lvl 2)\n"))
			} else {
				prog.WriteString(lockedStyle.Render("🔒 Syndicate Investing (Lvl 2)\n"))
			}
			
			if s.profile.Level >= 5 {
				prog.WriteString(unlockStyle.Render("✓ Hard Difficulty (Lvl 5)\n"))
			} else {
				prog.WriteString(lockedStyle.Render("🔒 Hard Difficulty (Lvl 5)\n"))
			}
			
			if s.profile.Level >= 10 {
				prog.WriteString(unlockStyle.Render("✓ Expert Difficulty (Lvl 10)\n"))
			} else {
				prog.WriteString(lockedStyle.Render("🔒 Expert Difficulty (Lvl 10)\n"))
			}
		}
		
		b.WriteString(lipgloss.NewStyle().Width(s.width).Align(lipgloss.Center).Render(progBox.Render(prog.String())))
	}
	
	b.WriteString("\n\n")
	
	// Help
	helpStyle := lipgloss.NewStyle().Foreground(styles.Gray).Width(s.width).Align(lipgloss.Center)
	b.WriteString(helpStyle.Render("esc back"))
	
	return b.String()
}
