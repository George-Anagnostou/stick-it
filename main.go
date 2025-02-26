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

	deck := GenerateSpotItGeneric(stickers)

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
