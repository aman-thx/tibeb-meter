package main

import (
	"log"
	"net/http"

	"tibeb-meter/server/internal/config"
	"tibeb-meter/server/internal/handlers"
)

func main() {
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/ping", handlers.Ping)
	mux.HandleFunc("GET /api/download", handlers.Download)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	log.Printf("tibeb-meter server listening on :%s", cfg.Port)
	log.Fatal(srv.ListenAndServe())
}
