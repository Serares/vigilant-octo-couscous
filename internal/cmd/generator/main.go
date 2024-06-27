package main

import (
	"errors"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/exchanger/exchange"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/generator"
	publish "github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
)

func convertAmount(
	event *casino.Event,
	cache *exchange.AmountsStore,
	client *exchange.ExchangeClient,
	cacheFileName string,
	log *slog.Logger,
	l sync.Locker,
) (float64, error) {
	l.Lock()
	defer l.Unlock()
	eurAmount, err := cache.GetOne(event.Currency, float64(event.Amount))
	if err != nil && !errors.Is(err, exchange.ErrNotFound) {
		log.Error("Error trying to get the exchange rate from cache", err)
	}
	if errors.Is(err, exchange.ErrNotFound) {
		// get the exchange from api
		resp, err := client.GetConvertAmount(
			exchange.ExchangeRequest{
				From:   event.Currency,
				Amount: float64(event.Amount),
				To:     "EUR",
			},
		)
		if err != nil {
			log.Error("Error trying to get the exchange rate from api", err)
		}
		eurAmount = float64(resp.Result)
	}
	// event.AmountEUR = int(eurAmount)

	cache.Add(exchange.CurrencyToEuro{
		From:       event.Currency,
		Amount:     float64(event.Amount),
		AmountEuro: eurAmount,
	})
	defer cache.Save(cacheFileName)
	if err != nil {
		log.Error("Error saving cache", err)
	}

	return eurAmount, nil
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	err := godotenv.Load(".env")
	if err != nil {
		log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		log.Error("Error loading .env file")
	}
	connectionString := os.Getenv("MQ_CONNECTION_STRING")
	cacheFileName := os.Getenv("CACHE_FILE_NAME")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.WithGroup("Generator")
	lock := &sync.Mutex{}
	pubSubServ := publish.NewPubService(
		connectionString,
		log,
	)

	cache := &exchange.AmountsStore{}
	err = cache.Get(cacheFileName)
	if err != nil {
		log.Error("Error getting initial cache", err)
	}

	exClient := exchange.NewExchangeClient(log)
	defer cache.Save(cacheFileName)

	eventCh := generator.Generate(ctx)

	for event := range eventCh {
		if event.Currency == "EUR" {
			// store the amount euro
			event.AmountEUR = event.Amount
		} else {
			convertedAmount, err := convertAmount(
				&event,
				cache,
				exClient,
				cacheFileName,
				log,
				lock,
			)
			if err != nil {
				log.Error("error trying to convert the amount", err)
				continue
			}
			event.AmountEUR = int(convertedAmount)
		}
		pubSubServ.Publish(event)
	}

	log.Info("finished")
}
