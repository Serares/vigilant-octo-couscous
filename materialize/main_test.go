package main

import (
	"fmt"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/Bitstarz-eng/event-processing-challenge/materialize/mux"
)

func setupApi(t *testing.T) (string, func()) {
	t.Helper()

	m := mux.SetupMux()

	ts := httptest.NewServer(m)

	return ts.URL, func() {
		ts.Close()
	}
}

func TestApi(t *testing.T) {
	url, teardown := setupApi(t)
	fmt.Println("The server url", url)
	// var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	defer teardown()
}
