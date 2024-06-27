package main

import (
	"log/slog"
	"os"

	publish "github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
)

// it's not actually used
func main() {
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
}
