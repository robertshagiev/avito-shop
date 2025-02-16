package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log/slog"
)

func Connect(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		slog.Error("sql.Open", "error", err.Error())
		return nil, err
	}

	return db, nil
}
