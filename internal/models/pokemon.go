// Package models provides domain models for Pokemon data.
// These models are designed to be clean, simple, and focused on what the UI needs.
// They represent Pokemon in a way that's easy to work with and serialize to JSON.
package models

import (
	"fmt"
	"math"
)

// Pokemon represents a single Pokemon with all its details.
// This is the primary domain model used throughout the application.
//
// Example:
//
//	pokemon := &Pokemon{
//	    ID:     1,
//	    Name:   "Bulbasaur",
//	    Height: 7,    // decimeters (0.7m)
//	    Weight: 69,   // hectograms (6.9kg)
//	    Types:  []*Type{{Name: "Grass"}, {Name: "Poison"}},
//	    Stats:  []*Stat{{Name: "HP", BaseValue: 45}, ...},
//	}
type Pokemon struct {
	ID             int        `json:"id"`
	Name           string     `json:"name"`
	Height         int        `json:"height"` // in decimeters
	Weight         int        `json:"weight"` // in hectograms
	BaseExperience int        `json:"base_experience"`
	IsDefault      bool       `json:"is_default"`
	Types          []*Type    `json:"types"`
	Abilities      []*Ability `json:"abilities"`
	Stats          []*Stat    `json:"stats"`
	Sprites        *Sprite    `json:"sprites"`
	Forms          []string   `json:"forms,omitempty"`
	Generation     string     `json:"generation,omitempty"`
	Habitat        string     `json:"habitat,omitempty"`
	Color          string     `json:"color,omitempty"`
	Shape          string     `json:"shape,omitempty"`
	Description    string     `json:"description,omitempty"`
	EvolutionChain []int      `json:"evolution_chain,omitempty"`
}

// Type represents a Pokemon type (e.g., Fire, Water, Grass).
// Each type has a name and a color for UI rendering.
type Type struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// Ability represents a Pokemon ability (e.g., Blaze, Torrent).
type Ability struct {
	Name     string `json:"name"`
	IsHidden bool   `json:"is_hidden"`
	Slot     int    `json:"slot"`
}

// Stat represents a Pokemon stat (e.g., HP, Attack, Defense).
// BaseValue is the base stat value used for calculations.
type Stat struct {
	Name      string `json:"name"`
	BaseValue int    `json:"base_value"`
}

// Sprite contains URLs for various Pokemon image formats.
// This covers different image sources and formats for UI rendering.
type Sprite struct {
	FrontDefault     string `json:"front_default,omitempty"`
	FrontFemale      string `json:"front_female,omitempty"`
	FrontShiny       string `json:"front_shiny,omitempty"`
	FrontShinyFemale string `json:"front_shiny_female,omitempty"`
	BackDefault      string `json:"back_default,omitempty"`
	BackFemale       string `json:"back_female,omitempty"`
	BackShiny        string `json:"back_shiny,omitempty"`
	BackShinyFemale  string `json:"back_shiny_female,omitempty"`
	OfficialArtwork  string `json:"official_artwork,omitempty"`
}

// TypeColors maps Pokemon type names to their corresponding UI colors.
// These colors follow the official Pokemon game color scheme.
var TypeColors = map[string]string{
	"normal":   "#A8A878", // Gray
	"fire":     "#F08030", // Red/Orange
	"water":    "#6890F0", // Blue
	"grass":    "#78C850", // Green
	"electric": "#F8D030", // Yellow
	"ice":      "#98D8D8", // Cyan
	"fighting": "#C03028", // Dark Red
	"poison":   "#A040A0", // Purple
	"ground":   "#E0C068", // Brown/Tan
	"flying":   "#A890F0", // Light Purple
	"psychic":  "#F85888", // Pink
	"bug":      "#A8B820", // Olive
	"rock":     "#B8A038", // Brown
	"ghost":    "#705898", // Dark Purple
	"dragon":   "#7038F8", // Dark Blue/Purple
	"dark":     "#705848", // Dark Brown
	"steel":    "#B8B8D0", // Light Gray
	"fairy":    "#EE99AC", // Light Pink
}

// TypeNames returns a slice of type names for this Pokemon.
//
// Example:
//
//	types := pokemon.TypeNames() // []string{"Grass", "Poison"}
func (p *Pokemon) TypeNames() []string {
	if p == nil || p.Types == nil {
		return []string{}
	}

	names := make([]string, len(p.Types))
	for i, t := range p.Types {
		names[i] = t.Name
	}
	return names
}

// FormattedHeight returns the height in a human-readable format.
// Input is in decimeters, output is formatted as meters.
//
// Example:
//
//	pokemon.Height = 7
//	pokemon.FormattedHeight() // "0.7m"
//
//	pokemon.Height = 200
//	pokemon.FormattedHeight() // "20.0m"
func (p *Pokemon) FormattedHeight() string {
	if p == nil || p.Height == 0 {
		return "Unknown"
	}

	meters := float64(p.Height) / 10.0
	return fmt.Sprintf("%.1fm", meters)
}

// FormattedWeight returns the weight in a human-readable format.
// Input is in hectograms, output is formatted as kilograms.
//
// Example:
//
//	pokemon.Weight = 69
//	pokemon.FormattedWeight() // "6.9kg"
//
//	pokemon.Weight = 1000
//	pokemon.FormattedWeight() // "100.0kg"
func (p *Pokemon) FormattedWeight() string {
	if p == nil || p.Weight == 0 {
		return "Unknown"
	}

	kilograms := float64(p.Weight) / 10.0
	return fmt.Sprintf("%.1fkg", kilograms)
}

// MaxStat returns the highest base stat value for this Pokemon.
// This is useful for rendering stat bars in the UI.
//
// Example:
//
//	max := pokemon.MaxStat() // e.g., 115
//	// Use this to scale stat bar widths: width = (stat / max) * 100%
func (p *Pokemon) MaxStat() int {
	if p == nil || len(p.Stats) == 0 {
		return 0
	}

	max := 0
	for _, stat := range p.Stats {
		if stat != nil && stat.BaseValue > max {
			max = stat.BaseValue
		}
	}
	return max
}

// StatByName returns the stat with the given name, or nil if not found.
// Comparison is case-insensitive.
//
// Example:
//
//	hpStat := pokemon.StatByName("HP")
//	if hpStat != nil {
//	    fmt.Printf("HP: %d\n", hpStat.BaseValue)
//	}
func (p *Pokemon) StatByName(name string) *Stat {
	if p == nil || len(p.Stats) == 0 {
		return nil
	}

	for _, stat := range p.Stats {
		if stat != nil && stat.Name == name {
			return stat
		}
	}
	return nil
}

// TotalStats returns the sum of all base stat values.
// Higher total stats indicate stronger Pokemon overall.
//
// Example:
//
//	total := pokemon.TotalStats() // e.g., 318
func (p *Pokemon) TotalStats() int {
	if p == nil || len(p.Stats) == 0 {
		return 0
	}

	total := 0
	for _, stat := range p.Stats {
		if stat != nil {
			total += stat.BaseValue
		}
	}
	return total
}

// PrimaryType returns the first type for this Pokemon.
// Every Pokemon has at least one type.
//
// Example:
//
//	primary := pokemon.PrimaryType() // "Grass"
func (p *Pokemon) PrimaryType() string {
	if p == nil || len(p.Types) == 0 {
		return "Unknown"
	}
	if p.Types[0] == nil {
		return "Unknown"
	}
	return p.Types[0].Name
}

// SecondaryType returns the second type for this Pokemon, if it exists.
// Some Pokemon have only one type.
//
// Example:
//
//	secondary := pokemon.SecondaryType() // "Poison" or ""
func (p *Pokemon) SecondaryType() string {
	if p == nil || len(p.Types) < 2 {
		return ""
	}
	if p.Types[1] == nil {
		return ""
	}
	return p.Types[1].Name
}

// IsDualType returns true if the Pokemon has two types.
func (p *Pokemon) IsDualType() bool {
	return p != nil && len(p.Types) >= 2 && p.Types[1] != nil
}

// BestSprite returns the best available sprite URL for this Pokemon.
// Preference order: Official Artwork > Front Shiny > Front Default
//
// Example:
//
//	imageURL := pokemon.BestSprite()
func (p *Pokemon) BestSprite() string {
	if p == nil || p.Sprites == nil {
		return ""
	}

	if p.Sprites.OfficialArtwork != "" {
		return p.Sprites.OfficialArtwork
	}
	if p.Sprites.FrontShiny != "" {
		return p.Sprites.FrontShiny
	}
	if p.Sprites.FrontDefault != "" {
		return p.Sprites.FrontDefault
	}
	return ""
}

// GetTypeColor returns the UI color for the Pokemon's primary type.
// Returns a default gray color if the type is not found.
//
// Example:
//
//	color := pokemon.GetTypeColor() // "#78C850"
func (p *Pokemon) GetTypeColor() string {
	if p == nil || len(p.Types) == 0 || p.Types[0] == nil {
		return TypeColors["normal"]
	}

	typeName := p.Types[0].Name
	if color, exists := TypeColors[typeName]; exists {
		return color
	}
	return TypeColors["normal"]
}

// TypeColorMap returns a map of type name to color for all types in this Pokemon.
//
// Example:
//
//	colorMap := pokemon.TypeColorMap()
//	// {"Grass": "#78C850", "Poison": "#A040A0"}
func (p *Pokemon) TypeColorMap() map[string]string {
	colorMap := make(map[string]string)
	if p == nil || len(p.Types) == 0 {
		return colorMap
	}

	for _, t := range p.Types {
		if t != nil {
			typeName := t.Name
			if color, exists := TypeColors[typeName]; exists {
				colorMap[typeName] = color
			} else {
				colorMap[typeName] = TypeColors["normal"]
			}
		}
	}
	return colorMap
}

// AbilityNames returns a slice of ability names for this Pokemon.
//
// Example:
//
//	abilities := pokemon.AbilityNames() // []string{"Overgrow"}
func (p *Pokemon) AbilityNames() []string {
	if p == nil || len(p.Abilities) == 0 {
		return []string{}
	}

	names := make([]string, len(p.Abilities))
	for i, a := range p.Abilities {
		if a != nil {
			names[i] = a.Name
		}
	}
	return names
}

// HiddenAbility returns the hidden ability for this Pokemon, if it exists.
//
// Example:
//
//	hidden := pokemon.HiddenAbility() // "Chlorophyll" or ""
func (p *Pokemon) HiddenAbility() string {
	if p == nil || len(p.Abilities) == 0 {
		return ""
	}

	for _, ability := range p.Abilities {
		if ability != nil && ability.IsHidden {
			return ability.Name
		}
	}
	return ""
}

// NormalAbilities returns all non-hidden abilities for this Pokemon.
//
// Example:
//
//	normal := pokemon.NormalAbilities() // []string{"Overgrow"}
func (p *Pokemon) NormalAbilities() []string {
	if p == nil || len(p.Abilities) == 0 {
		return []string{}
	}

	var abilities []string
	for _, ability := range p.Abilities {
		if ability != nil && !ability.IsHidden {
			abilities = append(abilities, ability.Name)
		}
	}
	return abilities
}

// HasAbility returns true if the Pokemon has an ability with the given name.
// Comparison is case-insensitive.
//
// Example:
//
//	has := pokemon.HasAbility("Overgrow") // true
func (p *Pokemon) HasAbility(name string) bool {
	if p == nil || len(p.Abilities) == 0 {
		return false
	}

	for _, ability := range p.Abilities {
		if ability != nil && ability.Name == name {
			return true
		}
	}
	return false
}

// StatPercentage returns the base stat as a percentage of the maximum base stat value.
// This is useful for rendering stat bars that scale to 100%.
// The maximum base stat is 255 (Blissey's HP).
//
// Example:
//
//	percent := pokemon.StatPercentage(pokemon.StatByName("HP"))
//	// 45/255 * 100 ≈ 17.65%
func (p *Pokemon) StatPercentage(stat *Stat) float64 {
	if stat == nil {
		return 0
	}
	// Maximum base stat across all Pokemon stats is 255
	return (float64(stat.BaseValue) / 255.0) * 100.0
}

// AverageStatValue returns the average of all base stats.
// This gives a sense of the Pokemon's overall stat distribution.
//
// Example:
//
//	avg := pokemon.AverageStatValue() // e.g., 53
func (p *Pokemon) AverageStatValue() float64 {
	if p == nil || len(p.Stats) == 0 {
		return 0
	}

	return float64(p.TotalStats()) / float64(len(p.Stats))
}

// IsHighStat returns true if the Pokemon has a stat value that exceeds the given threshold.
// Useful for identifying Pokemon strengths.
//
// Example:
//
//	strong := pokemon.IsHighStat("Attack", 100) // true if Attack >= 100
func (p *Pokemon) IsHighStat(statName string, threshold int) bool {
	stat := p.StatByName(statName)
	return stat != nil && stat.BaseValue >= threshold
}

// StrongestStat returns the name and value of the Pokemon's highest stat.
//
// Example:
//
//	name, value := pokemon.StrongestStat() // "HP", 45
func (p *Pokemon) StrongestStat() (string, int) {
	if p == nil || len(p.Stats) == 0 {
		return "", 0
	}

	strongest := p.Stats[0]
	for _, stat := range p.Stats {
		if stat != nil && stat.BaseValue > strongest.BaseValue {
			strongest = stat
		}
	}

	if strongest == nil {
		return "", 0
	}
	return strongest.Name, strongest.BaseValue
}

// WeakestStat returns the name and value of the Pokemon's lowest stat.
//
// Example:
//
//	name, value := pokemon.WeakestStat() // "Defense", 40
func (p *Pokemon) WeakestStat() (string, int) {
	if p == nil || len(p.Stats) == 0 {
		return "", 0
	}

	weakest := p.Stats[0]
	for _, stat := range p.Stats {
		if stat != nil && stat.BaseValue < weakest.BaseValue {
			weakest = stat
		}
	}

	if weakest == nil {
		return "", 0
	}
	return weakest.Name, weakest.BaseValue
}

// StatVariance returns the variance of all stat values.
// Lower variance means stats are more balanced.
// Higher variance means the Pokemon is specialized in certain stats.
//
// Example:
//
//	variance := pokemon.StatVariance() // e.g., 456.7
func (p *Pokemon) StatVariance() float64 {
	if p == nil || len(p.Stats) <= 1 {
		return 0
	}

	avg := p.AverageStatValue()
	sumSquaredDiff := 0.0

	for _, stat := range p.Stats {
		if stat != nil {
			diff := float64(stat.BaseValue) - avg
			sumSquaredDiff += diff * diff
		}
	}

	return sumSquaredDiff / float64(len(p.Stats))
}

// StatStandardDeviation returns the standard deviation of all stat values.
// This is the square root of the variance.
//
// Example:
//
//	stdDev := pokemon.StatStandardDeviation() // e.g., 21.37
func (p *Pokemon) StatStandardDeviation() float64 {
	return math.Sqrt(p.StatVariance())
}

// IsBalanced returns true if the Pokemon's stats are relatively balanced.
// A Pokemon is considered balanced if its stat standard deviation is below 20.
//
// Example:
//
//	balanced := pokemon.IsBalanced() // true for well-rounded Pokemon
func (p *Pokemon) IsBalanced() bool {
	return p.StatStandardDeviation() < 20
}

// IsSpecialized returns true if the Pokemon is specialized in certain stats.
// A Pokemon is considered specialized if its stat standard deviation is above 30.
//
// Example:
//
//	specialized := pokemon.IsSpecialized() // true for Alakazam (high Sp. Atk)
func (p *Pokemon) IsSpecialized() bool {
	return p.StatStandardDeviation() > 30
}
