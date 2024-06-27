package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"testing"

	publish "github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
)

func setupService() func() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.WithGroup("PlayerData module")
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file", err)
		os.Exit(1)
	}

	connectionString := os.Getenv("MQ_CONNECTION_STRING")
	pubSub := publish.NewPubService(
		connectionString,
		log,
	)
	// init the exchange
	pubSub.Init()

	return func() {
		pubSub.Close()
	}
}

func TestInit(t *testing.T) {
	teardown := setupService()
	defer teardown()
	// var wg sync.WaitGroup
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
}
