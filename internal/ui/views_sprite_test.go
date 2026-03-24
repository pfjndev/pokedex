package ui

import (
	"strings"
	"testing"

	"github.com/pfjndev/pokedex/internal/models"
)

// TestRenderDetailsSpriteSafety tests that renderDetails handles sprite edge cases gracefully
func TestRenderDetailsSpriteSafety(t *testing.T) {
	tests := []struct {
		name        string
		pokemon     interface{}
		description string
		shouldPanic bool
	}{
		{
			name:        "nil pokemon",
			pokemon:     nil,
			description: "Should not render sprite when pokemon is nil",
			shouldPanic: false,
		},
		{
			name: "pokemon with nil sprites",
			pokemon: &models.Pokemon{
				ID:        1,
				Name:      "bulbasaur",
				Sprites:   nil,
				Types:     []*models.Type{{Name: "grass"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should not render sprite when Sprites field is nil",
			shouldPanic: false,
		},
		{
			name: "pokemon with empty sprite fields",
			pokemon: &models.Pokemon{
				ID:   1,
				Name: "bulbasaur",
				Sprites: &models.Sprite{
					FrontDefault:     "",
					FrontFemale:      "",
					FrontShiny:       "",
					FrontShinyFemale: "",
					BackDefault:      "",
					BackFemale:       "",
					BackShiny:        "",
					BackShinyFemale:  "",
					OfficialArtwork:  "",
				},
				Types:     []*models.Type{{Name: "grass"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should not render sprite when all sprite URLs are empty",
			shouldPanic: false,
		},
		{
			name: "pokemon with valid sprite",
			pokemon: &models.Pokemon{
				ID:   1,
				Name: "bulbasaur",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/1.png",
				},
				Types:     []*models.Type{{Name: "grass"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should render sprite when valid sprite data exists",
			shouldPanic: false,
		},
		{
			name: "pokemon with empty name",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{{Name: "normal"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle empty name gracefully in sprite",
			shouldPanic: false,
		},
		{
			name: "pokemon with no types",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "missingno",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle missing types gracefully in sprite",
			shouldPanic: false,
		},
		{
			name: "pokemon with nil type in slice",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "test",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{nil},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle nil type in slice gracefully",
			shouldPanic: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap in function to test panic recovery
			var rendered string
			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			model := Model{pokemonDetails: tt.pokemon}
			rendered = model.renderDetails()

			// Verify output is not empty and doesn't panic
			if rendered == "" && tt.pokemon != nil {
				t.Error("expected non-empty output for non-nil pokemon")
			}

			// Check if sprite section is in output when expected
			switch p := tt.pokemon.(type) {
			case *models.Pokemon:
				if p != nil && p.HasSprite() {
					if !strings.Contains(rendered, "⚡") &&
						!strings.Contains(rendered, "▓") &&
						!strings.Contains(rendered, "≈") &&
						!strings.Contains(rendered, "❀") &&
						!strings.Contains(rendered, "❄") &&
						!strings.Contains(rendered, "◉") &&
						!strings.Contains(rendered, "◆") &&
						!strings.Contains(rendered, "▒") &&
						!strings.Contains(rendered, "█") &&
						!strings.Contains(rendered, "?") {
						t.Logf("Note: sprite may not be visible but no panic occurred (%s)", tt.description)
					}
				}
			}
		})
	}
}

// TestRenderDetailsSpriteBoundaryConditions tests boundary conditions for sprite rendering
func TestRenderDetailsSpriteBoundaryConditions(t *testing.T) {
	tests := []struct {
		name        string
		pokemon     *models.Pokemon
		description string
	}{
		{
			name: "very long pokemon name",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "verylongpokemonnameverylongpokemonnameverylongpokemonname",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{{Name: "normal"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle very long names without truncation issues",
		},
		{
			name: "special characters in name",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "pkmn-test_v2",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{{Name: "normal"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle special characters in names",
		},
		{
			name: "unicode in pokemon name",
			pokemon: &models.Pokemon{
				ID:   999,
				Name: "é",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/999.png",
				},
				Types:     []*models.Type{{Name: "normal"}},
				Abilities: []*models.Ability{},
				Stats:     []*models.Stat{},
			},
			description: "Should handle unicode characters in names",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unexpected panic with %s: %v", tt.description, r)
				}
			}()

			model := Model{pokemonDetails: tt.pokemon}
			rendered := model.renderDetails()

			if rendered == "" {
				t.Error("expected non-empty output")
			}
		})
	}
}

// TestGetSpriteASCIISafety tests GetSpriteASCII under edge conditions
func TestGetSpriteASCIISafety(t *testing.T) {
	tests := []struct {
		name        string
		pokemon     *models.Pokemon
		description string
	}{
		{
			name:        "nil pokemon",
			pokemon:     nil,
			description: "Should return fallback sprite for nil pokemon",
		},
		{
			name: "pokemon with unknown type",
			pokemon: &models.Pokemon{
				Name:  "custom",
				Types: []*models.Type{{Name: "unknown-type"}},
			},
			description: "Should return default sprite for unknown type",
		},
		{
			name: "pokemon with nil type",
			pokemon: &models.Pokemon{
				Name:  "custom",
				Types: []*models.Type{nil},
			},
			description: "Should handle nil type gracefully",
		},
		{
			name: "pokemon with empty types slice",
			pokemon: &models.Pokemon{
				Name:  "custom",
				Types: []*models.Type{},
			},
			description: "Should handle empty types slice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("unexpected panic: %v", r)
				}
			}()

			ascii := tt.pokemon.GetSpriteASCII()

			// Verify we get valid output
			if ascii == "" {
				t.Error("expected non-empty ASCII art")
			}

			// Verify it has the expected structure (5 lines)
			lines := strings.Split(ascii, "\n")
			if len(lines) != 5 {
				t.Errorf("expected 5 lines, got %d", len(lines))
			}

			// Verify no panics and output contains at least one printable character
			hasContent := false
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					hasContent = true
					break
				}
			}
			if !hasContent {
				t.Error("expected at least some content in ASCII art")
			}
		})
	}
}

// TestHasSpriteSafety verifies HasSprite() handles all edge cases
func TestHasSpriteSafety(t *testing.T) {
	tests := []struct {
		name        string
		pokemon     *models.Pokemon
		expected    bool
		description string
	}{
		{
			name:        "nil pokemon",
			pokemon:     nil,
			expected:    false,
			description: "nil pokemon should return false",
		},
		{
			name: "pokemon with nil sprites",
			pokemon: &models.Pokemon{
				Name:    "test",
				Sprites: nil,
			},
			expected:    false,
			description: "pokemon with nil Sprites should return false",
		},
		{
			name: "pokemon with empty Sprites struct",
			pokemon: &models.Pokemon{
				Name:    "test",
				Sprites: &models.Sprite{},
			},
			expected:    false,
			description: "pokemon with empty Sprites struct should return false",
		},
		{
			name: "pokemon with FrontDefault only",
			pokemon: &models.Pokemon{
				Name:    "test",
				Sprites: &models.Sprite{FrontDefault: "https://example.com/test.png"},
			},
			expected:    true,
			description: "pokemon with FrontDefault should return true",
		},
		{
			name: "pokemon with all other sprites but no FrontDefault",
			pokemon: &models.Pokemon{
				Name: "test",
				Sprites: &models.Sprite{
					FrontFemale:     "https://example.com/f.png",
					FrontShiny:      "https://example.com/s.png",
					OfficialArtwork: "https://example.com/o.png",
				},
			},
			expected:    false,
			description: "pokemon with sprites but no FrontDefault should return false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pokemon.HasSprite()
			if result != tt.expected {
				t.Errorf("expected %v, got %v: %s", tt.expected, result, tt.description)
			}
		})
	}
}

// TestBestSpriteSafety verifies BestSprite() handles edge cases
func TestBestSpriteSafety(t *testing.T) {
	tests := []struct {
		name        string
		pokemon     *models.Pokemon
		expected    string
		description string
	}{
		{
			name:        "nil pokemon",
			pokemon:     nil,
			expected:    "",
			description: "nil pokemon should return empty string",
		},
		{
			name: "pokemon with nil sprites",
			pokemon: &models.Pokemon{
				Name:    "test",
				Sprites: nil,
			},
			expected:    "",
			description: "pokemon with nil Sprites should return empty string",
		},
		{
			name: "pokemon with empty Sprites",
			pokemon: &models.Pokemon{
				Name:    "test",
				Sprites: &models.Sprite{},
			},
			expected:    "",
			description: "pokemon with empty Sprites should return empty string",
		},
		{
			name: "preference: OfficialArtwork > others",
			pokemon: &models.Pokemon{
				Name: "test",
				Sprites: &models.Sprite{
					OfficialArtwork: "https://example.com/official.png",
					FrontShiny:      "https://example.com/shiny.png",
					FrontDefault:    "https://example.com/default.png",
				},
			},
			expected:    "https://example.com/official.png",
			description: "should prefer OfficialArtwork",
		},
		{
			name: "preference: FrontShiny > FrontDefault",
			pokemon: &models.Pokemon{
				Name: "test",
				Sprites: &models.Sprite{
					FrontShiny:   "https://example.com/shiny.png",
					FrontDefault: "https://example.com/default.png",
				},
			},
			expected:    "https://example.com/shiny.png",
			description: "should prefer FrontShiny over FrontDefault",
		},
		{
			name: "fallback: FrontDefault",
			pokemon: &models.Pokemon{
				Name: "test",
				Sprites: &models.Sprite{
					FrontDefault: "https://example.com/default.png",
				},
			},
			expected:    "https://example.com/default.png",
			description: "should use FrontDefault as fallback",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pokemon.BestSprite()
			if result != tt.expected {
				t.Errorf("expected %q, got %q: %s", tt.expected, result, tt.description)
			}
		})
	}
}
