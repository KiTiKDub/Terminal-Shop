package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

type downloadEntry struct {
	filePath  string
	expiresAt time.Time
}

// sync is a package that allows you to safetly use things between sepearate go routines
// aka, our ssh and http routines. We use the map to add our download struct
var tokenStore sync.Map

func createDownloadToken(filepath string) string {

	//initalize a byte slice
	b := make([]byte, 16)
	//cryptographically generate bytes into b
	_, _ = rand.Read(b)
	//encode to string for url purposes
	token := hex.EncodeToString(b)

	//store in sync map that safely allows use for the SSH and HTTP server
	tokenStore.Store(token, downloadEntry{
		filePath:  filepath,
		expiresAt: time.Now().Add(5 * time.Minute),
	})

	return token
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Path[len("/dl/"):]
	val, ok := tokenStore.LoadAndDelete(token)
	if !ok {
		http.Error(w, "Invalid or expired download link", http.StatusNotFound)
		return
	}

	entry := val.(downloadEntry)
	if time.Now().After(entry.expiresAt) {
		http.Error(w, "Entry has expired", http.StatusGone)
		return
	}

	//prepare thte response to send the right file back.
	filename := filepath.Base(entry.filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename=%s`, filename))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, entry.filePath)
}
