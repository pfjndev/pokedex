package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Type badge styles with colors matching Pokemon types
var (
	typeBadgeStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				MarginTop(1).
				MarginBottom(1)

	statNameStyle = lipgloss.NewStyle().
			Width(18)

	statValueStyle = lipgloss.NewStyle().
			Width(3).
			Align(lipgloss.Right)
)

// RenderTypeBadge renders a colored type badge like "[⚡ Electric]"
func RenderTypeBadge(typeName, color string) string {
	// Map type symbols for visual appeal
	typeSymbols := map[string]string{
		"electric": "⚡",
		"fire":     "🔥",
		"water":    "💧",
		"grass":    "🌿",
		"ice":      "❄️",
		"fighting": "🥊",
		"poison":   "☠️",
		"ground":   "🌍",
		"flying":   "🦅",
		"psychic":  "🔮",
		"bug":      "🐛",
		"rock":     "🪨",
		"ghost":    "👻",
		"dragon":   "🐉",
		"dark":     "🌙",
		"steel":    "⚙️",
		"fairy":    "✨",
		"normal":   "⭐",
	}

	symbol := typeSymbols[strings.ToLower(typeName)]
	if symbol == "" {
		symbol = "●"
	}

	badgeText := fmt.Sprintf("%s %s", symbol, cases.Title(language.English).String(typeName))
	
	style := typeBadgeStyle
	if color != "" {
		style = style.Foreground(lipgloss.Color(color))
	}
	
	return style.Render(badgeText)
}

// RenderSectionHeader renders a section header with consistent styling
func RenderSectionHeader(title string) string {
	// Create a decorative header with box drawing characters
	separator := strings.Repeat("═", 3)
	headerText := fmt.Sprintf("%s %s %s", separator, title, separator)
	
	return sectionHeaderStyle.Render(headerText)
}

// RenderStatBar renders a stat with name, value, and visual bar
func RenderStatBar(name string, value, maxValue int, color string) string {
	// Format the stat name (left-aligned, fixed width)
	statName := statNameStyle.Render(name + ":")
	
	// Format the value (right-aligned, fixed width)
	statValue := statValueStyle.Render(fmt.Sprintf("%d", value))
	
	// Calculate percentage
	percentage := 0
	if maxValue > 0 {
		percentage = (value * 100) / maxValue
		if percentage > 100 {
			percentage = 100
		}
	}
	
	// Render the progress bar
	barWidth := 20
	progressBar := RenderProgressBar(value, maxValue, barWidth, color)
	
	// Format percentage display
	percentStr := fmt.Sprintf("%3d%%", percentage)
	
	// Combine all parts
	return fmt.Sprintf("%s %s  %s %s", statName, statValue, progressBar, percentStr)
}

// RenderProgressBar renders just the progress bar
func RenderProgressBar(value, maxValue int, width int, color string) string {
	if maxValue <= 0 {
		maxValue = 1
	}
	
	// Calculate how many blocks to fill
	filledLength := (value * width) / maxValue
	if filledLength > width {
		filledLength = width
	}
	if filledLength < 0 {
		filledLength = 0
	}
	
	// Build the bar
	filled := strings.Repeat("█", filledLength)
	empty := strings.Repeat("░", width-filledLength)
	bar := filled + empty
	
	// Apply color if provided
	if color != "" {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color(color))
		return style.Render(bar)
	}
	
	return bar
}

// WrapText wraps long text to specified width
func WrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		// If adding this word would exceed width, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+len(word)+1 > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}
		
		// Add word to current line
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	
	// Add the last line
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	return strings.Join(lines, "\n")
}
