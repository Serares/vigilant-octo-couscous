package stats

import (
	"sync"
	"time"
)

type TopPlayer struct {
	ID    int `json:"id"`
	Count int `json:"count"`
}

type EventStats struct {
	EventsTotal                  int       `json:"events_total"`
	EventsPerMinute              float64   `json:"events_per_minute"`
	EventsPerSecondMovingAverage float64   `json:"events_per_second_moving_average"`
	TopPlayerBets                TopPlayer `json:"top_player_bets"`
	TopPlayerWins                TopPlayer `json:"top_player_wins"`
	TopPlayerDeposits            TopPlayer `json:"top_player_deposits"`
}

type MessageStats struct {
	TotalMessages   int             `json:"total_messages"`
	EventTimes      []time.Time     `json:"event_times"`
	EventsPerSecond float64         `json:"events_per_second"`
	EventsPerMinute float64         `json:"events_per_minute"`
	PlayerBets      map[int64]int64 `json:"player_bets"`
	PlayerWins      map[int64]int64 `json:"player_wins"`
	PlayerDeposits  map[int64]int64 `json:"player_deposits"`
}

var stats MessageStats
var once sync.Once

func GetInstance() *MessageStats {
	once.Do(func() {
		stats = MessageStats{}
	})
	return &stats
}

func (s *MessageStats) GetEventsPerSecond(l sync.Locker) float64 {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)
	count := 0
	for _, t := range s.EventTimes {
		if t.After(cutoff) {
			count++
		}
	}

	return float64(count) / 60.0
}

func (s *MessageStats) GetEventsPerMinute(l sync.Locker) float64 {
	l.Lock()
	defer l.Unlock()

	now := time.Now()
	cutoff := now.Add(-1 * time.Minute)
	count := 0
	for _, t := range s.EventTimes {
		if t.After(cutoff) {
			count++
		}
	}

	return float64(count)
}
