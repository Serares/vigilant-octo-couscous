package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/Bitstarz-eng/event-processing-challenge/playerData/db"

	_ "github.com/lib/pq"
)

type Players struct {
	db           *db.Queries
	dbConnection *sql.DB
}

func NewPsqlPlayersRepository() (*Players, error) {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		os.Getenv("PSQL_USER"),
		os.Getenv("PSQL_DBNAME"),
		os.Getenv("PSQL_PASSWORD"),
		os.Getenv("PSQL_HOST"),
	)
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error trying to create the db connection string %w", err)
	}
	dbQueris := db.New(dbConn)

	return &Players{
		db:           dbQueris,
		dbConnection: dbConn,
	}, nil
}

func (p *Players) GetPlayerById(ctx context.Context, id int64) (db.Player, error) {
	pl, err := p.db.GetPlayerById(ctx, id)
	if errors.Is(sql.ErrNoRows, err) {
		return db.Player{}, fmt.Errorf("player with id %d not found %w", id, err)
	}
	if err != nil {
		return db.Player{}, fmt.Errorf("error trying to get player with id %d %w", id, err)
	}
	return pl, nil
}

func (p *Players) CloseConnection() error {
	return p.dbConnection.Close()
}
