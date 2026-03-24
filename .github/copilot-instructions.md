# Copilot Instructions for Pokedex CLI

A command-line Pokedex application written in Go that queries Pokemon data from the PokeAPI.

## Build, Test, and Lint Commands

### Build
```bash
make build          # Build binary to ./bin/pokedex
go build -o ./bin/pokedex ./cmd/pokedex
```

### Run
```bash
make run            # Build and run
go run ./cmd/pokedex
```

### Testing
```bash
# Run all tests
make test
go test ./...

# Run tests for a specific package
go test ./internal/api
go test ./internal/cache

# Run tests with verbose output
go test -v ./...

# Run a single test
go test -v ./internal/api -run TestNewClient
go test -v ./internal/cache -run TestSetAndGet
```

### Linting
```bash
make lint           # Prefers golangci-lint, falls back to go vet
golangci-lint run ./...
go vet ./...
```

### Clean
```bash
make clean          # Remove build artifacts and run go mod tidy
```

## Architecture Overview

### Three-Layer Architecture

1. **API Layer** (`internal/api/`)
   - HTTP client for PokeAPI v2
   - Exponential backoff retry logic with jitter
   - Structured error handling with `APIError` type
   - Type mappings from PokeAPI JSON responses

2. **Cache Layer** (`internal/cache/`)
   - Thread-safe file-based cache using `sync.RWMutex`
   - Atomic writes via temp files + rename for data integrity
   - Pokemon data stored as JSON in `~/.cache/pokedex/pokemon/`
   - Cache operations: Get, Set, Has, Delete, Clear, List, Stats

3. **Models Layer** (`internal/models/`)
   - Clean domain models focused on UI needs
   - Rich helper methods for formatting and calculations
   - Type colors for UI rendering

### Data Flow

```
User Request → API Client → PokeAPI
                    ↓
           PokemonResponse (API types)
                    ↓
           Transform to models.Pokemon
                    ↓
              Cache → File System
```

### Key Design Patterns

**Functional Options Pattern** - Used in `api.NewClient()`:
```go
client := api.NewClient(
    api.WithTimeout(15 * time.Second),
    api.WithMaxRetries(5),
)
```

**Atomic File Writes** - Cache writes to temp file first, then atomic rename:
```go
tempFile → Write data → Close → Rename to final path
```

**Retry Logic** - Exponential backoff with jitter:
- Initial backoff: 100ms
- Max backoff: 10s
- Default retries: 3
- Does NOT retry on 404 errors

## Codebase Conventions

### Error Handling
- API errors wrapped in `APIError` struct with StatusCode + Message
- Cache operations return descriptive wrapped errors: `fmt.Errorf("failed to X: %w", err)`
- Test for specific error conditions before retrying (e.g., 404 = no retry)

### Testing
- Table-driven tests NOT used - each test is a separate function
- Test functions named `TestFunctionName` (e.g., `TestNewClient`, `TestSetAndGet`)
- Use `t.Fatalf()` for fatal errors, NOT `t.Errorf()` followed by return
- Clean up test data with `defer func() { _ = c.Clear() }()`
- Test concurrency with goroutines where applicable

### Code Organization
- `internal/` for unexported packages (standard Go practice)
- Separate `types.go` for API response structures
- Doc comments follow Go conventions: `// FunctionName does X` on exported functions
- Constants grouped at package level with descriptive comments

### Naming Conventions
- Use `nameOrID` for parameters accepting Pokemon name or ID
- Cache files named `<name-or-id>.json` (lowercase)
- Mutex named `mu` for brevity
- HTTP client named `httpClient` to distinguish from API `Client`

### File Operations
- Always use `filepath.Join()` for cross-platform path construction
- Check file existence with `os.Stat()` and `os.IsNotExist(err)`
- Use `os.ReadFile()` and `os.WriteFile()` for simple file I/O
- Set directory permissions to `0755`, file permissions to `0644`

### Concurrency
- Use `sync.RWMutex` for cache operations (reads can be concurrent)
- Lock patterns: `c.mu.Lock()` / `defer c.mu.Unlock()` for writes
- Lock patterns: `c.mu.RLock()` / `defer c.mu.RUnlock()` for reads

### Models Layer Conventions
- Rich domain models with helper methods (not anemic data classes)
- Formatting methods return "Unknown" for zero values
- Percentage/stats calculations use well-documented maximums (e.g., 255 for base stats)
- Methods are nil-safe: check `if p == nil` before accessing fields

### Sprite Display System

**ASCII Art Generation** (`internal/models/pokemon.go`):
- `GetSpriteASCII()` method generates 10x10 character ASCII art representation of Pokemon
- Pattern selection based on Pokemon primary type for visual variety
- Type-to-pattern mappings:
  - `⚡` Electric (lightning energy)
  - `▓` Fire (flame intensity)
  - `≈` Water (wave patterns)
  - `❀` Grass (flower/plant)
  - `❄` Ice (snowflake pattern)
  - `◆` Dragon (gem/crystal pattern)
  - `▒` Ghost (ethereal pattern)
  - `█` Default (solid blocks for unknown types)
- Fallback: Uses solid blocks (`█`) if type not recognized
- Returns "No sprite available" if Pokemon has no sprites

**Sprite Preference Order** (`BestSprite()` method):
1. Official Artwork (475x475px, highest quality professional render)
2. Front Shiny (color variant, usually matches official style)
3. Front Default (standard sprite, fallback)
- Returns first non-empty sprite in order of preference
- Returns empty string if no sprites available

**Data Flow & Extraction**:
- PokeAPI provides sprites in nested JSON structure under `sprites` object
- Service layer (`internal/api/service.go`) extracts all 9 sprite variants during transform
- Official Artwork located at: `sprites.other["official-artwork"].front_default`
- Domain model stores all variants in `Sprite` struct:
  ```go
  type Sprite struct {
      FrontDefault string // standard sprite URL
      FrontShiny   string // shiny variant URL
      BackDefault  string // rear-facing sprite
      BackShiny    string // rear shiny variant
      OfficialArtwork string // 475x475px professional artwork
      // ... additional variants
  }
  ```

**Sprite Caching**:
- Sprites stored as URLs (not cached as images)
- Full Pokemon data including sprite URLs cached by `internal/cache/`
- Cache key: Pokemon name or ID (e.g., `pikachu.json`)
- Cache location: `~/.cache/pokedex/pokemon/`

**UI Integration** (`internal/ui/views.go`):
- `renderDetails()` function displays sprite information:
  - Calls `pokemon.GetSpriteASCII()` for ASCII art display
  - Shows ASCII art in terminal above Pokemon stats section
  - Displays official artwork URL with `BestSprite()` for reference
  - Gracefully handles nil sprites with "(No sprite available)" message
  - ASCII art useful for low-bandwidth scenarios or terminal-only viewing

### API Client Conventions
- Always set `User-Agent` header: `"Pokedex CLI (+https://github.com/pfjndev/pokedex)"`
- Accept JSON: `req.Header.Set("Accept", "application/json")`
- Read full response body before checking errors
- Close response body with `defer func() { _ = resp.Body.Close() }()`

### Constants
```go
const (
    DefaultBaseURL = "https://pokeapi.co/api/v2"
    DefaultTimeout = 10 * time.Second
    DefaultMaxRetries = 3
    InitialBackoff = 100 * time.Millisecond
    MaxBackoff = 10 * time.Second
)
```

## CI/CD

### GitHub Actions Workflows
- `.github/workflows/test.yml` - Run tests on push/PR
- `.github/workflows/lint.yml` - Run linters on push/PR
- `.github/workflows/release.yml` - Release automation
- `.github/workflows/dependabot-auto-merge.yml` - Auto-merge Dependabot PRs

### Dependabot Configuration
- Weekly updates on Mondays (09:00 UTC for Actions, 10:00 UTC for Go modules)
- Auto-merges patch and minor updates
- Major updates require manual review
- Groups: `bubbletea` (charmbracelet deps) and `test-deps`

## External API

**PokeAPI v2**: https://pokeapi.co/api/v2
- Endpoints used:
  - `GET /pokemon/{name-or-id}` - Fetch single Pokemon
  - `GET /pokemon?limit={limit}&offset={offset}` - List Pokemon (paginated)
- Rate limiting: 429 status code triggers retry with backoff
- No authentication required

## Project Status

**Current state**: Foundation complete, CLI interface not yet implemented
- ✅ API client with retry logic
- ✅ File-based cache with atomic writes
- ✅ Domain models with helper methods
- ⏳ CLI/REPL implementation pending (see `cmd/pokedex/main.go`)

## Module Information

- **Module**: `github.com/pfjndev/pokedex`
- **Go Version**: 1.26.1
- **Dependencies**: Standard library only (no external dependencies)
