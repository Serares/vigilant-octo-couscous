package repo

import (
	"context"

	"github.com/Bitstarz-eng/event-processing-challenge/playerData/db"
)

type IPlayersRepo interface {
	GetPlayerById(ctx context.Context, id int64) (db.Player, error)
}
