# Pokedex CLI

[![Release](https://github.com/pfjndev/pokedex/actions/workflows/release.yml/badge.svg)](https://github.com/pfjndev/pokedex/actions/workflows/release.yml)
[![Lint](https://github.com/pfjndev/pokedex/actions/workflows/lint.yml/badge.svg)](https://github.com/pfjndev/pokedex/actions/workflows/lint.yml)
[![Tests](https://github.com/pfjndev/pokedex/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/pfjndev/pokedex/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/pfjndev/pokedex)](https://goreportcard.com/report/github.com/pfjndev/pokedex)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A beautiful command-line Pokedex application written in Go with an interactive terminal UI built using Bubble Tea. Search and explore Pokemon data with style!

## ✨ Features

- 🔍 **Search Pokemon** - Find Pokemon by name or ID with smart search
- 📋 **Browse Pokemon** - Navigate through paginated lists of all Pokemon
- 📊 **Detailed Stats** - View comprehensive Pokemon information including:
  - Types with color-coded badges
  - Base stats with visual bars
  - Abilities (including hidden abilities)
  - Physical characteristics (height/weight)
  - Base experience
- 🎨 **ASCII Sprite Display** - ASCII art rendering with type-specific visual patterns
- 🖼️ **Official Artwork URLs** - Links to high-resolution sprites for external viewing
- 🎨 **Beautiful TUI** - Interactive terminal interface with Bubble Tea
- ⚡ **Fast & Cached** - File-based caching for instant repeated lookups
- 🛡️ **Reliable** - Automatic retry logic with exponential backoff
- 🌐 **Complete Data** - Powered by [PokéAPI](https://pokeapi.co/) v2

## 🎨 Sprite Display

The Pokedex displays Pokemon sprites as ASCII art directly in your terminal! Each Pokemon type has unique visual patterns that help identify the Pokemon at a glance:

**Type-Specific Sprite Patterns:**
- ⚡ **Electric** - Lightning bolt patterns
- 🔥 **Fire** - Flame patterns (▓)
- 💧 **Water** - Wave patterns (≈)
- 🌿 **Grass** - Flower patterns (❀)
- ❄️ **Ice** - Snowflake patterns (❄)
- 🐉 **Dragon** - Scale patterns (◆)
- 👻 **Ghost** - Spooky patterns (▒)
- ⚪ **Normal/Other** - Solid blocks (█)

When viewing Pokemon details, you'll see:
- **ASCII Sprite** - A small ASCII art representation of the Pokemon in the terminal
- **Official Artwork URL** - A direct link to the official Pokemon artwork so you can view high-resolution sprites in your browser

This combines the best of both worlds: quick ASCII visualization in your terminal and access to official artwork when you want more detail!

## 📋 Prerequisites

- Go 1.26.1 or higher
- macOS, Linux, or Windows
- Terminal with color support

## 🚀 Installation

### From Source

```bash
git clone https://github.com/pfjndev/pokedex.git
cd pokedex
make build
```

The binary will be available at `./bin/pokedex`.

### Using Go Install

```bash
go install github.com/pfjndev/pokedex/cmd/pokedex@latest
```

## 🎮 Usage

### Start the Interactive CLI

```bash
./bin/pokedex
# or if installed globally
pokedex
```

### Command-Line Options

```bash
pokedex [options]

Options:
  --version          Show version information
  --cache-dir PATH   Custom cache directory (default: ~/.cache/pokedex)
  --no-cache         Disable caching
  --help             Show help message
```

### Keyboard Shortcuts

#### Main Menu
- `s` - Search for a Pokemon
- `b` - Browse Pokemon list
- `q` - Quit application

#### Search View
- Type to search for a Pokemon
- `Enter` - Execute search
- `Esc` - Return to main menu

#### Browse List
- `↑/k` - Move cursor up
- `↓/j` - Move cursor down
- `Enter` - View Pokemon details
- `n` - Next page
- `p` - Previous page
- `q/Esc` - Return to main menu

#### Details View
- `Backspace/Esc/q` - Return to list

#### Global
- `Ctrl+C` - Force quit

## 📁 Project Structure

```
.
├── cmd/
│   └── pokedex/              # Main CLI application entry point
├── internal/
│   ├── api/                  # API client for PokeAPI v2
│   │   ├── client.go         # HTTP client with retry logic
│   │   ├── types.go          # API response types
│   │   └── client_test.go    # API client tests
│   ├── cache/                # File-based caching layer
│   │   ├── cache.go          # Cache operations
│   │   ├── path.go           # Path utilities
│   │   └── cache_test.go     # Cache tests
│   ├── models/               # Domain models
│   │   └── pokemon.go        # Pokemon data structures
│   ├── service/              # Business logic layer
│   │   ├── service.go        # Service coordination
│   │   └── service_test.go   # Service tests
│   └── ui/                   # Terminal UI
│       ├── app.go            # Bubble Tea application
│       └── views.go          # UI view rendering
├── Makefile                  # Build commands
├── go.mod                    # Go module definition
└── README.md                 # This file
```

## 🛠️ Development

### Build

```bash
make build          # Build binary to ./bin/pokedex
```

### Run

```bash
make run            # Build and run
```

### Test

```bash
make test           # Run all tests
go test ./...       # Alternative

# Run specific package tests
go test ./internal/api
go test ./internal/cache
go test ./internal/service

# Run with coverage
go test -cover ./...
```

### Lint

```bash
make lint           # Run linter
```

### Clean

```bash
make clean          # Remove build artifacts
```

### Cache Management

The Pokedex caches Pokemon data in `~/.cache/pokedex/pokemon/` by default.

To clear the cache:
```bash
rm -rf ~/.cache/pokedex
```

## 🏗️ Architecture

The application follows a clean three-layer architecture:

1. **API Layer** (`internal/api/`)
   - HTTP client for PokeAPI v2
   - Exponential backoff retry logic
   - Structured error handling

2. **Cache Layer** (`internal/cache/`)
   - Thread-safe file-based cache
   - Atomic writes for data integrity
   - JSON storage format

3. **Service Layer** (`internal/service/`)
   - Coordinates API and cache
   - Cache-first strategy with API fallback
   - Domain model transformation

4. **UI Layer** (`internal/ui/`)
   - Bubble Tea TUI framework
   - Interactive navigation
   - State management

### Data Flow

```
User Input → UI (Bubble Tea) → Service → Cache → API → PokeAPI
                                   ↓
                            Domain Models
```

## 🌐 API Source

This project uses the [PokéAPI](https://pokeapi.co/) v2 for Pokemon data. No authentication required.

## 🐛 Troubleshooting

### Cache Issues

If you experience cache-related errors:
```bash
./bin/pokedex --no-cache
# or
rm -rf ~/.cache/pokedex
```

### API Rate Limiting

The application implements automatic retry with exponential backoff for rate limiting (429 errors).

### Network Errors

Check your internet connection. The app requires API access for first-time Pokemon lookups.

## 🤝 Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👤 Author

[pfjndev](https://github.com/pfjndev)

## 🙏 Acknowledgments

- [PokéAPI](https://pokeapi.co/) - For providing the Pokemon data
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - For the amazing TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - For beautiful terminal styling

