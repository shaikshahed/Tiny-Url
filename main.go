package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct {
	ID           string    `json:"id"`
	OriginalUrl  string    `json:"original_url"`
	ShortUrl     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

var urlDB = make(map[string]URL)

func generateShortURL(OriginalUrl string) string {
	hasher := md5.New()               // Create a new MD5 hasher instance
	hasher.Write([]byte(OriginalUrl)) // Write the OriginalUrl as bytes to the hasher
	fmt.Println("hasherr:", hasher)
	data := hasher.Sum(nil) // Get the MD5 hash (as a byte slice)
	fmt.Println("hasher data:", data)
	hash := hex.EncodeToString(data) // Convert the hash to a hexadecimal string
	fmt.Println("hasher data string:", hash)
	fmt.Println("hasher data final string:", hash[:8])
	return hash[:8] // Return the first 8 characters of the hash as the short URL
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalUrl:  originalURL,
		ShortUrl:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) {
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "get")
}

func ShortURLHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadGateway)
		return
	}

	shortURL_ := createURL(data.URL)
	response := struct {
		ShortUrl string `json:"short_url"`
	}{ShortUrl: shortURL_}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalUrl, http.StatusFound)
}

func main() {
	// fmt.Println("Starting URL shortner")
	// OriginalUrl := "https://github.com"
	// generateShortURL(OriginalUrl)

	// Register the handler function to handle all requests to the root URL ("/")
	http.HandleFunc("/", handler)
	http.HandleFunc("/shorten", ShortURLHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)
	// Start the HTTP server on port 3000
	fmt.Println("Starting the server on port 3000...")
	http.ListenAndServe(":3000", nil)
}
