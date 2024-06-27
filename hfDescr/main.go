package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/Bitstarz-eng/event-processing-challenge/hfDescr/messages"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.WithGroup("Human Readable module")
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file", err)
		os.Exit(1)
	}
	connectionString := os.Getenv("MQ_CONNECTION_STRING")
	pubSub := rabbitservice.NewPubService(connectionString, log)

	msgCh, err := pubSub.Subscribe()
	casinoEventCh := make(chan casino.Event)
	go func() {
		for msg := range msgCh {
			var casinoEvent casino.Event
			err := json.Unmarshal(msg.Body, &casinoEvent)
			if err != nil {
				log.Error("Error unmarshalling message", err)
				msg.Nack(false, false)
				continue
			}
			casinoEventCh <- casinoEvent
			// log.Info("Received message", "message", msg)
		}
	}()
	defer close(casinoEventCh)
	for event := range casinoEventCh {
		if err != nil {
			log.Error("Error getting player from database", err)
			continue
		}
		logMessage := messages.EventToMessage(event)

		log.Info(logMessage)
	}
}
