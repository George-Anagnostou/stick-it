package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
)

const uploadDir = "uploads"

var stickers []string // Dynamic size, no fixed allocation

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
	if len(files) < 7 { // Minimum for n=2
		http.Error(w, "Need at least 7 stickers", 400)
		return
	}

	os.RemoveAll(uploadDir)
	os.Mkdir(uploadDir, 0755)

	// Allocate stickers based on uploaded files
	stickers = make([]string, len(files))
	for i, fileHeader := range files { // Use full slice, not fixed 57
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

	deck := GenerateSpotItGeneric(stickers)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"deck":     deck,
		"stickers": stickers,
	})
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	if len(stickers) < 7 {
		http.Error(w, "No deck generated yet", 400)
		return
	}

	deck := GenerateSpotItGeneric(stickers)

	pdf := gofpdf.New("P", "mm", "Letter", "") // 8.5' x 11"
	pdf.SetMargins(10, 10, 10)
	pdf.AddPage()

	cardWidth := 50.0
	cardHeight := 50.0
	xOffset := 10.0
	yOffset := 10.0
	cardsPerRow := 4
	row := 0
	col := 0

	for _, card := range deck {
		if col >= cardsPerRow {
			col = 0
			row++
			if row >= cardsPerRow { // New page after 16 cards
				pdf.AddPage()
				row = 0
				yOffset = 10.0
			}
		}

		x := xOffset + float64(col)*cardWidth
		y := yOffset + float64(row)*cardHeight

		// Draw card border
		pdf.SetFillColor(255, 255, 255)
		pdf.Rect(x, y, cardWidth, cardHeight, "FD")

		// Draw symbols in a circular layout
		radius := 20.0
		centerX := x + cardWidth/2
		centerY := y + cardHeight/2
		symbolCount := len(card)
		for i, sticker := range card {
			angle := float64(i) / float64(symbolCount) * 2 * 3.14159
			sx := centerX + radius*float64(math.Cos(angle)) - 10 // Adjust for image size
			sy := centerY + radius*float64(math.Sin(angle)) - 10
			pdf.Image(filepath.Join(uploadDir, sticker), sx, sy, 20, 20, false, "", 0, "")
		}

		col++
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		http.Error(w, "Error generating PDF", 500)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=spotit_deck.pdf")
	w.Write(buf.Bytes())
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
	http.HandleFunc("/export", exportHandler)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
