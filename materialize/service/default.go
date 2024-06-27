package service

import (
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/materialize/stats"
	"github.com/Bitstarz-eng/event-processing-challenge/pubsub/rabbitservice"
)

type MaterializeService struct {
	Logger *slog.Logger
	PubSub *rabbitservice.PubSubService
	Stats  *stats.MessageStats
	Lock   sync.Locker
}

func NewMaterializeService(
	logger *slog.Logger,
	pubsub *rabbitservice.PubSubService,
	messageStats *stats.MessageStats,
	locker sync.Locker,
) *MaterializeService {
	return &MaterializeService{
		Logger: logger.WithGroup("Materialize Service"),
		PubSub: pubsub,
		Stats:  messageStats,
		Lock:   locker,
	}
}

func (s *MaterializeService) InitSubscribe() (*stats.MessageStats, error) {
	msgCh, err := s.PubSub.Subscribe()
	if err != nil {
		s.Logger.Error("Error subscribing to queue", err)
	}

	go func() {
		for msg := range msgCh {
			var evnt casino.Event
			json.Unmarshal(msg.Body, &evnt)
			s.handleMessage(evnt)
			s.ingestMessage(evnt)
			s.updateTimes()
		}
	}()

	return s.Stats, nil
}

func (s *MaterializeService) ingestMessage(msg casino.Event) error {
	switch msg.Type {
	case "bet":
		s.Lock.Lock()
		defer s.Lock.Unlock()
		if msg.HasWon {
			s.Stats.PlayerWins[int64(msg.PlayerID)] += int64(msg.AmountEUR)
		} else {
			s.Stats.PlayerBets[int64(msg.PlayerID)] += int64(msg.AmountEUR)
		}

	case "deposit":
		s.Lock.Lock()
		defer s.Lock.Unlock()
		s.Stats.PlayerDeposits[int64(msg.PlayerID)] += int64(msg.AmountEUR)
	default:
		return nil
	}
	return nil
}

func (s *MaterializeService) handleMessage(msg casino.Event) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	s.Stats.TotalMessages++
	now := time.Now()
	s.Stats.EventTimes = append(s.Stats.EventTimes, now)
	// Remove events older than 1 minute
	cutoff := now.Add(-1 * time.Minute)
	i := 0
	for _, t := range s.Stats.EventTimes {
		if t.After(cutoff) {
			break
		}
		i++
	}
	s.Stats.EventTimes = s.Stats.EventTimes[i:]
	// Process the message here
	s.Logger.Info("Got message with id: ", "id", msg.ID)
}

// the methods are locking
func (s *MaterializeService) updateTimes() {
	s.Stats.EventsPerMinute = s.Stats.GetEventsPerMinute(s.Lock)
	s.Stats.EventsPerSecond = s.Stats.GetEventsPerSecond(s.Lock)
}

func (s *MaterializeService) GetStats() *stats.MessageStats {
	return s.Stats
}
