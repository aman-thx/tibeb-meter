package handlers

import (
	"io"
	"net/http"
	"strconv"
)

const maxUpload = 100 << 20 // refuse anything over 100 MB

// Upload reads the whole request body and discards it.
// The client measures how long the upload took.
func Upload(w http.ResponseWriter, r *http.Request) {
	// MaxBytesReader stops reading after maxUpload bytes and returns an error,
	// so nobody can push gigabytes at your server.
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)

	n, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		http.Error(w, "upload too large or interrupted", http.StatusRequestEntityTooLarge)
		return
	}

	w.Header().Set("Cache-Control", "no-store")
	w.Write([]byte(strconv.FormatInt(n, 10)))
}
