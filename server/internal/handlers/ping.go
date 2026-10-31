package handlers

import "net/http"

func Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	w.Write([]byte("pong"))
}
