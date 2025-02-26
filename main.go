package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const uploadDir = "uploads"

var stickers = make([]string, 57)

func generateSpotIt() [][]string {
	const n = 7
	const totalCards = n*n + n + 1 // 57
	const symbolsPerCard = n + 1   // 8

	// Points: 0-48 are [x, y, 1], 49-55 are [1, m, 0], 56 is [0, 1, 0]
	points := make([][3]int, totalCards)
	pIndex := 0
	for x := 0; x < n; x++ {
		for y := 0; y < n; y++ {
			points[pIndex] = [3]int{x, y, 1}
			pIndex++
		}
	}
	for m := 0; m < n; m++ {
		points[pIndex] = [3]int{1, m, 0}
		pIndex++
	}
	points[pIndex] = [3]int{0, 1, 0}

	// Lines
	cards := make([][]int, totalCards)
	cardIndex := 0

	// Step 1: Lines through [0, 0, 1] with slopes 0-6
	for m := 0; m < n; m++ {
		card := make([]int, symbolsPerCard)
		card[0] = n*n + m // [1, m, 0], 49-55
		for x := 0; x < n; x++ {
			y := (m * x) % n
			card[x+1] = x + y*n
		}
		cards[cardIndex] = card
		cardIndex++
	}
	// Vertical line through [0, 0, 1]
	card := make([]int, symbolsPerCard)
	card[0] = n*n + n // [0, 1, 0], 56
	for y := 0; y < n; y++ {
		card[y+1] = y * n
	}
	cards[cardIndex] = card
	cardIndex++

	// Step 2: Lines through [0, b, 1] for b = 0 to 6, slopes 0-6
	for b := 0; b < n; b++ {
		for m := 0; m < n; m++ {
			card := make([]int, symbolsPerCard)
			card[0] = n*n + m // [1, m, 0], 49-55
			for x := 0; x < n; x++ {
				y := (m*x + b) % n
				card[x+1] = x + y*n
			}
			cards[cardIndex] = card
			cardIndex++
			if cardIndex >= totalCards {
				break // Prevent overrun
			}
		}
		if cardIndex >= totalCards {
			break
		}
	}

	// Step 3: Line at infinity (only if room)
	if cardIndex < totalCards {
		card := make([]int, symbolsPerCard)
		for i := 0; i < symbolsPerCard; i++ {
			card[i] = n*n + i // 49-56
		}
		cards[cardIndex] = card
		cardIndex++
	}

	if cardIndex != totalCards {
		log.Printf("Generated %d cards, expected %d", cardIndex, totalCards)
	}

	// Map to stickers and check duplicates
	deck := make([][]string, totalCards)
	for i := 0; i < totalCards; i++ {
		deck[i] = make([]string, symbolsPerCard)
		seen := make(map[int]bool)
		for j, sym := range cards[i] {
			if sym >= 57 {
				log.Printf("Out of range symbol %d on card %d", sym, i)
				sym = sym % 57 // Safety clamp
			}
			if seen[sym] {
				log.Printf("Duplicate symbol index %d (%s) on card %d", sym, stickers[sym], i)
			}
			seen[sym] = true
			deck[i][j] = stickers[sym]
		}
	}

	// Verify exactly one match
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Bad Request", 400)
		return
	}

	files := r.MultipartForm.File["stickers"]
	if len(files) < 57 {
		http.Error(w, "Need at least 57 stickers", 400)
		return
	}

	os.RemoveAll(uploadDir)
	os.Mkdir(uploadDir, 0755)

	for i, fileHeader := range files[:57] {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error reading file", 500)
			return
		}
		defer file.Close()

		filename := filepath.Base(fileHeader.Filename)
		out, err := os.Create(filepath.Join(uploadDir, filename))
		if err != nil {
			http.Error(w, "Error saving file", 500)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			http.Error(w, "Error saving file", 500)
			return
		}
		stickers[i] = filename
	}

	deck := generateSpotIt()
	for i := 0; i < 56; i++ {
		for j := i + 1; j < 57; j++ {
			matches := 0
			for _, s1 := range deck[i] {
				for _, s2 := range deck[j] {
					if s1 == s2 {
						matches++
					}
				}
			}
			if matches != 1 {
				log.Printf("Cards %d and %d have %d matches", i, j, matches)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deck":     deck,
		"stickers": stickers,
	})
}

func main() {
	fs := http.FileServer(http.Dir("uploads"))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".png":
			w.Header().Set("Content-Type", "image/png")
		case ".jpg", ".jpeg":
			w.Header().Set("Content-Type", "image/jpeg")
		default:
			w.Header().Set("Content-Type", "application/octet-stream")
		}
		fs.ServeHTTP(w, r)
	})))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
