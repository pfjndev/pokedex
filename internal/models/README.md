# Pokemon Domain Models

This package provides clean, UI-focused domain models for Pokemon data. The models are designed to be simple, well-documented, and easy to serialize to JSON for API responses.

## Core Models

### Pokemon
The main domain model representing a single Pokemon with all its details.

```go
type Pokemon struct {
    ID                   int       // Unique identifier
    Name                 string    // Pokemon name
    Height               int       // in decimeters (divide by 10 for meters)
    Weight               int       // in hectograms (divide by 10 for kilograms)
    BaseExperience       int       // Base experience gained
    IsDefault            bool      // Is this the default form?
    Types                []*Type   // Type(s) of the Pokemon
    Abilities            []*Ability // Abilities this Pokemon can have
    Stats                []*Stat   // Base stats (HP, Attack, Defense, etc.)
    Sprites              *Sprite   // Image URLs for different formats
    Forms                []string  // Alternative forms
    Generation           string    // Which generation introduced this Pokemon
    Habitat              string    // Natural habitat
    Color                string    // Official Pokemon color
    Shape                string    // Body shape category
    Description          string    // Pokedex description
    EvolutionChain       []int     // IDs of evolution chain members
}
```

### Type
Represents a Pokemon type (Fire, Water, Grass, etc.) with UI color information.

```go
type Type struct {
    Name  string // Type name (e.g., "Fire", "Water")
    Color string // Hex color code for UI rendering (e.g., "#F08030")
}
```

### Ability
Represents a Pokemon ability (skill).

```go
type Ability struct {
    Name     string // Ability name (e.g., "Blaze", "Overgrow")
    IsHidden bool   // Is this a hidden ability?
    Slot     int    // Which ability slot (1, 2, or 3)
}
```

### Stat
Represents a single base stat for a Pokemon.

```go
type Stat struct {
    Name      string // Stat name (e.g., "HP", "Attack", "Defense")
    BaseValue int    // Base stat value (0-255)
}
```

### Sprite
Contains URLs for various Pokemon image formats.

```go
type Sprite struct {
    FrontDefault      string // Default front-facing image
    FrontFemale       string // Front female variant (if exists)
    FrontShiny        string // Shiny front variant
    FrontShinyFemale  string // Shiny female variant (if exists)
    BackDefault       string // Default back-facing image
    BackFemale        string // Back female variant (if exists)
    BackShiny         string // Shiny back variant
    BackShinyFemale   string // Shiny female back variant (if exists)
    OfficialArtwork   string // Official Pokemon artwork
}
```

## Type Colors

All 18 Pokemon types have predefined colors for UI rendering:

```go
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
```

## Helper Methods

### Type Information
- **`TypeNames() []string`** - Returns slice of type names
- **`PrimaryType() string`** - Returns the first type
- **`SecondaryType() string`** - Returns the second type (if exists)
- **`IsDualType() bool`** - Returns true if Pokemon has two types
- **`GetTypeColor() string`** - Returns hex color for primary type
- **`TypeColorMap() map[string]string`** - Returns map of all types to colors

### Physical Attributes
- **`FormattedHeight() string`** - Returns human-readable height (e.g., "0.7m")
- **`FormattedWeight() string`** - Returns human-readable weight (e.g., "6.9kg")

### Sprites
- **`BestSprite() string`** - Returns best available sprite URL (prefers official artwork)

### Abilities
- **`AbilityNames() []string`** - Returns all ability names
- **`HiddenAbility() string`** - Returns hidden ability name (if exists)
- **`NormalAbilities() []string`** - Returns non-hidden ability names
- **`HasAbility(name string) bool`** - Checks if Pokemon has specific ability

### Stats
- **`MaxStat() int`** - Returns highest base stat value
- **`TotalStats() int`** - Returns sum of all base stats
- **`StatByName(name string) *Stat`** - Returns specific stat by name
- **`StatPercentage(stat *Stat) float64`** - Returns stat as percentage of max (255)
- **`AverageStatValue() float64`** - Returns average of all stats
- **`IsHighStat(name string, threshold int) bool`** - Checks if stat exceeds threshold
- **`StrongestStat() (string, int)`** - Returns name and value of highest stat
- **`WeakestStat() (string, int)`** - Returns name and value of lowest stat

### Stat Analysis
- **`StatVariance() float64`** - Returns variance of stat values
- **`StatStandardDeviation() float64`** - Returns standard deviation of stats
- **`IsBalanced() bool`** - Returns true if stats are well-balanced (stddev < 20)
- **`IsSpecialized() bool`** - Returns true if Pokemon specializes in certain stats (stddev > 30)

## Examples

### Creating a Pokemon
```go
pokemon := &Pokemon{
    ID:     1,
    Name:   "Bulbasaur",
    Height: 7,    // 0.7m
    Weight: 69,   // 6.9kg
    Types: []*Type{
        {Name: "Grass", Color: "#78C850"},
        {Name: "Poison", Color: "#A040A0"},
    },
    Stats: []*Stat{
        {Name: "HP", BaseValue: 45},
        {Name: "Attack", BaseValue: 49},
        {Name: "Defense", BaseValue: 49},
        {Name: "Sp. Atk", BaseValue: 65},
        {Name: "Sp. Def", BaseValue: 65},
        {Name: "Speed", BaseValue: 45},
    },
    Sprites: &Sprite{
        FrontDefault: "https://raw.githubusercontent.com/PokeAPI/sprites/master/pokemon/1.png",
        OfficialArtwork: "https://raw.githubusercontent.com/PokeAPI/sprites/master/pokemon/other/official-artwork/1.png",
    },
}
```

### Displaying Pokemon Info
```go
fmt.Printf("%s (%d)\n", pokemon.Name, pokemon.ID)
fmt.Printf("Types: %v\n", pokemon.TypeNames())
fmt.Printf("Height: %s\n", pokemon.FormattedHeight())
fmt.Printf("Weight: %s\n", pokemon.FormattedWeight())
fmt.Printf("Total Stats: %d\n", pokemon.TotalStats())
fmt.Printf("Average Stat: %.1f\n", pokemon.AverageStatValue())

strongest, val := pokemon.StrongestStat()
fmt.Printf("Strongest: %s (%d)\n", strongest, val)
```

### Rendering Stats UI
```go
// For stat bars
maxStat := pokemon.MaxStat()
for _, stat := range pokemon.Stats {
    percentage := (float64(stat.BaseValue) / float64(maxStat)) * 100
    fmt.Printf("%-8s: %3d [%-50s] %.1f%%\n", 
        stat.Name, stat.BaseValue, 
        strings.Repeat("█", int(percentage/2)), 
        percentage)
}
```

### Type Color UI
```go
// Get colors for styling
colorMap := pokemon.TypeColorMap()
for typeName, color := range colorMap {
    fmt.Printf("<span style='background-color: %s'>%s</span>\n", color, typeName)
}
```

### Analyzing Pokemon Stats
```go
if pokemon.IsBalanced() {
    fmt.Println("This is a well-rounded Pokemon")
} else if pokemon.IsSpecialized() {
    fmt.Println("This Pokemon specializes in certain stats")
}

// Find Pokemon weak points
weakest, weakVal := pokemon.WeakestStat()
fmt.Printf("Watch out for low %s: %d\n", weakest, weakVal)
```

## Usage in API Responses

These models serialize cleanly to JSON:

```go
// Example JSON output
{
    "id": 1,
    "name": "Bulbasaur",
    "height": 7,
    "weight": 69,
    "base_experience": 64,
    "is_default": true,
    "types": [
        {"name": "Grass", "color": "#78C850"},
        {"name": "Poison", "color": "#A040A0"}
    ],
    "abilities": [
        {"name": "Overgrow", "is_hidden": false, "slot": 1},
        {"name": "Chlorophyll", "is_hidden": true, "slot": 3}
    ],
    "stats": [
        {"name": "HP", "base_value": 45},
        {"name": "Attack", "base_value": 49},
        {"name": "Defense", "base_value": 49},
        {"name": "Sp. Atk", "base_value": 65},
        {"name": "Sp. Def", "base_value": 65},
        {"name": "Speed", "base_value": 45}
    ],
    "sprites": {
        "front_default": "https://...",
        "official_artwork": "https://..."
    }
}
```

## Design Principles

1. **UI-Focused** - Models include information needed for rendering (colors, formatted values)
2. **Simple Types** - Uses standard Go types (int, string, bool)
3. **Nil-Safe** - Methods handle nil pointers gracefully
4. **Well-Documented** - Every type and method has clear examples
5. **Serializable** - All types have JSON tags for API responses
6. **Helpful Methods** - Common operations are built-in (calculations, lookups, formatting)

## Notes

- Heights are in decimeters (divide by 10 for meters)
- Weights are in hectograms (divide by 10 for kilograms)
- Base stats range from 0-255, with 255 being the maximum
- Type colors follow the official Pokemon game color scheme
- All slices are zero-length (not nil) when empty to allow iteration
