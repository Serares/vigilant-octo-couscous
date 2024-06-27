package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Bitstarz-eng/event-processing-challenge/materialize/service"
	"github.com/Bitstarz-eng/event-processing-challenge/materialize/stats"
	"github.com/Bitstarz-eng/event-processing-challenge/materialize/utils"
)

type MaterializeHandler struct {
	Logger  *slog.Logger
	Service *service.MaterializeService
}

func NewMaterializeHandler(
	logger *slog.Logger,
	mService *service.MaterializeService,
) *MaterializeHandler {
	return &MaterializeHandler{
		Logger:  logger,
		Service: mService,
	}
}

func getHighestPlayer(playerStats map[int64]int64) stats.TopPlayer {
	var maxKey, maxValue int64
	for key, value := range playerStats {
		maxKey = key
		maxValue = value
		break
	}

	for key, value := range playerStats {
		if value > maxValue {
			maxKey = key
			maxValue = value
		}
	}

	return stats.TopPlayer{
		ID:    int(maxKey),
		Count: int(maxValue),
	}

}

func (h *MaterializeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Logger.Info("MaterializeHandler.ServeHTTP", slog.String("method", r.Method), slog.String("url", r.URL.String()))
	switch r.Method {
	case http.MethodGet:
		var message stats.EventStats
		stats := h.Service.GetStats()

		message.EventsTotal = stats.TotalMessages
		message.EventsPerMinute = stats.EventsPerMinute
		message.EventsPerSecondMovingAverage = stats.EventsPerSecond
		message.TopPlayerBets = getHighestPlayer(stats.PlayerBets)
		message.TopPlayerDeposits = getHighestPlayer(stats.PlayerDeposits)
		message.TopPlayerWins = getHighestPlayer(stats.PlayerWins)
		utils.ReplySuccess(w, r, http.StatusAccepted, message)
	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
	}
}
