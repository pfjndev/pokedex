# Pokemon Cache Package

A thread-safe, file-based caching system for Pokemon data written in Go.

## Overview

The `cache` package provides efficient local storage and retrieval of Pokemon data using the file system. Each Pokemon is stored as an individual JSON file in `~/.pokedex/cache/pokemon/`, making it easy to manage, inspect, and backup cached data.

## Features

- **Thread-Safe Operations**: All cache operations are protected by a read-write mutex (RWMutex) for safe concurrent access from multiple goroutines
- **Atomic Writes**: Data is written to temporary files and then atomically renamed, ensuring data integrity even during unexpected failures
- **Case-Insensitive Keys**: Pokemon names are automatically normalized to lowercase for consistent lookups
- **Efficient Lookups**: Direct file system access with minimal overhead
- **Comprehensive API**: Get, Set, Delete, Clear, List, and statistical queries
- **Error Handling**: Detailed error messages for debugging and recovery

## Directory Structure

```
~/.pokedex/cache/
└── pokemon/
    ├── pikachu.json
    ├── charizard.json
    ├── blastoise.json
    └── ...
```

## API Reference

### Cache Initialization

```go
// Create a new cache instance
cache, err := cache.NewCache()
if err != nil {
    log.Fatal(err)
}
```

### Core Operations

#### Get(nameOrID string) ([]byte, error)
Retrieves Pokemon data from the cache.

```go
data, err := cache.Get("pikachu")
if err != nil {
    log.Fatal("Pokemon not found in cache")
}
```

#### Set(nameOrID string, data []byte) error
Stores Pokemon data in the cache using atomic writes.

```go
data := []byte(`{"id": 25, "name": "Pikachu"}`)
if err := cache.Set("pikachu", data); err != nil {
    log.Fatal(err)
}
```

#### Has(nameOrID string) bool
Checks if a Pokemon exists in the cache without reading the full data.

```go
if cache.Has("pikachu") {
    fmt.Println("Pikachu is cached!")
}
```

#### Delete(nameOrID string) error
Removes a specific Pokemon from the cache.

```go
if err := cache.Delete("pikachu"); err != nil {
    log.Fatal(err)
}
```

### Cache Management

#### Clear() error
Removes all cached Pokemon data and recreates the cache directory.

```go
if err := cache.Clear(); err != nil {
    log.Fatal(err)
}
```

#### List() ([]string, error)
Returns a list of all cached Pokemon names/IDs.

```go
pokemon, err := cache.List()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Cached Pokemon: %v\n", pokemon)
```

#### GetCacheStats() (CacheStats, error)
Returns statistics about the cache including total size and file count.

```go
stats, err := cache.GetCacheStats()
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Cache size: %d bytes\n", stats.TotalSize)
fmt.Printf("Cached Pokemon: %d\n", stats.FileCount)
fmt.Printf("Cache location: %s\n", stats.CacheDir)
```

### File Import/Export

#### LoadFromFile(nameOrID string, sourcePath string) error
Imports Pokemon data from an external file into the cache.

```go
if err := cache.LoadFromFile("pikachu", "/path/to/pikachu.json"); err != nil {
    log.Fatal(err)
}
```

#### CopyToFile(nameOrID string, destPath string) error
Exports cached Pokemon data to an external file.

```go
if err := cache.CopyToFile("pikachu", "/path/to/backup/pikachu.json"); err != nil {
    log.Fatal(err)
}
```

## Storage Functions

The `storage.go` module provides utility functions for cache path management:

### GetCacheDir() (string, error)
Returns the path to the cache directory (`~/.pokedex/cache/`).

### EnsureCacheDir() (string, error)
Creates the cache directory structure if it doesn't exist.

### GetPokemonCachePath(nameOrID string) (string, error)
Returns the file path for a specific Pokemon's cache file.

### GetCachePathInfo() (CachePathInfo, error)
Returns information about the cache directory structure.

## Thread Safety

All cache operations are protected by a read-write mutex (RWMutex), allowing:
- Multiple concurrent reads via `Get()`, `Has()`, `List()`, and `GetCacheStats()`
- Exclusive write access during `Set()`, `Delete()`, `Clear()`, and `LoadFromFile()`

```go
// Safe for concurrent reads from multiple goroutines
go func() { cache.Get("pikachu") }()
go func() { cache.Has("charizard") }()
go func() { cache.List() }()

// Writes are serialized
go func() { cache.Set("blastoise", data) }()
```

## Error Handling

Common error scenarios:

```go
// Pokemon not found
data, err := cache.Get("nonexistent")
if err != nil {
    fmt.Println("Pokemon not found in cache")
}

// Empty data
err := cache.Set("pikachu", []byte{})
if err != nil {
    fmt.Println("Cannot cache empty data")
}

// Invalid input
_, err := GetPokemonCachePath("")
if err != nil {
    fmt.Println("Pokemon name cannot be empty")
}

// File system errors
stats, err := cache.GetCacheStats()
if err != nil {
    fmt.Println("Failed to calculate cache stats:", err)
}
```

## Example Usage

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "github.com/pfjndev/pokedex/internal/cache"
)

type Pokemon struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Type string `json:"type"`
}

func main() {
    // Initialize cache
    c, err := cache.NewCache()
    if err != nil {
        log.Fatalf("Failed to initialize cache: %v", err)
    }

    // Create and cache a Pokemon
    pikachu := Pokemon{ID: 25, Name: "Pikachu", Type: "Electric"}
    data, err := json.Marshal(pikachu)
    if err != nil {
        log.Fatal(err)
    }

    if err := c.Set("pikachu", data); err != nil {
        log.Fatalf("Failed to cache Pokemon: %v", err)
    }

    // Retrieve from cache
    if cachedData, err := c.Get("pikachu"); err == nil {
        var cached Pokemon
        json.Unmarshal(cachedData, &cached)
        fmt.Printf("Retrieved: %+v\n", cached)
    }

    // Check cache statistics
    stats, err := c.GetCacheStats()
    if err == nil {
        fmt.Printf("Cache contains %d Pokemon using %d bytes\n", 
                   stats.FileCount, stats.TotalSize)
    }

    // List all cached Pokemon
    list, err := c.List()
    if err == nil {
        fmt.Printf("Cached: %v\n", list)
    }

    // Clear cache when done
    if err := c.Clear(); err != nil {
        log.Fatal(err)
    }
}
```

## Performance Considerations

- **Read-Heavy Workloads**: The RWMutex allows concurrent reads, improving performance for read-heavy operations
- **Large Cache Files**: Consider splitting very large Pokemon datasets into multiple cache entries
- **Disk I/O**: File system operations are relatively fast for JSON files, but batch operations when possible
- **Cache Cleanup**: Periodically call `GetCacheStats()` to monitor cache size and clear if needed

## Testing

The package includes comprehensive unit tests covering:
- Basic get/set operations
- Case-insensitive lookups
- Concurrent access
- Cache statistics
- Error handling
- Data integrity

Run tests with:
```bash
go test ./internal/cache -v
```

## Future Enhancements

Potential improvements for future versions:
- TTL (Time-To-Live) for cached entries
- Compression support for large cache files
- Batch operations for improved performance
- Cache eviction policies (LRU, LFU)
- Encryption for sensitive data
- Migration tools for cache versions
