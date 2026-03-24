package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pfjndev/pokedex/internal/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	selectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(lipgloss.Color("170")).
				Bold(true)

	// Enhanced menu item styles
	unselectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("246"))

	menuIconStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))

	selectedMenuIconStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("170"))

	menuSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true).
			Padding(1)

	// Error type specific styles
	notFoundErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E74C3C")).
				Bold(true)

	networkErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#E67E22")).
				Bold(true)

	apiErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9B59B6")).
			Bold(true)

	validationErrorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F39C12")).
				Bold(true)

	errorSuggestionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("246")).
				Italic(true).
				MarginTop(1)

	pokemonNameStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("170")).
				MarginBottom(1)

	sectionStyle = lipgloss.NewStyle().
			MarginTop(1)

	statBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39"))
)

// renderMenu renders the main menu view
func (m Model) renderMenu() string {
	var b strings.Builder

	// Render section header with decorative styling
	b.WriteString(RenderSectionHeader("🌟 Pokedex CLI"))
	b.WriteString("\n")

	// Menu items with emoji/symbols
	menuItems := []struct {
		key   string
		emoji string
		desc  string
	}{
		{"s", "🔍", "Search for a Pokemon"},
		{"b", "📋", "Browse Pokemon list"},
		{"q", "👋", "Quit"},
	}

	// Visual separator before menu items
	separator := menuSeparatorStyle.Render(strings.Repeat("─", 40))
	b.WriteString(separator)
	b.WriteString("\n\n")

	// Render menu items with enhanced styling
	for i, item := range menuItems {
		if i == m.menuCursor {
			// Selected item - prominent styling
			icon := selectedMenuIconStyle.Render(item.emoji)
			text := fmt.Sprintf("%s  %s", icon, item.desc)
			line := selectedMenuItemStyle.Render(fmt.Sprintf("▸ %s", text))
			b.WriteString(line)
		} else {
			// Unselected item - subdued styling
			icon := menuIconStyle.Render(item.emoji)
			text := fmt.Sprintf("%s  %s", icon, item.desc)
			shortcut := helpStyle.Render(fmt.Sprintf("[%s]", item.key))
			line := unselectedMenuItemStyle.Render(fmt.Sprintf("  %s %s", text, shortcut))
			b.WriteString(line)
		}
		b.WriteString("\n")
	}

	// Visual separator after menu items
	b.WriteString("\n")
	b.WriteString(separator)
	b.WriteString("\n")

	// Enhanced help text with styled key hints
	b.WriteString("\n")
	helpParts := []string{
		fmt.Sprintf("%s: Navigate", helpKeyStyle.Render("↑↓/j/k")),
		fmt.Sprintf("%s: Select", helpKeyStyle.Render("Enter")),
		fmt.Sprintf("%s: Quick access", helpKeyStyle.Render("Letter keys")),
	}
	helpText := strings.Join(helpParts, " │ ")
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

// renderSearch renders the search input view
func (m Model) renderSearch() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("🔍 Search Pokemon"))
	b.WriteString("\n\n")

	b.WriteString("Enter Pokemon name or ID: ")
	b.WriteString(m.searchQuery)
	b.WriteString("_\n")

	// Show search preview
	if m.searchQuery == "" {
		// Show helpful prompt when empty
		b.WriteString(helpStyle.Render("\nTry: pikachu, 25, or char"))
	} else {
		// Show match count and preview
		if m.searchCount == 0 {
			b.WriteString("\n")
			b.WriteString(errorStyle.Render("No matches found"))
		} else if m.searchCount == 1 && len(m.searchMatches) == 1 {
			// Single match - show preview with ID lookup
			pokemon := m.searchMatches[0]
			name := cases.Title(language.English).String(pokemon.Name)
			b.WriteString("\n")
			b.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("Found: %s", name)))
			b.WriteString(helpStyle.Render(" (press Enter)"))
		} else {
			// Multiple matches - show count and preview list
			b.WriteString("\n")
			b.WriteString(selectedMenuItemStyle.Render(fmt.Sprintf("Found %d matches:", m.searchCount)))
			b.WriteString("\n")
			for i, pokemon := range m.searchMatches {
				name := cases.Title(language.English).String(pokemon.Name)
				b.WriteString(menuItemStyle.Render(fmt.Sprintf("  • %s", name)))
				if i < len(m.searchMatches)-1 {
					b.WriteString("\n")
				}
			}
			if m.searchCount > len(m.searchMatches) {
				b.WriteString("\n")
				b.WriteString(helpStyle.Render(fmt.Sprintf("  ... and %d more", m.searchCount-len(m.searchMatches))))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("\nPress Enter to search, Esc to go back"))

	return b.String()
}

// renderList renders the Pokemon list view
func (m Model) renderList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("📋 Pokemon List"))
	b.WriteString("\n\n")

	if len(m.pokemonList) == 0 {
		b.WriteString("No Pokemon found.\n")
	} else {
		// Show current page info
		start := m.listOffset + 1
		end := m.listOffset + len(m.pokemonList)
		fmt.Fprintf(&b, "Showing %d-%d of %d\n\n", start, end, m.totalCount)

		// Render list items
		visibleStart := 0
		visibleEnd := len(m.pokemonList)

		// If terminal is small, show a window around cursor
		maxVisible := m.height - 10
		if maxVisible < len(m.pokemonList) {
			halfWindow := maxVisible / 2
			visibleStart = m.cursor - halfWindow
			visibleEnd = m.cursor + halfWindow

			if visibleStart < 0 {
				visibleStart = 0
				visibleEnd = maxVisible
			}
			if visibleEnd > len(m.pokemonList) {
				visibleEnd = len(m.pokemonList)
				visibleStart = visibleEnd - maxVisible
				if visibleStart < 0 {
					visibleStart = 0
				}
			}
		}

		for i := visibleStart; i < visibleEnd; i++ {
			pokemon := m.pokemonList[i]
			line := fmt.Sprintf("%d. %s", m.listOffset+i+1, cases.Title(language.English).String(pokemon.Name))

			if i == m.cursor {
				b.WriteString(selectedMenuItemStyle.Render("▸ " + line))
			} else {
				b.WriteString(menuItemStyle.Render("  " + line))
			}
			b.WriteString("\n")
		}
	}

	// Build help text with contextual pagination hints
	helpText := "\n↑/k: up, ↓/j: down, Enter: view details"

	// Show appropriate pagination hints
	canGoPrev := m.listOffset >= m.listLimit
	canGoNext := m.listOffset+m.listLimit < m.totalCount

	if canGoPrev && canGoNext {
		helpText += ", n: next page, p: prev page"
	} else if canGoPrev {
		helpText += ", p: prev page (at end)"
	} else if canGoNext {
		helpText += ", n: next page"
	} else {
		// Only one page of results
		helpText += " (single page)"
	}

	helpText += ", q/Esc: back"
	b.WriteString(helpStyle.Render(helpText))

	return b.String()
}

// renderDetails renders the Pokemon details view
func (m Model) renderDetails() string {
	var b strings.Builder

	if m.pokemonDetails == nil {
		return "No Pokemon loaded."
	}

	pokemon, ok := m.pokemonDetails.(*models.Pokemon)
	if !ok || pokemon == nil {
		return "Invalid Pokemon data."
	}

	// Title
	name := cases.Title(language.English).String(pokemon.Name)
	b.WriteString(pokemonNameStyle.Render(fmt.Sprintf("#%d %s", pokemon.ID, name)))
	b.WriteString("\n\n")

	// Display ASCII sprite if available (with defensive error handling)
	// Only show sprite section if HasSprite check passes
	if pokemon.HasSprite() {
		asciiArt := pokemon.GetSpriteASCII()
		if asciiArt != "" {
			b.WriteString("\n")
			b.WriteString(asciiArt)
			b.WriteString("\n\n")
		}
	}

	// Display official artwork URL if available
	bestSpriteURL := pokemon.BestSprite()
	if bestSpriteURL != "" {
		urlStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			MarginLeft(2)

		b.WriteString(urlStyle.Render("🖼️  Official Artwork: " + bestSpriteURL))
		b.WriteString("\n\n")
	}

	// Types
	if len(pokemon.Types) > 0 {
		b.WriteString("Types: ")
		typeNames := make([]string, 0, len(pokemon.Types))
		for _, t := range pokemon.Types {
			if t != nil {
				typeNames = append(typeNames, cases.Title(language.English).String(t.Name))
			}
		}
		if len(typeNames) > 0 {
			b.WriteString(strings.Join(typeNames, ", "))
			b.WriteString("\n")
		}
	}

	// Physical characteristics
	b.WriteString(sectionStyle.Render(fmt.Sprintf("Height: %s", pokemon.FormattedHeight())))
	b.WriteString("\n")
	fmt.Fprintf(&b, "Weight: %s\n", pokemon.FormattedWeight())

	if pokemon.BaseExperience > 0 {
		fmt.Fprintf(&b, "Base Experience: %d\n", pokemon.BaseExperience)
	}

	// Abilities
	if len(pokemon.Abilities) > 0 {
		b.WriteString(sectionStyle.Render("\nAbilities:"))
		b.WriteString("\n")
		for _, ability := range pokemon.Abilities {
			name := cases.Title(language.English).String(strings.ReplaceAll(ability.Name, "-", " "))
			if ability.IsHidden {
				name += " (Hidden)"
			}
			b.WriteString("  • " + name + "\n")
		}
	}

	// Stats
	if len(pokemon.Stats) > 0 {
		b.WriteString(sectionStyle.Render("\nBase Stats:"))
		b.WriteString("\n")
		for _, stat := range pokemon.Stats {
			name := cases.Title(language.English).String(strings.ReplaceAll(stat.Name, "-", " "))

			// Create a simple stat bar
			barLength := 20
			filledLength := int(float64(stat.BaseValue) / 255.0 * float64(barLength))
			if filledLength > barLength {
				filledLength = barLength
			}
			bar := strings.Repeat("█", filledLength) + strings.Repeat("░", barLength-filledLength)

			fmt.Fprintf(&b, "  %-18s %3d %s\n", name+":", stat.BaseValue, statBarStyle.Render(bar))
		}
	}

	b.WriteString(helpStyle.Render("\nPress Backspace/Esc/q to go back"))

	return b.String()
}

// renderLoading renders the loading view
func (m Model) renderLoading() string {
	message := m.loadingMessage
	if message == "" {
		message = "Loading..."
	}
	return titleStyle.Render(m.spinner.View()+" "+message) + "\n\nPlease wait..."
}

// errorType represents different categories of errors
type errorType int

const (
	errorTypeNotFound errorType = iota
	errorTypeNetwork
	errorTypeAPI
	errorTypeValidation
	errorTypeGeneric
)

// detectErrorType categorizes an error based on its message content
func detectErrorType(errorMsg string) errorType {
	lowerMsg := strings.ToLower(errorMsg)

	// Check for NotFoundError indicators
	if strings.Contains(lowerMsg, "not found") ||
		strings.Contains(lowerMsg, "no pokemon") ||
		strings.Contains(lowerMsg, "404") {
		return errorTypeNotFound
	}

	// Check for NetworkError indicators
	if strings.Contains(lowerMsg, "connection") ||
		strings.Contains(lowerMsg, "timeout") ||
		strings.Contains(lowerMsg, "network") {
		return errorTypeNetwork
	}

	// Check for APIError indicators
	if strings.Contains(lowerMsg, "api") ||
		strings.Contains(lowerMsg, "failed to fetch") ||
		strings.Contains(lowerMsg, "status code") {
		return errorTypeAPI
	}

	// Check for ValidationError indicators
	if strings.Contains(lowerMsg, "invalid") ||
		strings.Contains(lowerMsg, "empty") ||
		strings.Contains(lowerMsg, "cannot be") {
		return errorTypeValidation
	}

	return errorTypeGeneric
}

// getErrorDisplay returns the title, style, and suggestion for an error type
func getErrorDisplay(errType errorType) (title string, style lipgloss.Style, suggestion string) {
	switch errType {
	case errorTypeNotFound:
		return "❌ Pokemon Not Found",
			notFoundErrorStyle,
			"💡 Try checking the spelling or ID number"
	case errorTypeNetwork:
		return "⚠️  Connection Issue",
			networkErrorStyle,
			"💡 Check your internet connection and try again"
	case errorTypeAPI:
		return "🚫 API Error",
			apiErrorStyle,
			"💡 PokeAPI may be temporarily unavailable. Please try again later"
	case errorTypeValidation:
		return "⚡ Invalid Input",
			validationErrorStyle,
			"💡 Pokemon name or ID must not be empty"
	default:
		return "❌ Error",
			errorStyle,
			""
	}
}

// renderError renders the error view with categorized styling and helpful suggestions
func (m Model) renderError() string {
	var b strings.Builder

	// Detect error type and get display elements
	errType := detectErrorType(m.errorMsg)
	title, style, suggestion := getErrorDisplay(errType)

	// Render title
	b.WriteString(titleStyle.Render(title))
	b.WriteString("\n\n")

	// Render error message with appropriate styling
	b.WriteString(style.Render(m.errorMsg))
	b.WriteString("\n")

	// Render suggestion if available
	if suggestion != "" {
		b.WriteString("\n")
		b.WriteString(errorSuggestionStyle.Render(suggestion))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press any key to return to menu"))

	return b.String()
}
