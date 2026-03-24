// Package service provides a service layer that coordinates the API client and cache.
// It implements cache-first fetching with fallback to the API.
package service

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pfjndev/pokedex/internal/api"
	"github.com/pfjndev/pokedex/internal/models"
)

// APIClient defines the interface for API operations
type APIClient interface {
	GetPokemon(nameOrID string) (*api.PokemonResponse, error)
	ListPokemon(limit, offset int) (*api.PokemonListResponse, error)
}

// Cache defines the interface for cache operations
type Cache interface {
	Get(nameOrID string) ([]byte, error)
	Set(nameOrID string, data []byte) error
}

// Service coordinates API and cache operations for Pokemon data.
type Service struct {
	apiClient APIClient
	cache     Cache
}

// NewService creates a new Service with the given API client and cache.
func NewService(apiClient APIClient, cache Cache) *Service {
	return &Service{
		apiClient: apiClient,
		cache:     cache,
	}
}

// GetPokemon fetches a Pokemon by name or ID using cache-first strategy.
// It first checks the cache, and if not found, fetches from the API and caches the result.
// Returns a models.Pokemon or an error if the Pokemon is not found or a failure occurs.
//
// Example:
//
//	pokemon, err := service.GetPokemon("pikachu")
//	pokemon, err := service.GetPokemon("25")
func (s *Service) GetPokemon(nameOrID string) (*models.Pokemon, error) {
	nameOrID = strings.ToLower(strings.TrimSpace(nameOrID))
	if nameOrID == "" {
		return nil, fmt.Errorf("pokemon name or ID cannot be empty")
	}

	// Try cache first
	cachedData, err := s.cache.Get(nameOrID)
	if err == nil {
		var pokemon models.Pokemon
		if err := json.Unmarshal(cachedData, &pokemon); err != nil {
			return nil, fmt.Errorf("failed to unmarshal cached pokemon: %w", err)
		}
		return &pokemon, nil
	}

	// Cache miss, fetch from API
	apiResp, err := s.apiClient.GetPokemon(nameOrID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon from API: %w", err)
	}

	// Transform API response to domain model
	pokemon := transformPokemonResponse(apiResp)

	// Cache the result
	pokemonData, err := json.Marshal(pokemon)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pokemon for caching: %w", err)
	}

	if err := s.cache.Set(nameOrID, pokemonData); err != nil {
		// Log error but don't fail the request
		// The user still got their Pokemon data
		_ = err
	}

	return pokemon, nil
}

// ListPokemon fetches a paginated list of Pokemon from the API.
// It returns basic Pokemon information (ID and Name) for browsing.
// Use limit and offset for pagination.
//
// Example:
//
//	list, err := service.ListPokemon(20, 0)   // First 20 Pokemon
//	list, err := service.ListPokemon(20, 20)  // Next 20 Pokemon
func (s *Service) ListPokemon(limit, offset int) (*PokemonList, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	apiResp, err := s.apiClient.ListPokemon(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list pokemon from API: %w", err)
	}

	items := make([]*PokemonListItem, 0, len(apiResp.Results))
	for _, result := range apiResp.Results {
		id := extractPokemonIDFromURL(result.URL)
		items = append(items, &PokemonListItem{
			ID:   id,
			Name: result.Name,
		})
	}

	return &PokemonList{
		Count:    apiResp.Count,
		Next:     apiResp.Next,
		Previous: apiResp.Previous,
		Items:    items,
	}, nil
}

// SearchPokemon searches for Pokemon by name or ID.
// If the query is numeric, it fetches the Pokemon by ID directly.
// Otherwise, it fetches the list from the API and filters by name (contains match).
// Returns matching Pokemon names.
//
// Example:
//
//	results, err := service.SearchPokemon("pika")  // Returns pikachu, pikachu-gmax, etc.
//	results, err := service.SearchPokemon("zard")  // Returns charizard, charizard-mega-x, etc.
//	results, err := service.SearchPokemon("25")    // Returns pikachu (ID 25)
func (s *Service) SearchPokemon(query string) ([]*PokemonListItem, error) {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Check if query is numeric (Pokemon ID)
	if _, err := strconv.Atoi(query); err == nil {
		// Query is numeric, fetch by ID
		pokemon, err := s.GetPokemon(query)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch pokemon by ID: %w", err)
		}

		// Return as single-item slice
		return []*PokemonListItem{
			{
				ID:   pokemon.ID,
				Name: pokemon.Name,
			},
		}, nil
	}

	// Query is not numeric, search by name prefix
	// Fetch a large list to search through
	// PokeAPI has ~1000 Pokemon, so we fetch them all
	apiResp, err := s.apiClient.ListPokemon(10000, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon list for search: %w", err)
	}

	// Filter by name (contains match)
	matches := make([]*PokemonListItem, 0)
	for _, result := range apiResp.Results {
		if strings.Contains(result.Name, query) {
			id := extractPokemonIDFromURL(result.URL)
			matches = append(matches, &PokemonListItem{
				ID:   id,
				Name: result.Name,
			})
		}
	}

	return matches, nil
}

// FilterByType filters a list of Pokemon by type.
// Note: This requires fetching each Pokemon individually to check its types.
// For MVP, this is kept simple but could be optimized with a type index.
func (s *Service) FilterByType(pokemonList []*PokemonListItem, typeName string) ([]*PokemonListItem, error) {
	typeName = strings.ToLower(strings.TrimSpace(typeName))
	if typeName == "" {
		return pokemonList, nil
	}

	filtered := make([]*PokemonListItem, 0)
	for _, item := range pokemonList {
		pokemon, err := s.GetPokemon(item.Name)
		if err != nil {
			continue
		}

		hasType := false
		for _, t := range pokemon.Types {
			if strings.ToLower(t.Name) == typeName {
				hasType = true
				break
			}
		}

		if hasType {
			filtered = append(filtered, item)
		}
	}

	return filtered, nil
}

// PokemonList represents a paginated list of Pokemon.
type PokemonList struct {
	Count    int                `json:"count"`
	Next     *string            `json:"next"`
	Previous *string            `json:"previous"`
	Items    []*PokemonListItem `json:"items"`
}

// PokemonListItem represents a single Pokemon in a list.
type PokemonListItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// extractPokemonIDFromURL extracts the Pokemon ID from a PokeAPI URL.
// Example: "https://pokeapi.co/api/v2/pokemon/25/" -> 25
func extractPokemonIDFromURL(url string) int {
	// Remove trailing slash
	url = strings.TrimSuffix(url, "/")
	
	// Split by "/" and get the last segment
	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return 0
	}
	
	// Parse the ID from the last segment
	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0
	}
	
	return id
}

// transformPokemonResponse converts an API response to a domain model.
func transformPokemonResponse(resp *api.PokemonResponse) *models.Pokemon {
	pokemon := &models.Pokemon{
		ID:             resp.ID,
		Name:           resp.Name,
		Height:         resp.Height,
		Weight:         resp.Weight,
		BaseExperience: resp.BaseExperience,
		IsDefault:      resp.IsDefault,
	}

	// Transform types
	pokemon.Types = make([]*models.Type, 0, len(resp.Types))
	for _, t := range resp.Types {
		color := models.TypeColors[t.Type.Name]
		if color == "" {
			color = "#A8A878" // Default to normal type color
		}
		pokemon.Types = append(pokemon.Types, &models.Type{
			Name:  t.Type.Name,
			Color: color,
		})
	}

	// Transform abilities
	pokemon.Abilities = make([]*models.Ability, 0, len(resp.Abilities))
	for _, a := range resp.Abilities {
		pokemon.Abilities = append(pokemon.Abilities, &models.Ability{
			Name:     a.Ability.Name,
			IsHidden: a.IsHidden,
			Slot:     a.Slot,
		})
	}

	// Transform stats
	pokemon.Stats = make([]*models.Stat, 0, len(resp.Stats))
	for _, s := range resp.Stats {
		pokemon.Stats = append(pokemon.Stats, &models.Stat{
			Name:      s.Stat.Name,
			BaseValue: s.BaseStat,
		})
	}

	// Transform sprites
	if resp.Sprites.FrontDefault != "" || resp.Sprites.FrontShiny != "" || resp.Sprites.OfficialArtwork.Artwork.FrontDefault != "" {
		pokemon.Sprites = &models.Sprite{
			FrontDefault:     resp.Sprites.FrontDefault,
			FrontFemale:      resp.Sprites.FrontFemale,
			FrontShiny:       resp.Sprites.FrontShiny,
			FrontShinyFemale: resp.Sprites.FrontShinyFemale,
			BackDefault:      resp.Sprites.BackDefault,
			BackFemale:       resp.Sprites.BackFemale,
			BackShiny:        resp.Sprites.BackShiny,
			BackShinyFemale:  resp.Sprites.BackShinyFemale,
			OfficialArtwork:  resp.Sprites.OfficialArtwork.Artwork.FrontDefault,
		}
	}

	return pokemon
}
