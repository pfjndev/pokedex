package cache

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestNewCache tests cache initialization
func TestNewCache(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	if c == nil {
		t.Fatal("Cache instance is nil")
	}

	// Verify cache directory was created
	cacheDir, err := GetCacheDir()
	if err != nil {
		t.Fatalf("Failed to get cache dir: %v", err)
	}

	pokemonDir := filepath.Join(cacheDir, PokemonDirName)
	if _, err := os.Stat(pokemonDir); os.IsNotExist(err) {
		t.Fatalf("Pokemon directory was not created: %s", pokemonDir)
	}
}

// TestSetAndGet tests writing and reading from cache
func TestSetAndGet(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	testData := []byte(`{"id": 1, "name": "Bulbasaur"}`)
	pokemonName := "bulbasaur"

	// Write to cache
	if err := c.Set(pokemonName, testData); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Read from cache
	cachedData, err := c.Get(pokemonName)
	if err != nil {
		t.Fatalf("Failed to get from cache: %v", err)
	}

	// Verify data matches
	if !bytes.Equal(testData, cachedData) {
		t.Fatalf("Cached data does not match: expected %s, got %s", string(testData), string(cachedData))
	}
}

// TestGetNonExistent tests retrieving non-existent Pokemon
func TestGetNonExistent(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	_, err = c.Get("nonexistent")
	if err == nil {
		t.Fatal("Expected error for non-existent Pokemon, got nil")
	}
}

// TestHas tests checking if Pokemon exists in cache
func TestHas(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	pokemonName := "pikachu"
	testData := []byte(`{"id": 25, "name": "Pikachu"}`)

	// Pokemon should not exist initially
	if c.Has(pokemonName) {
		t.Fatal("Pokemon should not exist in cache")
	}

	// Add Pokemon to cache
	if err := c.Set(pokemonName, testData); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Pokemon should now exist
	if !c.Has(pokemonName) {
		t.Fatal("Pokemon should exist in cache")
	}
}

// TestCaseInsensitivity tests that cache keys are case-insensitive
func TestCaseInsensitivity(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	testData := []byte(`{"id": 6, "name": "Charizard"}`)

	// Set with lowercase
	if err := c.Set("charizard", testData); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Get with uppercase
	cachedData, err := c.Get("CHARIZARD")
	if err != nil {
		t.Fatalf("Failed to get with uppercase: %v", err)
	}

	// Get with mixed case
	cachedData2, err := c.Get("ChArIzArD")
	if err != nil {
		t.Fatalf("Failed to get with mixed case: %v", err)
	}

	if !bytes.Equal(testData, cachedData) || !bytes.Equal(testData, cachedData2) {
		t.Fatal("Case-insensitive lookup failed")
	}
}

// TestDelete tests deleting a Pokemon from cache
func TestDelete(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	pokemonName := "squirtle"
	testData := []byte(`{"id": 7, "name": "Squirtle"}`)

	// Set Pokemon
	if err := c.Set(pokemonName, testData); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Verify it's cached
	if !c.Has(pokemonName) {
		t.Fatal("Pokemon should exist before deletion")
	}

	// Delete Pokemon
	if err := c.Delete(pokemonName); err != nil {
		t.Fatalf("Failed to delete: %v", err)
	}

	// Verify it's not cached anymore
	if c.Has(pokemonName) {
		t.Fatal("Pokemon should not exist after deletion")
	}

	// Verify Get returns error
	_, err = c.Get(pokemonName)
	if err == nil {
		t.Fatal("Expected error when getting deleted Pokemon")
	}
}

// TestList tests listing all cached Pokemon
func TestList(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	// Initially empty
	list, err := c.List()
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("Expected empty list, got %d items", len(list))
	}

	// Add some Pokemon
	pokemon := []string{"pikachu", "charizard", "blastoise"}
	testData := []byte(`{"id": 1}`)

	for _, p := range pokemon {
		if err := c.Set(p, testData); err != nil {
			t.Fatalf("Failed to set %s: %v", p, err)
		}
	}

	// List and verify
	list, err = c.List()
	if err != nil {
		t.Fatalf("Failed to list: %v", err)
	}

	if len(list) != len(pokemon) {
		t.Fatalf("Expected %d items, got %d", len(pokemon), len(list))
	}
}

// TestGetCacheStats tests getting cache statistics
func TestGetCacheStats(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	// Check stats on empty cache
	stats, err := c.GetCacheStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.FileCount != 0 || stats.TotalSize != 0 {
		t.Fatalf("Expected empty cache stats, got FileCount=%d, TotalSize=%d", stats.FileCount, stats.TotalSize)
	}

	// Add some data
	testData := []byte(`{"id": 1, "name": "Bulbasaur"}`)
	if err := c.Set("bulbasaur", testData); err != nil {
		t.Fatalf("Failed to set cache: %v", err)
	}

	// Check stats again
	stats, err = c.GetCacheStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.FileCount != 1 {
		t.Fatalf("Expected 1 file, got %d", stats.FileCount)
	}

	if stats.TotalSize != int64(len(testData)) {
		t.Fatalf("Expected total size %d, got %d", len(testData), stats.TotalSize)
	}
}

// TestClear tests clearing the entire cache
func TestClear(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	// Add multiple Pokemon
	pokemon := []string{"pikachu", "charizard", "blastoise", "venusaur"}
	testData := []byte(`{"id": 1}`)

	for _, p := range pokemon {
		if err := c.Set(p, testData); err != nil {
			t.Fatalf("Failed to set %s: %v", p, err)
		}
	}

	// Verify they're cached
	stats, err := c.GetCacheStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.FileCount != len(pokemon) {
		t.Fatalf("Expected %d files, got %d", len(pokemon), stats.FileCount)
	}

	// Clear cache
	if err := c.Clear(); err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify cache is empty
	for _, p := range pokemon {
		if c.Has(p) {
			t.Fatalf("Pokemon %s should not exist after clear", p)
		}
	}

	stats, err = c.GetCacheStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	if stats.FileCount != 0 {
		t.Fatalf("Expected empty cache after clear, got %d files", stats.FileCount)
	}
}

// TestEmptyData tests that setting empty data is rejected
func TestEmptyData(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	// Try to set empty data
	err = c.Set("test", []byte{})
	if err == nil {
		t.Fatal("Expected error when setting empty data")
	}

	// Try to set nil data
	err = c.Set("test", nil)
	if err == nil {
		t.Fatal("Expected error when setting nil data")
	}
}

// TestConcurrentAccess tests that cache is thread-safe
func TestConcurrentAccess(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	done := make(chan bool, 10)

	// Launch multiple goroutines doing concurrent operations
	for i := range 5 {
		go func(id int) {
			pokemonName := "pokemon" + string(rune('a'+id))
			testData := []byte(`{"id": ` + string(rune('0'+id)) + `}`)

			if err := c.Set(pokemonName, testData); err != nil {
				t.Errorf("Failed to set: %v", err)
			}

			if _, err := c.Get(pokemonName); err != nil {
				t.Errorf("Failed to get: %v", err)
			}

			if !c.Has(pokemonName) {
				t.Errorf("Pokemon not found after set")
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines
	for range 5 {
		<-done
	}
}

// TestInvalidInput tests handling of invalid input
func TestInvalidInput(t *testing.T) {
	c, err := NewCache()
	if err != nil {
		t.Fatalf("Failed to create cache: %v", err)
	}

	defer func() { _ = c.Clear() }()

	// Empty Pokemon name
	_, err = GetPokemonCachePath("")
	if err == nil {
		t.Fatal("Expected error for empty Pokemon name")
	}

	// Whitespace only
	_, err = GetPokemonCachePath("   ")
	if err == nil {
		t.Fatal("Expected error for whitespace-only Pokemon name")
	}

	// Try to set with empty name
	err = c.Set("", []byte("test"))
	if err == nil {
		t.Fatal("Expected error when setting with empty name")
	}
}
