package models

import (
	"math"
	"testing"
)

func TestGetStatByName(t *testing.T) {
pikachu := &Pokemon{
Stats: []*Stat{
{Name: "hp", BaseValue: 35},
{Name: "attack", BaseValue: 55},
{Name: "defense", BaseValue: 40},
},
}

// Test case-insensitive matching
stat := pikachu.GetStatByName("HP")
if stat == nil || stat.BaseValue != 35 {
t.Fatalf("Expected HP stat with value 35, got %v", stat)
}

stat = pikachu.GetStatByName("hp")
if stat == nil || stat.BaseValue != 35 {
t.Fatalf("Expected hp stat with value 35, got %v", stat)
}

stat = pikachu.GetStatByName("ATTACK")
if stat == nil || stat.BaseValue != 55 {
t.Fatalf("Expected ATTACK stat with value 55, got %v", stat)
}

// Test non-existent stat
stat = pikachu.GetStatByName("magic")
if stat != nil {
t.Fatalf("Expected nil for non-existent stat, got %v", stat)
}

// Test nil Pokemon
var nilPokemon *Pokemon
stat = nilPokemon.GetStatByName("hp")
if stat != nil {
t.Fatalf("Expected nil for nil Pokemon, got %v", stat)
}
}

func TestGetStatPercentage(t *testing.T) {
pikachu := &Pokemon{
Stats: []*Stat{
{Name: "hp", BaseValue: 35},
{Name: "speed", BaseValue: 90},
{Name: "defense", BaseValue: 255}, // Max value
},
}

// Test HP percentage
percent := pikachu.GetStatPercentage("hp")
expected := (35.0 / 255.0) * 100.0
if math.Abs(percent-expected) > 1e-9 {
t.Fatalf("Expected %.2f%%, got %.2f%%", expected, percent)
}

// Test Speed percentage
percent = pikachu.GetStatPercentage("speed")
expected = (90.0 / 255.0) * 100.0
if math.Abs(percent-expected) > 1e-9 {
t.Fatalf("Expected %.2f%%, got %.2f%%", expected, percent)
}

// Test max value
percent = pikachu.GetStatPercentage("defense")
if math.Abs(percent-100.0) > 1e-9 {
t.Fatalf("Expected 100%%, got %.2f%%", percent)
}

// Test non-existent stat
percent = pikachu.GetStatPercentage("magic")
if percent != 0.0 {
t.Fatalf("Expected 0.0 for non-existent stat, got %.2f", percent)
}

// Test nil Pokemon
var nilPokemon *Pokemon
percent = nilPokemon.GetStatPercentage("hp")
if percent != 0.0 {
t.Fatalf("Expected 0.0 for nil Pokemon, got %.2f", percent)
}
}

func TestGetStatColor(t *testing.T) {
pikachu := &Pokemon{}

// Test low range (0-85)
color := pikachu.GetStatColor(0)
if color != "#E74C3C" {
t.Fatalf("Expected #E74C3C for value 0, got %s", color)
}

color = pikachu.GetStatColor(85)
if color != "#E74C3C" {
t.Fatalf("Expected #E74C3C for value 85, got %s", color)
}

// Test medium range (86-170)
color = pikachu.GetStatColor(86)
if color != "#F39C12" {
t.Fatalf("Expected #F39C12 for value 86, got %s", color)
}

color = pikachu.GetStatColor(170)
if color != "#F39C12" {
t.Fatalf("Expected #F39C12 for value 170, got %s", color)
}

// Test high range (171-255)
color = pikachu.GetStatColor(171)
if color != "#27AE60" {
t.Fatalf("Expected #27AE60 for value 171, got %s", color)
}

color = pikachu.GetStatColor(255)
if color != "#27AE60" {
t.Fatalf("Expected #27AE60 for value 255, got %s", color)
}
}

func TestGetTotalStats(t *testing.T) {
pikachu := &Pokemon{
Stats: []*Stat{
{Name: "hp", BaseValue: 35},
{Name: "attack", BaseValue: 55},
{Name: "defense", BaseValue: 40},
{Name: "special-attack", BaseValue: 50},
{Name: "special-defense", BaseValue: 50},
{Name: "speed", BaseValue: 90},
},
}

total := pikachu.GetTotalStats()
expected := 35 + 55 + 40 + 50 + 50 + 90
if total != expected {
t.Fatalf("Expected total %d, got %d", expected, total)
}

// Test nil Pokemon
var nilPokemon *Pokemon
total = nilPokemon.GetTotalStats()
if total != 0 {
t.Fatalf("Expected 0 for nil Pokemon, got %d", total)
}
}

func TestGetAverageStatPercentage(t *testing.T) {
pikachu := &Pokemon{
Stats: []*Stat{
{Name: "hp", BaseValue: 35},
{Name: "attack", BaseValue: 55},
{Name: "defense", BaseValue: 40},
{Name: "special-attack", BaseValue: 50},
{Name: "special-defense", BaseValue: 50},
{Name: "speed", BaseValue: 90},
},
}

avgPercent := pikachu.GetAverageStatPercentage()

// Calculate expected: sum of all percentages / count
total := 0.0
for _, stat := range pikachu.Stats {
total += (float64(stat.BaseValue) / 255.0) * 100.0
}
expected := total / float64(len(pikachu.Stats))

if avgPercent != expected {
t.Fatalf("Expected %.2f%%, got %.2f%%", expected, avgPercent)
}

// Test nil Pokemon
var nilPokemon *Pokemon
avgPercent = nilPokemon.GetAverageStatPercentage()
if avgPercent != 0.0 {
t.Fatalf("Expected 0.0 for nil Pokemon, got %.2f", avgPercent)
}

// Test empty stats
emptyPokemon := &Pokemon{Stats: []*Stat{}}
avgPercent = emptyPokemon.GetAverageStatPercentage()
if avgPercent != 0.0 {
t.Fatalf("Expected 0.0 for empty stats, got %.2f", avgPercent)
}
}
