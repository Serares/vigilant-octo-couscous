package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/materialize/mux"
)

func main() {
	m := mux.SetupMux()

	server := &http.Server{
		Addr:         "localhost:8080",
		Handler:      m,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	slog.Info("Starting server on :8080")
	server.ListenAndServe()
}
