package cache

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Cache provides a thread-safe file-based cache for Pokemon data
type Cache struct {
	mu sync.RWMutex // Protects access to cache operations
}

// NewCache creates a new Cache instance and ensures the cache directory exists
func NewCache() (*Cache, error) {
	if _, err := EnsureCacheDir(); err != nil {
		return nil, err
	}
	return &Cache{}, nil
}

// Get reads Pokemon data from the cache
// Returns the cached data as bytes or an error if not found or if reading fails
func (c *Cache) Get(nameOrID string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath, err := GetPokemonCachePath(nameOrID)
	if err != nil {
		return nil, err
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("pokemon %q not found in cache", nameOrID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to stat cache file: %w", err)
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file for %q: %w", nameOrID, err)
	}

	return data, nil
}

// Set writes Pokemon data to the cache using atomic writes
// Data is written to a temporary file and then renamed to ensure atomicity
func (c *Cache) Set(nameOrID string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(data) == 0 {
		return fmt.Errorf("data cannot be empty")
	}

	filePath, err := GetPokemonCachePath(nameOrID)
	if err != nil {
		return err
	}

	// Ensure the directory exists
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Write to a temporary file first (atomic operation)
	tempFile, err := os.CreateTemp(filepath.Dir(filePath), ".tmp-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %w", err)
	}
	tempPath := tempFile.Name()

	// Write data to temp file
	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Close the temp file before renaming
	if err := tempFile.Close(); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to close temporary file: %w", err)
	}

	// Atomically rename temp file to actual cache file
	if err := os.Rename(tempPath, filePath); err != nil {
		_ = os.Remove(tempPath)
		return fmt.Errorf("failed to write cache file for %q: %w", nameOrID, err)
	}

	return nil
}

// Has checks if a Pokemon exists in the cache
func (c *Cache) Has(nameOrID string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath, err := GetPokemonCachePath(nameOrID)
	if err != nil {
		return false
	}

	_, err = os.Stat(filePath)
	return err == nil
}

// Clear removes all cached Pokemon data
// Returns an error if the operation fails
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheDir, err := GetCacheDir()
	if err != nil {
		return err
	}

	pokemonDir := filepath.Join(cacheDir, PokemonDirName)

	// Check if directory exists before removing
	if _, err := os.Stat(pokemonDir); os.IsNotExist(err) {
		return nil // Nothing to clear
	}

	// Remove the entire pokemon directory
	if err := os.RemoveAll(pokemonDir); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	// Recreate the empty pokemon directory
	if err := os.MkdirAll(pokemonDir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	return nil
}

// CacheStats holds statistics about the cache
type CacheStats struct {
	TotalSize int64  // Total size of all cached files in bytes
	FileCount int    // Number of cached Pokemon files
	CacheDir  string // Path to the cache directory
}

// GetCacheStats returns statistics about the cache
// This includes total cache size and the number of cached Pokemon
func (c *Cache) GetCacheStats() (CacheStats, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheDir, err := GetCacheDir()
	if err != nil {
		return CacheStats{}, err
	}

	pokemonDir := filepath.Join(cacheDir, PokemonDirName)

	stats := CacheStats{
		CacheDir: cacheDir,
	}

	// Check if directory exists
	if _, err := os.Stat(pokemonDir); os.IsNotExist(err) {
		return stats, nil
	}

	// Walk through the pokemon directory
	err = filepath.Walk(pokemonDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only count files (not directories)
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
			stats.TotalSize += info.Size()
			stats.FileCount++
		}

		return nil
	})

	if err != nil {
		return CacheStats{}, fmt.Errorf("failed to calculate cache stats: %w", err)
	}

	return stats, nil
}

// CopyToFile copies cached Pokemon data to an external file
// This is useful for exporting cached data
func (c *Cache) CopyToFile(nameOrID string, destPath string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Read from cache
	data, err := c.Get(nameOrID)
	if err != nil {
		return err
	}

	// Write to destination file
	if err := os.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write to destination file: %w", err)
	}

	return nil
}

// LoadFromFile loads data from an external file and stores it in the cache
// This is useful for importing Pokemon data
func (c *Cache) LoadFromFile(nameOrID string, sourcePath string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Read from source file
	data, err := os.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	if len(data) == 0 {
		return fmt.Errorf("source file is empty")
	}

	// Write to cache
	return c.Set(nameOrID, data)
}

// List returns a list of all cached Pokemon names/IDs
func (c *Cache) List() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cacheDir, err := GetCacheDir()
	if err != nil {
		return nil, err
	}

	pokemonDir := filepath.Join(cacheDir, PokemonDirName)

	// Check if directory exists
	if _, err := os.Stat(pokemonDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	entries, err := os.ReadDir(pokemonDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache directory: %w", err)
	}

	var pokemon []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			// Remove .json extension
			name := strings.TrimSuffix(entry.Name(), ".json")
			pokemon = append(pokemon, name)
		}
	}

	return pokemon, nil
}

// Delete removes a specific Pokemon from the cache
func (c *Cache) Delete(nameOrID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	filePath, err := GetPokemonCachePath(nameOrID)
	if err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("pokemon %q not found in cache", nameOrID)
		}
		return fmt.Errorf("failed to delete cache file for %q: %w", nameOrID, err)
	}

	return nil
}
