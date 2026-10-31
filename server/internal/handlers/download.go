package handlers

import (
	"crypto/rand"
	"net/http"
	"strconv"
)

const (
	chunkSize   = 1 << 20   // 1 MB (1 shifted left 20 bits = 1,048,576)
	maxSize     = 200 << 20 // never send more than 200 MB per request
	defaultSize = 25 << 20  // 25 MB if the client doesn't specify
)

// One random chunk, built once when the server starts, and reused for every request.
var chunk = makeChunk()

func makeChunk() []byte {
	b := make([]byte, chunkSize)
	rand.Read(b)
	return b
}

// Download streams `size` bytes of random data. Usage: /api/download?size=10000000
func Download(w http.ResponseWriter, r *http.Request) {
	size := defaultSize
	if v := r.URL.Query().Get("size"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			http.Error(w, "invalid size", http.StatusBadRequest)
			return
		}
		size = n
	}
	if size > maxSize {
		size = maxSize
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(size))
	w.Header().Set("Cache-Control", "no-store")

	remaining := size
	for remaining > 0 {
		n := chunkSize
		if remaining < n {
			n = remaining
		}
		if _, err := w.Write(chunk[:n]); err != nil {
			return // client disconnected, stop sending
		}
		remaining -= n
	}
}
