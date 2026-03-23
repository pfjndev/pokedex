package api

// PokemonResponse represents the full Pokemon data from the PokeAPI v2 schema.
// This is the complete response when fetching a single Pokemon by name or ID.
type PokemonResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	BaseExperience int              `json:"base_experience"`
	IsDefault      bool             `json:"is_default"`
	Abilities      []PokemonAbility `json:"abilities"`
	Stats          []PokemonStat    `json:"stats"`
	Types          []PokemonType    `json:"types"`
	Sprites        Sprites          `json:"sprites"`
	Forms          []Form           `json:"forms"`
	HeldItems      []HeldItem       `json:"held_items"`
	Moves          []Move           `json:"moves"`
	GameIndices    []GameIndex      `json:"game_indices"`
	Cries          Cries            `json:"cries"`
}

// PokemonAbility represents an ability that a Pokemon can have.
type PokemonAbility struct {
	Ability  NamedResource `json:"ability"`
	IsHidden bool          `json:"is_hidden"`
	Slot     int           `json:"slot"`
}

// PokemonStat represents a stat of a Pokemon (HP, Attack, Defense, etc).
type PokemonStat struct {
	Stat     NamedResource `json:"stat"`
	Effort   int           `json:"effort"`
	BaseStat int           `json:"base_stat"`
}

// PokemonType represents a type that a Pokemon can have (Fire, Water, Grass, etc).
type PokemonType struct {
	Slot int           `json:"slot"`
	Type NamedResource `json:"type"`
}

// Sprites contains the image URLs for a Pokemon in different formats.
type Sprites struct {
	FrontDefault     string          `json:"front_default"`
	FrontShiny       string          `json:"front_shiny"`
	FrontFemale      string          `json:"front_female"`
	FrontShinyFemale string          `json:"front_shiny_female"`
	BackDefault      string          `json:"back_default"`
	BackShiny        string          `json:"back_shiny"`
	BackFemale       string          `json:"back_female"`
	BackShinyFemale  string          `json:"back_shiny_female"`
	OfficialArtwork  OfficialArtwork `json:"other"`
}

// OfficialArtwork contains official artwork URLs for a Pokemon.
type OfficialArtwork struct {
	Artwork ArtworkData `json:"official-artwork"`
}

// ArtworkData contains the official artwork URLs.
type ArtworkData struct {
	FrontDefault string `json:"front_default"`
	FrontShiny   string `json:"front_shiny"`
}

// Form represents a Pokemon form (e.g., Alola form, Galar form).
type Form struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// HeldItem represents an item that a Pokemon can hold.
type HeldItem struct {
	Item           NamedResource   `json:"item"`
	VersionDetails []VersionDetail `json:"version_details"`
}

// VersionDetail contains version-specific details for a held item.
type VersionDetail struct {
	Rarity  int           `json:"rarity"`
	Version NamedResource `json:"version"`
}

// Move represents a move that a Pokemon can learn.
type Move struct {
	Move                NamedResource        `json:"move"`
	VersionGroupDetails []VersionGroupDetail `json:"version_group_details"`
}

// VersionGroupDetail contains version group specific details for a move.
type VersionGroupDetail struct {
	LevelLearnedAt  int           `json:"level_learned_at"`
	MoveLearnMethod NamedResource `json:"move_learn_method"`
	Order           *int          `json:"order"`
	VersionGroup    NamedResource `json:"version_group"`
}

// GameIndex represents a Pokemon's index in a specific game.
type GameIndex struct {
	GameIndex int           `json:"game_index"`
	Version   NamedResource `json:"version"`
}

// Cries contains the URLs to cry audio files for a Pokemon.
type Cries struct {
	Latest string `json:"latest"`
	Legacy string `json:"legacy"`
}

// NamedResource is a common structure used throughout the PokeAPI
// that contains the name and URL of a resource.
type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PokemonListResponse represents a paginated list of Pokemon from the PokeAPI.
type PokemonListResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []PokemonListItem `json:"results"`
}

// PokemonListItem represents a single entry in the Pokemon list response.
type PokemonListItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
