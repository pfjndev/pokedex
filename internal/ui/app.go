// Package ui provides the terminal user interface for the Pokedex.
package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pfjndev/pokedex/internal/service"
)

// AppState represents the current state of the application
type AppState int

const (
	// StateMenu shows the main menu
	StateMenu AppState = iota
	// StateSearch shows the search input
	StateSearch
	// StateList shows the Pokemon list
	StateList
	// StateDetails shows Pokemon details
	StateDetails
	// StateLoading shows a loading spinner
	StateLoading
	// StateError shows an error message
	StateError
)

// Model represents the Bubble Tea application model
type Model struct {
	service *service.Service
	state   AppState
	
	// Menu state
	menuCursor int
	
	// Search state
	searchQuery   string
	searchMatches []*service.PokemonListItem
	searchCount   int
	
	// List state
	pokemonList []*service.PokemonListItem
	listOffset  int
	listLimit   int
	totalCount  int
	cursor      int
	
	// Details state
	currentPokemon *service.PokemonListItem
	pokemonDetails interface{}
	
	// Error state
	errorMsg string
	
	// Loading state
	spinner        spinner.Model
	loadingMessage string
	
	// UI dimensions
	width  int
	height int
}

// NewModel creates a new Model with the given service
func NewModel(svc *service.Service) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	
	return Model{
		service:    svc,
		state:      StateMenu,
		menuCursor: 0,
		listLimit:  20,
		listOffset: 0,
		cursor:     0,
		spinner:    s,
	}
}

// Init initializes the model (required by Bubble Tea)
func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

// Msg types for the application

// pokemonListMsg is sent when a Pokemon list is fetched
type pokemonListMsg struct {
	list       []*service.PokemonListItem
	totalCount int
	err        error
}

// pokemonDetailsMsg is sent when Pokemon details are fetched
type pokemonDetailsMsg struct {
	pokemon interface{}
	err     error
}

// searchResultsMsg is sent when search results are ready
type searchResultsMsg struct {
	results []*service.PokemonListItem
	err     error
}

// searchPreviewMsg is sent when live search preview is ready
type searchPreviewMsg struct {
	matches []*service.PokemonListItem
	count   int
}

// Update handles messages and updates the model (required by Bubble Tea)
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case spinner.TickMsg:
		if m.state == StateLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
		return m, nil

	case pokemonListMsg:
		return m.handlePokemonList(msg)

	case pokemonDetailsMsg:
		return m.handlePokemonDetails(msg)

	case searchResultsMsg:
		return m.handleSearchResults(msg)

	case searchPreviewMsg:
		return m.handleSearchPreview(msg)
	}

	return m, nil
}

// handleKeyPress handles keyboard input based on current state
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit keys
	if msg.Type == tea.KeyCtrlC || (m.state == StateMenu && msg.String() == "q") {
		return m, tea.Quit
	}

	switch m.state {
	case StateMenu:
		return m.handleMenuKeys(msg)
	case StateSearch:
		return m.handleSearchKeys(msg)
	case StateList:
		return m.handleListKeys(msg)
	case StateDetails:
		return m.handleDetailsKeys(msg)
	case StateError:
		return m.handleErrorKeys(msg)
	}

	return m, nil
}

// handleMenuKeys handles keyboard input in the menu state
func (m Model) handleMenuKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		} else {
			m.menuCursor = 2
		}
		return m, nil
	case "down", "j":
		if m.menuCursor < 2 {
			m.menuCursor++
		} else {
			m.menuCursor = 0
		}
		return m, nil
	case "enter":
		switch m.menuCursor {
		case 0:
			m.state = StateSearch
			m.searchQuery = ""
			return m, nil
		case 1:
			m.state = StateLoading
			m.loadingMessage = "Fetching Pokemon list..."
			return m, m.fetchPokemonList()
		case 2:
			return m, tea.Quit
		}
		return m, nil
	case "s":
		m.state = StateSearch
		m.searchQuery = ""
		return m, nil
	case "b":
		m.state = StateLoading
		m.loadingMessage = "Fetching Pokemon list..."
		return m, m.fetchPokemonList()
	case "q":
		return m, tea.Quit
	}
	return m, nil
}

// handleSearchKeys handles keyboard input in the search state
func (m Model) handleSearchKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		m.state = StateMenu
		m.searchQuery = ""
		m.searchMatches = nil
		m.searchCount = 0
		return m, nil
	case tea.KeyEnter:
		if m.searchQuery != "" {
			// If we have exactly one match, go directly to it
			if len(m.searchMatches) == 1 {
				m.currentPokemon = m.searchMatches[0]
				m.state = StateLoading
				m.loadingMessage = fmt.Sprintf("Loading %s...", m.searchMatches[0].Name)
				return m, m.fetchPokemonDetails(m.searchMatches[0].Name)
			}
			// Multiple or no matches, trigger full search
			m.state = StateLoading
			m.loadingMessage = "Searching..."
			return m, m.fetchPokemonByName(m.searchQuery)
		}
		return m, nil
	case tea.KeyBackspace:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			// Update live search preview
			return m, m.updateSearchPreview()
		}
		return m, nil
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
			// Update live search preview
			return m, m.updateSearchPreview()
		}
	}
	return m, nil
}

// handleListKeys handles keyboard input in the list state
func (m Model) handleListKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc":
		m.state = StateMenu
		return m, nil
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
		return m, nil
	case "down", "j":
		if m.cursor < len(m.pokemonList)-1 {
			m.cursor++
		}
		return m, nil
	case "enter":
		if len(m.pokemonList) > 0 && m.cursor < len(m.pokemonList) {
			m.currentPokemon = m.pokemonList[m.cursor]
			m.state = StateLoading
			m.loadingMessage = fmt.Sprintf("Loading %s...", m.currentPokemon.Name)
			return m, m.fetchPokemonDetails(m.currentPokemon.Name)
		}
		return m, nil
	case "n":
		// Next page
		if m.listOffset+m.listLimit < m.totalCount {
			m.listOffset += m.listLimit
			m.cursor = 0
			m.state = StateLoading
			m.loadingMessage = "Fetching Pokemon list..."
			return m, m.fetchPokemonList()
		}
		return m, nil
	case "p":
		// Previous page
		if m.listOffset >= m.listLimit {
			m.listOffset -= m.listLimit
			m.cursor = 0
			m.state = StateLoading
			m.loadingMessage = "Fetching Pokemon list..."
			return m, m.fetchPokemonList()
		}
		return m, nil
	}
	return m, nil
}

// handleDetailsKeys handles keyboard input in the details state
func (m Model) handleDetailsKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "backspace":
		m.state = StateList
		return m, nil
	}
	return m, nil
}

// handleErrorKeys handles keyboard input in the error state
func (m Model) handleErrorKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "esc", "enter":
		m.state = StateMenu
		return m, nil
	}
	return m, nil
}

// handlePokemonList handles the Pokemon list message
func (m Model) handlePokemonList(msg pokemonListMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.state = StateError
		m.errorMsg = fmt.Sprintf("Failed to fetch Pokemon list: %v", msg.err)
		return m, nil
	}

	m.pokemonList = msg.list
	m.totalCount = msg.totalCount
	m.cursor = 0
	m.state = StateList
	return m, nil
}

// handlePokemonDetails handles the Pokemon details message
func (m Model) handlePokemonDetails(msg pokemonDetailsMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.state = StateError
		m.errorMsg = fmt.Sprintf("Failed to fetch Pokemon details: %v", msg.err)
		return m, nil
	}

	m.pokemonDetails = msg.pokemon
	m.state = StateDetails
	return m, nil
}

// handleSearchResults handles the search results message
func (m Model) handleSearchResults(msg searchResultsMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.state = StateError
		m.errorMsg = fmt.Sprintf("Search failed: %v", msg.err)
		return m, nil
	}

	if len(msg.results) == 0 {
		m.state = StateError
		m.errorMsg = "No Pokemon found matching your search"
		return m, nil
	}

	if len(msg.results) == 1 {
		// Direct match, show details
		m.currentPokemon = msg.results[0]
		m.state = StateLoading
		m.loadingMessage = fmt.Sprintf("Loading %s...", msg.results[0].Name)
		return m, m.fetchPokemonDetails(msg.results[0].Name)
	}

	// Multiple matches, show list
	m.pokemonList = msg.results
	m.cursor = 0
	m.state = StateList
	return m, nil
}

// View renders the UI (required by Bubble Tea)
func (m Model) View() string {
	switch m.state {
	case StateMenu:
		return m.renderMenu()
	case StateSearch:
		return m.renderSearch()
	case StateList:
		return m.renderList()
	case StateDetails:
		return m.renderDetails()
	case StateLoading:
		return m.renderLoading()
	case StateError:
		return m.renderError()
	default:
		return "Unknown state"
	}
}

// fetchPokemonList fetches a paginated list of Pokemon
func (m Model) fetchPokemonList() tea.Cmd {
	return func() tea.Msg {
		list, err := m.service.ListPokemon(m.listLimit, m.listOffset)
		if err != nil {
			return pokemonListMsg{err: err}
		}
		return pokemonListMsg{list: list.Items, totalCount: list.Count}
	}
}

// fetchPokemonDetails fetches details for a specific Pokemon
func (m Model) fetchPokemonDetails(nameOrID string) tea.Cmd {
	return func() tea.Msg {
		pokemon, err := m.service.GetPokemon(nameOrID)
		if err != nil {
			return pokemonDetailsMsg{err: err}
		}
		return pokemonDetailsMsg{pokemon: pokemon}
	}
}

// fetchPokemonByName searches for a Pokemon by name
func (m Model) fetchPokemonByName(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.service.SearchPokemon(query)
		if err != nil {
			return searchResultsMsg{err: err}
		}
		return searchResultsMsg{results: results}
	}
}

// updateSearchPreview performs a live search preview as the user types
func (m Model) updateSearchPreview() tea.Cmd {
	// Don't search if query is empty
	if m.searchQuery == "" {
		return func() tea.Msg {
			return searchPreviewMsg{matches: nil, count: 0}
		}
	}

	query := m.searchQuery
	return func() tea.Msg {
		results, err := m.service.SearchPokemon(query)
		if err != nil {
			// Silent failure for preview
			return searchPreviewMsg{matches: nil, count: 0}
		}
		// Return first 5 matches for preview
		preview := results
		if len(preview) > 5 {
			preview = results[:5]
		}
		return searchPreviewMsg{matches: preview, count: len(results)}
	}
}

// handleSearchPreview handles the search preview message
func (m Model) handleSearchPreview(msg searchPreviewMsg) (tea.Model, tea.Cmd) {
	m.searchMatches = msg.matches
	m.searchCount = msg.count
	return m, nil
}
