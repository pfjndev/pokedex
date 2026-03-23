package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Cache directory constants
const (
	// CacheDirName is the name of the cache directory
	CacheDirName = ".pokedex"
	// PokemonDirName is the subdirectory for Pokemon cache
	PokemonDirName = "pokemon"
)

// GetCacheDir returns the path to the cache directory (~/.pokedex/cache/)
func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	cacheDir := filepath.Join(homeDir, CacheDirName, "cache")
	return cacheDir, nil
}

// EnsureCacheDir creates the cache directory structure if it doesn't exist
func EnsureCacheDir() (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}

	pokemonDir := filepath.Join(cacheDir, PokemonDirName)
	if err := os.MkdirAll(pokemonDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory %s: %w", pokemonDir, err)
	}

	return cacheDir, nil
}

// GetPokemonCachePath returns the path to a Pokemon's cache file
// nameOrID can be either a numeric ID or a Pokemon name
// Names are normalized to lowercase
func GetPokemonCachePath(nameOrID string) (string, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return "", err
	}

	// Normalize the key to lowercase
	key := strings.ToLower(strings.TrimSpace(nameOrID))

	// Ensure the key is not empty
	if key == "" {
		return "", fmt.Errorf("nameOrID cannot be empty")
	}

	// Create the file path
	filePath := filepath.Join(cacheDir, PokemonDirName, fmt.Sprintf("%s.json", key))
	return filePath, nil
}

// CachePathInfo holds information about cache files
type CachePathInfo struct {
	Dir        string // Cache directory path
	PokemonDir string // Pokemon subdirectory path
}

// GetCachePathInfo returns information about the cache directory structure
func GetCachePathInfo() (CachePathInfo, error) {
	cacheDir, err := GetCacheDir()
	if err != nil {
		return CachePathInfo{}, err
	}

	return CachePathInfo{
		Dir:        cacheDir,
		PokemonDir: filepath.Join(cacheDir, PokemonDirName),
	}, nil
}
