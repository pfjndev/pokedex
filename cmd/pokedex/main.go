package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pfjndev/pokedex/internal/api"
	"github.com/pfjndev/pokedex/internal/cache"
	"github.com/pfjndev/pokedex/internal/service"
	"github.com/pfjndev/pokedex/internal/ui"
)

var (
	version   = "dev"
	showVer   = flag.Bool("version", false, "Show version information")
	cacheDir  = flag.String("cache-dir", "", "Custom cache directory (default: ~/.cache/pokedex)")
	noCache   = flag.Bool("no-cache", false, "Disable caching")
)

func main() {
	flag.Parse()

	// Show version
	if *showVer {
		fmt.Printf("Pokedex CLI %s\n", version)
		os.Exit(0)
	}

	// Initialize API client
	apiClient := api.NewClient()

	// Initialize cache
	var pokemonCache service.Cache
	if *noCache {
		pokemonCache = &noOpCache{}
	} else {
		c, err := cache.NewCache()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing cache: %v\n", err)
			fmt.Fprintf(os.Stderr, "Continuing without cache...\n")
			pokemonCache = &noOpCache{}
		} else {
			pokemonCache = c
		}
	}

	// Initialize service
	svc := service.NewService(apiClient, pokemonCache)

	// Create and run Bubble Tea program
	model := ui.NewModel(svc)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

// noOpCache is a cache implementation that does nothing
type noOpCache struct{}

func (n *noOpCache) Get(nameOrID string) ([]byte, error) {
	return nil, fmt.Errorf("cache disabled")
}

func (n *noOpCache) Set(nameOrID string, data []byte) error {
	return nil
}

