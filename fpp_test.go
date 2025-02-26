package main

import (
	"testing"
)

func TestGenerateSpotItGeneric(t *testing.T) {
	tests := []struct {
		stickerCount int
		expectedN    int
		expectedSize int
	}{
		{7, 2, 7},
		{12, 2, 7},
		{13, 3, 13},
		{30, 3, 13},
		{31, 5, 31},
		{56, 5, 31},
		{57, 7, 57},
		{100, 7, 57},
	}

	for _, tt := range tests {
		stickers := make([]string, tt.stickerCount)
		for i := 0; i < tt.stickerCount; i++ {
			stickers[i] = string(rune('A' + i))
		}

		deck := GenerateSpotItGeneric(stickers)

		if len(deck) != tt.expectedSize {
			t.Errorf("For %d stickers, expected %d cards, got %d", tt.stickerCount, tt.expectedSize, len(deck))
		}

		symbolsPerCard := tt.expectedN + 1
		for i, card := range deck {
			if len(card) != symbolsPerCard {
				t.Errorf("Card %d has %d symbols, expected %d", i, len(card), symbolsPerCard)
			}
			seen := make(map[string]bool)
			for _, sym := range card {
				if seen[sym] {
					t.Errorf("Duplicate symbol %s on card %d", sym, i)
				}
				seen[sym] = true
			}
		}

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
					t.Errorf("For %d stickers, cards %d and %d have %d matches: %v vs %v", tt.stickerCount, i, j, matches, deck[i], deck[j])
				}
			}
		}
	}
}

func TestGenerateSpotItGenericTooFew(t *testing.T) {
	stickers := make([]string, 6)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for too few stickers, but none occurred")
		} else if r != "Need at least 7 stickers" {
			t.Errorf("Expected panic message 'Need at least 7 stickers', got %v", r)
		}
	}()
	GenerateSpotItGeneric(stickers)
}
