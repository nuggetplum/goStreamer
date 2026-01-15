package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func streamHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Open the file on disk
	// We use os.Open because it gives us a pointer to the file,
	// it doesn't read the whole thing yet.
	file, err := os.Open("songs/test_song.mp3")
	if err != nil {
		http.Error(w, "File not found.", http.StatusNotFound)
		return
	}
	// Always close file descriptors to avoid leaks!
	defer file.Close()

	// 2. Set the header so the browser knows it's audio
	w.Header().Set("Content-Type", "audio/mpeg")

	// 3. The Core Logic: io.Copy
	// This function reads from 'file' in small chunks (usually 32kb)
	// and writes immediately to 'w' (the response).
	// It repeats this until the file is done.
	// RAM usage stays low, regardless of file size.
	log.Println("Client connected, starting stream...")
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Error while streaming:", err)
	}
}

func main() {
	http.HandleFunc("/stream", streamHandler)

	log.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
