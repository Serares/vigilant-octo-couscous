package messages

import "github.com/Bitstarz-eng/event-processing-challenge/internal/casino"

func EventToMessage(message casino.Event) string {
	log := CreateMessage(message)
	if log != nil {
		return log.ToString()
	}

	return "Unknown event type"
}
