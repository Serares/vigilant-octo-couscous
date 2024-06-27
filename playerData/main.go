package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/playerData/repo"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	log.WithGroup("PlayerData module")
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file", err)
		os.Exit(1)
	}
	ctx := context.Background()

	// init db
	playerRepo, err := repo.NewPsqlPlayersRepository()
	if err != nil {
		log.Error("Error initializing player repository", err)
		os.Exit(1)
	}

	connectionString := os.Getenv("MQ_CONNECTION_STRING")
	pubSub := rabbitservice.NewPubService(
		connectionString,
		log,
	)

	msgCh, err := pubSub.Subscribe()
	if err != nil {
		log.Error("Error subscribing to queue", err)
		os.Exit(1)
	}

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
	defer playerRepo.CloseConnection()

	for event := range casinoEventCh {
		dbPlayer, err := playerRepo.GetPlayerById(ctx, int64(event.PlayerID))
		if err != nil {
			log.Error("Error getting player from database", err)
			continue
		}
		event.Player.Email = dbPlayer.Email
		event.Player.LastSignedInAt = dbPlayer.LastSignedInAt.Time
		log.Info("event with player info", "event: ", event)
	}
}
