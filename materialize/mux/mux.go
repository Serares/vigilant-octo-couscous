package mux

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/Bitstarz-eng/event-processing-challenge/materialize/handler"
	"github.com/Bitstarz-eng/event-processing-challenge/materialize/service"
	"github.com/Bitstarz-eng/event-processing-challenge/materialize/stats"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
)

func SetupMux() *http.ServeMux {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
	}
	m := http.NewServeMux()
	CONNECTION_STRING := os.Getenv("MQ_CONNECTION_STRING")
	mu := &sync.Mutex{}

	messageStats := stats.MessageStats{
		PlayerBets:     make(map[int64]int64),
		PlayerWins:     make(map[int64]int64),
		PlayerDeposits: make(map[int64]int64),
	}

	pubSub := rabbitservice.NewPubService(
		CONNECTION_STRING,
		log,
	)

	ms := service.NewMaterializeService(
		log,
		pubSub,
		&messageStats,
		mu,
	)

	_, err = ms.InitSubscribe()
	if err != nil {
		log.Error("Error subscribing to queue", err)
	}

	mh := handler.NewMaterializeHandler(
		log,
		ms,
	)

	m.Handle("/materialize", mh)

	return m
}
