package main

import (
	"testing"
)

// TestGenerateSpotIt verifies the projective plane properties
func TestGenerateSpotIt(t *testing.T) {
	// Mock input: 57 stickers
	stickers := make([]string, 57)
	for i := 0; i < 57; i++ {
		stickers[i] = string(rune('A' + i))
	}

	deck := GenerateSpotIt(stickers)

	// Check 1: Correct number of cards
	if len(deck) != 57 {
		t.Errorf("Expected 57 cards, got %d", len(deck))
	}

	// Check 2: Each card has 8 unique symbols
	for i, card := range deck {
		if len(card) != 8 {
			t.Errorf("Card %d has %d symbols, expected 8", i, len(card))
		}
		seen := make(map[string]bool)
		for _, sym := range card {
			if seen[sym] {
				t.Errorf("Duplicate symbol %s on card %d", sym, i)
			}
			seen[sym] = true
		}
	}

	// Check 3: Exactly 1 match between each pair
	for i := 0; i < len(deck)-1; i++ {
		for j := i + 1; j < len(deck); j++ {
			matches := 0
			for _, s1 := range deck[i] {
				for _, s2 := range deck[j] {
					if s1 == s2 {
						matches++
					}
				}
			}
			if matches != 1 {
				t.Errorf("Cards %d and %d have %d matches, expected 1", i, j, matches)
			}
		}
	}

	// Check 4: All symbols are within range (0-56)
	for i, card := range deck {
		for _, sym := range card {
			found := false
			for _, s := range stickers[:57] {
				if s == sym {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Card %d has out-of-range symbol %s", i, sym)
			}
		}
	}
}

