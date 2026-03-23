// Package cache provides a thread-safe, file-based caching system for Pokemon data.
//
// The cache stores each Pokemon as an individual JSON file in ~/.pokedex/cache/pokemon/.
// Cache files are named using either Pokemon IDs or names (normalized to lowercase).
//
// Example usage:
//
//	package main
//
//	import (
//		"fmt"
//		"log"
//		"github.com/pfjndev/pokedex/internal/cache"
//	)
//
//	func main() {
//		// Create a new cache instance
//		c, err := cache.NewCache()
//		if err != nil {
//			log.Fatalf("Failed to create cache: %v", err)
//		}
//
//		// Store Pokemon data
//		data := []byte(`{"id": 1, "name": "Bulbasaur"}`)
//		if err := c.Set("bulbasaur", data); err != nil {
//			log.Fatalf("Failed to cache Pokemon: %v", err)
//		}
//
//		// Retrieve Pokemon data
//		cached, err := c.Get("bulbasaur")
//		if err != nil {
//			log.Fatalf("Pokemon not found in cache: %v", err)
//		}
//		fmt.Printf("Cached data: %s\n", string(cached))
//
//		// Check if Pokemon is cached
//		if c.Has("bulbasaur") {
//			fmt.Println("Bulbasaur is cached!")
//		}
//
//		// Get cache statistics
//		stats, err := c.GetCacheStats()
//		if err != nil {
//			log.Fatalf("Failed to get cache stats: %v", err)
//		}
//		fmt.Printf("Cache size: %d bytes, Files: %d\n", stats.TotalSize, stats.FileCount)
//
//		// List all cached Pokemon
//		pokemon, err := c.List()
//		if err != nil {
//			log.Fatalf("Failed to list cached Pokemon: %v", err)
//		}
//		fmt.Printf("Cached Pokemon: %v\n", pokemon)
//
//		// Delete a Pokemon from cache
//		if err := c.Delete("bulbasaur"); err != nil {
//			log.Fatalf("Failed to delete from cache: %v", err)
//		}
//
//		// Clear entire cache
//		if err := c.Clear(); err != nil {
//			log.Fatalf("Failed to clear cache: %v", err)
//		}
//	}
//
// Cache Features:
//
// Thread Safety: All cache operations are protected by a read-write mutex (RWMutex),
// ensuring safe concurrent access from multiple goroutines.
//
// Atomic Writes: Data is written to a temporary file and then atomically renamed
// to the final location, preventing partial writes in case of failures.
//
// Key Normalization: Pokemon names are automatically normalized to lowercase,
// allowing case-insensitive lookups (e.g., "Pikachu", "pikachu", and "PIKACHU" all map to the same file).
//
// Directory Structure:
//
//	~/.pokedex/cache/
//	├── pokemon/
//	│   ├── pikachu.json
//	│   ├── charizard.json
//	│   ├── blastoise.json
//	│   └── ...
//
// Cache Path Resolution:
//
// The cache directory is automatically created in the user's home directory if it doesn't exist.
// Cache files are stored as individual JSON files named after the Pokemon ID or name.
//
// Error Handling:
//
// All operations return error values that should be checked. Common error scenarios include:
// - File not found in cache (Get)
// - Permission denied (file system access errors)
// - Disk full (write operations)
// - Invalid input (empty Pokemon name/ID)
package cache
