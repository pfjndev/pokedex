# Pokedex CLI

[![Release](https://github.com/pfjndev/pokedex/actions/workflows/release.yml/badge.svg)](https://github.com/pfjndev/pokedex/actions/workflows/release.yml)
[![Lint](https://github.com/pfjndev/pokedex/actions/workflows/lint.yml/badge.svg)](https://github.com/pfjndev/pokedex/actions/workflows/lint.yml)
[![Tests](https://github.com/pfjndev/pokedex/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/pfjndev/pokedex/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/badge/go-1.21+-blue.svg)](https://golang.org/dl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/pfjndev/pokedex)](https://goreportcard.com/report/github.com/pfjndev/pokedex)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A command-line Pokedex application written in Go. Query and explore Pokemon data with an interactive CLI interface.

## Features

- 🔍 Search Pokemon by name or ID
- 📊 View detailed Pokemon statistics
- 🎯 Interactive command-line interface
- ⚡ Fast and efficient data retrieval
- 🐍 Support for multiple Pokemon generations

## Prerequisites

- Go 1.21 or higher
- macOS, Linux, or Windows

## Installation

### From Source

```bash
git clone https://github.com/pfjndev/pokedex.git
cd pokedex
go build -o pokedex ./cmd/pokedex
```

### Run Directly

```bash
go run ./cmd/pokedex
```

## Usage

Start the interactive Pokedex:

```bash
./pokedex
```

Available commands:
- `help` - Display available commands
- `search <name>` - Search for a Pokemon by name
- `stat <id>` - View detailed stats for a Pokemon
- `exit` - Exit the application

## Project Structure

```
.
├── cmd/
│   └── pokedex/              # Main CLI application
├── internal/
│   ├── api/                  # API client for Pokemon data
│   ├── models/               # Data models
│   └── repl/                 # REPL implementation
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
└── README.md                 # This file
```

## Development

### Running Tests

```bash
go test ./...
```

### Code Quality

```bash
go fmt ./...
go vet ./...
golint ./...
```

## API Source

This project uses the [PokéAPI](https://pokeapi.co/) for Pokemon data.

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

[pfjndev](https://github.com/pfjndev)
