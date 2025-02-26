package main

import (
	"log"
)

// GenerateSpotIt creates a finite projective plane PG(2,n) for n=7
func GenerateSpotIt(stickers []string) [][]string {
	const n = 7
	const totalCards = n*n + n + 1 // 57
	const symbolsPerCard = n + 1   // 8

	if len(stickers) < totalCards {
		log.Fatalf("Need at least %d stickers, got %d", totalCards, len(stickers))
	}

	// Step 1: Define points
	points := make([][3]int, totalCards)
	pIndex := 0
	for x := 0; x < n; x++ { // Finite points [x, y, 1]
		for y := 0; y < n; y++ {
			points[pIndex] = [3]int{x, y, 1}
			pIndex++
		}
	}
	for m := 0; m < n; m++ { // Slope infinities [1, m, 0]
		points[pIndex] = [3]int{1, m, 0}
		pIndex++
	}
	points[pIndex] = [3]int{0, 1, 0} // Vertical infinity [0, 1, 0]

	// Step 2: Generate lines
	cards := make([][]int, totalCards)
	cardIndex := 0

	// Helper: Get points on line ax + by + cz = 0 mod 7
	getLinePoints := func(a, b, c int) []int {
		line := make([]int, 0, symbolsPerCard)
		seen := make(map[int]bool)
		for i := 0; i < totalCards; i++ {
			x, y, z := points[i][0], points[i][1], points[i][2]
			if (a*x+b*y+c*z)%n == 0 && !seen[i] {
				line = append(line, i)
				seen[i] = true
				if len(line) == symbolsPerCard {
					break
				}
			}
		}
		if len(line) != symbolsPerCard {
			log.Printf("Line [%d, %d, %d] has %d points, expected %d", a, b, c, len(line), symbolsPerCard)
		}
		return line
	}

	// Generate 57 unique lines
	usedLines := make(map[[3]int]bool)
	for a := 0; a < n && cardIndex < totalCards; a++ {
		for b := 0; b < n && cardIndex < totalCards; b++ {
			for c := 0; c < n && cardIndex < totalCards; c++ {
				if a == 0 && b == 0 && c == 0 {
					continue // Invalid line
				}
				coeffs := [3]int{a, b, c}
				if usedLines[coeffs] {
					continue
				}
				// Normalize: Scale by inverse of first non-zero coefficient
				scale := 1
				if a != 0 {
					scale = modInverse(a, n)
				} else if b != 0 {
					scale = modInverse(b, n)
				} else {
					scale = modInverse(c, n)
				}
				norm := [3]int{(a * scale) % n, (b * scale) % n, (c * scale) % n}
				if usedLines[norm] {
					continue
				}
				usedLines[norm] = true

				cards[cardIndex] = getLinePoints(a, b, c)
				cardIndex++
			}
		}
	}

	if cardIndex != totalCards {
		log.Fatalf("Generated %d cards, expected %d", cardIndex, totalCards)
	}

	// Step 3: Map to stickers
	deck := make([][]string, totalCards)
	for i := 0; i < totalCards; i++ {
		deck[i] = make([]string, symbolsPerCard)
		seen := make(map[int]bool)
		for j, sym := range cards[i] {
			if sym >= totalCards {
				log.Printf("Out of range symbol %d on card %d", sym, i)
				sym = sym % totalCards
			}
			if seen[sym] {
				log.Printf("Duplicate symbol index %d (%s) on card %d", sym, stickers[sym], i)
			}
			seen[sym] = true
			deck[i][j] = stickers[sym]
		}
	}

	// Step 4: Verify
	for i := 0; i < totalCards-1; i++ {
		for j := i + 1; j < totalCards; j++ {
			matches := 0
			for _, s1 := range cards[i] {
				for _, s2 := range cards[j] {
					if s1 == s2 {
						matches++
					}
				}
			}
			if matches != 1 {
				log.Printf("Cards %d and %d have %d matches: %v vs %v", i, j, matches, cards[i], cards[j])
			}
		}
	}

	return deck
}

// modInverse computes the modular inverse in GF(n)
func modInverse(a, m int) int {
	a = a % m
	for x := 1; x < m; x++ {
		if (a*x)%m == 1 {
			return x
		}
	}
	return 1 // Shouldn’t happen for prime m
}
