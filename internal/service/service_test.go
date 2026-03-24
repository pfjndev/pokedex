package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pfjndev/pokedex/internal/api"
	"github.com/pfjndev/pokedex/internal/models"
)

// mockAPIClient is a mock implementation of the API client for testing
type mockAPIClient struct {
	getPokemonFunc func(nameOrID string) (*api.PokemonResponse, error)
	listPokemonFunc func(limit, offset int) (*api.PokemonListResponse, error)
}

func (m *mockAPIClient) GetPokemon(nameOrID string) (*api.PokemonResponse, error) {
	if m.getPokemonFunc != nil {
		return m.getPokemonFunc(nameOrID)
	}
	return nil, fmt.Errorf("mock not configured")
}

func (m *mockAPIClient) ListPokemon(limit, offset int) (*api.PokemonListResponse, error) {
	if m.listPokemonFunc != nil {
		return m.listPokemonFunc(limit, offset)
	}
	return nil, fmt.Errorf("mock not configured")
}

// mockCache is a mock implementation of the cache for testing
type mockCache struct {
	data map[string][]byte
}

func newMockCache() *mockCache {
	return &mockCache{
		data: make(map[string][]byte),
	}
}

func (m *mockCache) Get(nameOrID string) ([]byte, error) {
	if data, ok := m.data[nameOrID]; ok {
		return data, nil
	}
	return nil, fmt.Errorf("pokemon %q not found in cache", nameOrID)
}

func (m *mockCache) Set(nameOrID string, data []byte) error {
	m.data[nameOrID] = data
	return nil
}

// TestNewService verifies that NewService creates a service correctly
func TestNewService(t *testing.T) {
	apiClient := &mockAPIClient{}
	cache := newMockCache()

	service := NewService(nil, nil)
	if service == nil {
		t.Fatal("NewService returned nil")
	}

	service = NewService(apiClient, cache)
	if service == nil {
		t.Fatal("NewService returned nil")
	}

	if service.apiClient != apiClient {
		t.Error("apiClient not set correctly")
	}
}

// TestGetPokemonCacheHit verifies cache-first behavior when Pokemon is cached
func TestGetPokemonCacheHit(t *testing.T) {
	cache := newMockCache()
	apiClient := &mockAPIClient{
		getPokemonFunc: func(nameOrID string) (*api.PokemonResponse, error) {
			t.Error("API should not be called when cache hits")
			return nil, fmt.Errorf("should not be called")
		},
	}

	// Pre-populate cache
	cachedPokemon := &models.Pokemon{
		ID:   25,
		Name: "pikachu",
		Types: []*models.Type{
			{Name: "electric", Color: "#F8D030"},
		},
	}
	cachedData, _ := json.Marshal(cachedPokemon)
	_ = cache.Set("pikachu", cachedData)

	service := NewService(apiClient, cache)
	pokemon, err := service.GetPokemon("pikachu")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if pokemon.Name != "pikachu" {
		t.Errorf("expected name pikachu, got %s", pokemon.Name)
	}

	if pokemon.ID != 25 {
		t.Errorf("expected ID 25, got %d", pokemon.ID)
	}
}

// TestGetPokemonCacheMiss verifies API fallback when Pokemon is not cached
func TestGetPokemonCacheMiss(t *testing.T) {
	cache := newMockCache()
	apiCalled := false

	apiClient := &mockAPIClient{
		getPokemonFunc: func(nameOrID string) (*api.PokemonResponse, error) {
			apiCalled = true
			return &api.PokemonResponse{
				ID:     25,
				Name:   "pikachu",
				Height: 4,
				Weight: 60,
				Types: []api.PokemonType{
					{Type: api.NamedResource{Name: "electric"}},
				},
				Abilities: []api.PokemonAbility{
					{Ability: api.NamedResource{Name: "static"}, IsHidden: false, Slot: 1},
				},
				Stats: []api.PokemonStat{
					{Stat: api.NamedResource{Name: "hp"}, BaseStat: 35, Effort: 0},
				},
			}, nil
		},
	}

	service := NewService(apiClient, cache)
	pokemon, err := service.GetPokemon("pikachu")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !apiCalled {
		t.Error("API should have been called on cache miss")
	}

	if pokemon.Name != "pikachu" {
		t.Errorf("expected name pikachu, got %s", pokemon.Name)
	}

	if pokemon.ID != 25 {
		t.Errorf("expected ID 25, got %d", pokemon.ID)
	}

	// Verify it was cached
	cachedData, err := cache.Get("pikachu")
	if err != nil {
		t.Error("Pokemon should have been cached after API fetch")
	}

	var cachedPokemon models.Pokemon
	if err := json.Unmarshal(cachedData, &cachedPokemon); err != nil {
		t.Fatalf("failed to unmarshal cached pokemon: %v", err)
	}

	if cachedPokemon.Name != "pikachu" {
		t.Error("cached pokemon name mismatch")
	}
}

// TestGetPokemonAPIError verifies error propagation when API fails
func TestGetPokemonAPIError(t *testing.T) {
	cache := newMockCache()
	expectedErr := &api.APIError{
		StatusCode: 404,
		Message:    "Pokemon not found",
	}

	apiClient := &mockAPIClient{
		getPokemonFunc: func(nameOrID string) (*api.PokemonResponse, error) {
			return nil, expectedErr
		},
	}

	service := NewService(apiClient, cache)
	_, err := service.GetPokemon("invalid-pokemon")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestGetPokemonEmptyName verifies validation of empty names
func TestGetPokemonEmptyName(t *testing.T) {
	service := NewService(&mockAPIClient{}, newMockCache())

	_, err := service.GetPokemon("")
	if err == nil {
		t.Fatal("expected error for empty name")
	}

	_, err = service.GetPokemon("   ")
	if err == nil {
		t.Fatal("expected error for whitespace name")
	}
}

// TestListPokemon verifies pagination and list fetching
func TestListPokemon(t *testing.T) {
	apiClient := &mockAPIClient{
		listPokemonFunc: func(limit, offset int) (*api.PokemonListResponse, error) {
			return &api.PokemonListResponse{
				Count: 1000,
				Next:  strPtr("https://pokeapi.co/api/v2/pokemon?limit=20&offset=20"),
				Results: []api.PokemonListItem{
					{Name: "bulbasaur", URL: "https://pokeapi.co/api/v2/pokemon/1/"},
					{Name: "ivysaur", URL: "https://pokeapi.co/api/v2/pokemon/2/"},
					{Name: "venusaur", URL: "https://pokeapi.co/api/v2/pokemon/3/"},
				},
			}, nil
		},
	}

	service := NewService(apiClient, newMockCache())
	list, err := service.ListPokemon(20, 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if list.Count != 1000 {
		t.Errorf("expected count 1000, got %d", list.Count)
	}

	if len(list.Items) != 3 {
		t.Errorf("expected 3 items, got %d", len(list.Items))
	}

	if list.Items[0].Name != "bulbasaur" {
		t.Errorf("expected first item to be bulbasaur, got %s", list.Items[0].Name)
	}

	if list.Items[0].ID != 1 {
		t.Errorf("expected first item ID to be 1, got %d", list.Items[0].ID)
	}

	if list.Items[1].ID != 2 {
		t.Errorf("expected second item ID to be 2, got %d", list.Items[1].ID)
	}

	if list.Items[2].ID != 3 {
		t.Errorf("expected third item ID to be 3, got %d", list.Items[2].ID)
	}
}

// TestListPokemonDefaultLimit verifies default limit behavior
func TestListPokemonDefaultLimit(t *testing.T) {
	actualLimit := 0
	apiClient := &mockAPIClient{
		listPokemonFunc: func(limit, offset int) (*api.PokemonListResponse, error) {
			actualLimit = limit
			return &api.PokemonListResponse{
				Count:   100,
				Results: []api.PokemonListItem{},
			}, nil
		},
	}

	service := NewService(apiClient, newMockCache())
	_, err := service.ListPokemon(0, 0)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actualLimit != 20 {
		t.Errorf("expected default limit 20, got %d", actualLimit)
	}
}

// TestSearchPokemon verifies name contains search
func TestSearchPokemon(t *testing.T) {
	apiClient := &mockAPIClient{
		listPokemonFunc: func(limit, offset int) (*api.PokemonListResponse, error) {
			return &api.PokemonListResponse{
				Count: 5,
				Results: []api.PokemonListItem{
					{Name: "pikachu"},
					{Name: "pikachu-gmax"},
					{Name: "pichu"},
					{Name: "bulbasaur"},
					{Name: "charmander"},
				},
			}, nil
		},
	}

	service := NewService(apiClient, newMockCache())
	results, err := service.SearchPokemon("pika")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}

	if results[0].Name != "pikachu" {
		t.Errorf("expected first result pikachu, got %s", results[0].Name)
	}

	if results[1].Name != "pikachu-gmax" {
		t.Errorf("expected second result pikachu-gmax, got %s", results[1].Name)
	}
}

// TestSearchPokemonEmptyQuery verifies validation of empty search queries
func TestSearchPokemonEmptyQuery(t *testing.T) {
	service := NewService(&mockAPIClient{}, newMockCache())

	_, err := service.SearchPokemon("")
	if err == nil {
		t.Fatal("expected error for empty query")
	}

	_, err = service.SearchPokemon("   ")
	if err == nil {
		t.Fatal("expected error for whitespace query")
	}
}

// TestSearchPokemonByID verifies search by numeric ID
func TestSearchPokemonByID(t *testing.T) {
	apiClient := &mockAPIClient{
		getPokemonFunc: func(nameOrID string) (*api.PokemonResponse, error) {
			// Simulate PokeAPI behavior for different IDs
			switch nameOrID {
			case "25":
				return &api.PokemonResponse{
					ID:     25,
					Name:   "pikachu",
					Height: 4,
					Weight: 60,
					Types: []api.PokemonType{
						{Type: api.NamedResource{Name: "electric"}},
					},
				}, nil
			case "1":
				return &api.PokemonResponse{
					ID:   1,
					Name: "bulbasaur",
					Types: []api.PokemonType{
						{Type: api.NamedResource{Name: "grass"}},
						{Type: api.NamedResource{Name: "poison"}},
					},
				}, nil
			case "150":
				return &api.PokemonResponse{
					ID:   150,
					Name: "mewtwo",
					Types: []api.PokemonType{
						{Type: api.NamedResource{Name: "psychic"}},
					},
				}, nil
			default:
				return nil, &api.APIError{
					StatusCode: 404,
					Message:    "Not Found",
				}
			}
		},
	}

	service := NewService(apiClient, newMockCache())

	// Test searching for Pikachu by ID
	results, err := service.SearchPokemon("25")
	if err != nil {
		t.Fatalf("unexpected error searching for ID 25: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "pikachu" {
		t.Errorf("expected pikachu, got %s", results[0].Name)
	}

	// Test searching for Bulbasaur by ID
	results, err = service.SearchPokemon("1")
	if err != nil {
		t.Fatalf("unexpected error searching for ID 1: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "bulbasaur" {
		t.Errorf("expected bulbasaur, got %s", results[0].Name)
	}

	// Test searching for Mewtwo by ID
	results, err = service.SearchPokemon("150")
	if err != nil {
		t.Fatalf("unexpected error searching for ID 150: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].Name != "mewtwo" {
		t.Errorf("expected mewtwo, got %s", results[0].Name)
	}
}

// TestSearchPokemonByIDNotFound verifies error handling for invalid IDs
func TestSearchPokemonByIDNotFound(t *testing.T) {
	apiClient := &mockAPIClient{
		getPokemonFunc: func(nameOrID string) (*api.PokemonResponse, error) {
			return nil, &api.APIError{
				StatusCode: 404,
				Message:    "Not Found",
			}
		},
	}

	service := NewService(apiClient, newMockCache())
	_, err := service.SearchPokemon("99999")

	if err == nil {
		t.Fatal("expected error for invalid ID")
	}
}


// TestTransformPokemonResponse verifies transformation from API to domain model
func TestTransformPokemonResponse(t *testing.T) {
	apiResp := &api.PokemonResponse{
		ID:             25,
		Name:           "pikachu",
		Height:         4,
		Weight:         60,
		BaseExperience: 112,
		IsDefault:      true,
		Types: []api.PokemonType{
			{Type: api.NamedResource{Name: "electric"}},
		},
		Abilities: []api.PokemonAbility{
			{Ability: api.NamedResource{Name: "static"}, IsHidden: false, Slot: 1},
			{Ability: api.NamedResource{Name: "lightning-rod"}, IsHidden: true, Slot: 3},
		},
		Stats: []api.PokemonStat{
			{Stat: api.NamedResource{Name: "hp"}, BaseStat: 35, Effort: 0},
			{Stat: api.NamedResource{Name: "attack"}, BaseStat: 55, Effort: 0},
		},
		Sprites: api.Sprites{
			FrontDefault: "https://example.com/pikachu.png",
		},
	}

	pokemon := transformPokemonResponse(apiResp)

	if pokemon.ID != 25 {
		t.Errorf("expected ID 25, got %d", pokemon.ID)
	}

	if pokemon.Name != "pikachu" {
		t.Errorf("expected name pikachu, got %s", pokemon.Name)
	}

	if len(pokemon.Types) != 1 {
		t.Errorf("expected 1 type, got %d", len(pokemon.Types))
	}

	if pokemon.Types[0].Name != "electric" {
		t.Errorf("expected type electric, got %s", pokemon.Types[0].Name)
	}

	if len(pokemon.Abilities) != 2 {
		t.Errorf("expected 2 abilities, got %d", len(pokemon.Abilities))
	}

	if len(pokemon.Stats) != 2 {
		t.Errorf("expected 2 stats, got %d", len(pokemon.Stats))
	}

	if pokemon.Sprites == nil {
		t.Error("expected sprites to be set")
	}

	if pokemon.Sprites.FrontDefault != "https://example.com/pikachu.png" {
		t.Errorf("expected sprite URL, got %s", pokemon.Sprites.FrontDefault)
	}
}

// TestExtractPokemonIDFromURL verifies URL parsing for Pokemon IDs
func TestExtractPokemonIDFromURL(t *testing.T) {
	tests := []struct {
		url      string
		expected int
	}{
		{"https://pokeapi.co/api/v2/pokemon/1/", 1},
		{"https://pokeapi.co/api/v2/pokemon/25/", 25},
		{"https://pokeapi.co/api/v2/pokemon/150/", 150},
		{"https://pokeapi.co/api/v2/pokemon/1010/", 1010},
		{"https://pokeapi.co/api/v2/pokemon/1/", 1},      // With trailing slash
		{"https://pokeapi.co/api/v2/pokemon/1", 1},       // Without trailing slash
		{"invalid-url", 0},                                // Invalid URL
		{"", 0},                                           // Empty string
	}

	for _, tt := range tests {
		result := extractPokemonIDFromURL(tt.url)
		if result != tt.expected {
			t.Errorf("extractPokemonIDFromURL(%q) = %d, expected %d", tt.url, result, tt.expected)
		}
	}
}

// strPtr returns a pointer to a string
func strPtr(s string) *string {
	return &s
}
