package models

import (
	"strings"
	"testing"
)

func TestHasSprite(t *testing.T) {
	t.Run("nil pokemon", func(t *testing.T) {
		var p *Pokemon
		if p.HasSprite() {
			t.Error("nil Pokemon should not have sprite")
		}
	})

	t.Run("nil sprites", func(t *testing.T) {
		p := &Pokemon{Name: "Test"}
		if p.HasSprite() {
			t.Error("Pokemon with nil Sprites should not have sprite")
		}
	})

	t.Run("empty sprite URL", func(t *testing.T) {
		p := &Pokemon{
			Name:    "Test",
			Sprites: &Sprite{},
		}
		if p.HasSprite() {
			t.Error("Pokemon with empty FrontDefault should not have sprite")
		}
	})

	t.Run("has sprite", func(t *testing.T) {
		p := &Pokemon{
			Name: "Test",
			Sprites: &Sprite{
				FrontDefault: "https://example.com/sprite.png",
			},
		}
		if !p.HasSprite() {
			t.Error("Pokemon with FrontDefault should have sprite")
		}
	})
}

func TestGetSpriteURL(t *testing.T) {
	t.Run("nil pokemon", func(t *testing.T) {
		var p *Pokemon
		url := p.GetSpriteURL()
		if url != "" {
			t.Errorf("expected empty string, got %q", url)
		}
	})

	t.Run("nil sprites", func(t *testing.T) {
		p := &Pokemon{Name: "Test"}
		url := p.GetSpriteURL()
		if url != "" {
			t.Errorf("expected empty string, got %q", url)
		}
	})

	t.Run("returns front default", func(t *testing.T) {
		expected := "https://example.com/sprite.png"
		p := &Pokemon{
			Name: "Test",
			Sprites: &Sprite{
				FrontDefault: expected,
			},
		}
		url := p.GetSpriteURL()
		if url != expected {
			t.Errorf("expected %q, got %q", expected, url)
		}
	})
}

func TestGetSpriteASCII(t *testing.T) {
	t.Run("nil pokemon", func(t *testing.T) {
		var p *Pokemon
		ascii := p.GetSpriteASCII()
		if ascii == "" {
			t.Error("expected non-empty ASCII art for nil Pokemon")
		}
		// Should have 5 lines
		lines := strings.Split(ascii, "\n")
		if len(lines) != 5 {
			t.Errorf("expected 5 lines, got %d", len(lines))
		}
	})

	t.Run("electric type pikachu", func(t *testing.T) {
		p := &Pokemon{
			Name:  "pikachu",
			Types: []*Type{{Name: "electric"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "⚡") {
			t.Error("Electric type should contain lightning bolt (⚡)")
		}
		if !strings.Contains(ascii, "p") {
			t.Error("ASCII should contain first letter of name")
		}
		// Should have 5 lines
		lines := strings.Split(ascii, "\n")
		if len(lines) != 5 {
			t.Errorf("expected 5 lines, got %d", len(lines))
		}
	})

	t.Run("fire type charizard", func(t *testing.T) {
		p := &Pokemon{
			Name:  "charizard",
			Types: []*Type{{Name: "fire"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "▓") {
			t.Error("Fire type should contain block character (▓)")
		}
		if !strings.Contains(ascii, "c") {
			t.Error("ASCII should contain first letter of name")
		}
	})

	t.Run("water type squirtle", func(t *testing.T) {
		p := &Pokemon{
			Name:  "squirtle",
			Types: []*Type{{Name: "water"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "≈") {
			t.Error("Water type should contain wave character (≈)")
		}
		if !strings.Contains(ascii, "s") {
			t.Error("ASCII should contain first letter of name")
		}
	})

	t.Run("grass type bulbasaur", func(t *testing.T) {
		p := &Pokemon{
			Name:  "bulbasaur",
			Types: []*Type{{Name: "grass"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "❀") {
			t.Error("Grass type should contain flower character (❀)")
		}
		if !strings.Contains(ascii, "b") {
			t.Error("ASCII should contain first letter of name")
		}
	})

	t.Run("ice type", func(t *testing.T) {
		p := &Pokemon{
			Name:  "articuno",
			Types: []*Type{{Name: "ice"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "❄") {
			t.Error("Ice type should contain snowflake character (❄)")
		}
	})

	t.Run("psychic type", func(t *testing.T) {
		p := &Pokemon{
			Name:  "mewtwo",
			Types: []*Type{{Name: "psychic"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "◉") {
			t.Error("Psychic type should contain circle character (◉)")
		}
	})

	t.Run("dragon type", func(t *testing.T) {
		p := &Pokemon{
			Name:  "dragonite",
			Types: []*Type{{Name: "dragon"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "◆") {
			t.Error("Dragon type should contain diamond character (◆)")
		}
	})

	t.Run("ghost type", func(t *testing.T) {
		p := &Pokemon{
			Name:  "gengar",
			Types: []*Type{{Name: "ghost"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "▒") {
			t.Error("Ghost type should contain block character (▒)")
		}
	})

	t.Run("default type normal", func(t *testing.T) {
		p := &Pokemon{
			Name:  "raticate",
			Types: []*Type{{Name: "normal"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "█") {
			t.Error("Normal type should contain block character (█)")
		}
		if !strings.Contains(ascii, "r") {
			t.Error("ASCII should contain first letter of name")
		}
	})

	t.Run("empty name", func(t *testing.T) {
		p := &Pokemon{
			Name:  "",
			Types: []*Type{{Name: "normal"}},
		}
		ascii := p.GetSpriteASCII()
		if !strings.Contains(ascii, "?") {
			t.Error("Empty name should use '?' as placeholder")
		}
	})

	t.Run("no types", func(t *testing.T) {
		p := &Pokemon{
			Name:  "missingno",
			Types: []*Type{},
		}
		ascii := p.GetSpriteASCII()
		// Should not panic and should return default sprite
		if ascii == "" {
			t.Error("expected non-empty ASCII art")
		}
	})
}
