package messages

import (
	"fmt"

	"github.com/Bitstarz-eng/event-processing-challenge/hfDescr/utils"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

func CreateMessage(event casino.Event) Message {
	switch event.Type {
	case "game_start":
		return &MessageGameStart{
			CasinoEvent: event,
			GameName:    casino.Games[event.GameID].Title,
		}
	case "bet":
		return &MessageBet{
			CasinoEvent: event,
			GameName:    casino.Games[event.GameID].Title,
		}
	case "deposit":
		return &MessageDeposit{
			CasinoEvent: event,
		}
	default:
		return nil
	}
}

type Message interface {
	ToString() string
}

type MessageGameStart struct {
	CasinoEvent casino.Event
	GameName    string
}

// create a to string method
func (m *MessageGameStart) ToString() string {
	return fmt.Sprintf(
		"Player #%d started playing a game \"%s\" on %s.",
		m.CasinoEvent.PlayerID,
		m.GameName,
		utils.FormatDate(m.CasinoEvent.CreatedAt),
	)
}

type MessageBet struct {
	CasinoEvent casino.Event
	GameName    string
}

// create a to string method
func (m *MessageBet) ToString() string {
	return fmt.Sprintf(
		"Player #%d (john@example.com) placed a bet of %d %s (%d EUR) on a game \"%s\" on %s",
		m.CasinoEvent.PlayerID,
		m.CasinoEvent.Amount,
		m.CasinoEvent.Currency,
		m.CasinoEvent.AmountEUR,
		m.GameName,
		utils.FormatDate(m.CasinoEvent.CreatedAt),
	)
}

type MessageDeposit struct {
	CasinoEvent casino.Event
}

// create a to string method
// create a to string method
func (m *MessageDeposit) ToString() string {
	return fmt.Sprintf(
		"Player #%d made a deposit of %d %s on %s.",
		m.CasinoEvent.PlayerID,
		m.CasinoEvent.Amount,
		m.CasinoEvent.Currency,
		utils.FormatDate(m.CasinoEvent.CreatedAt),
	)
}
