package postgres

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(dbDriver, dbUri string) *sqlx.DB {
	conn, err := sqlx.Open(dbDriver, dbUri)
	if err != nil {
		log.Fatalf("sqlx.Connect: %v", err)
	}

	if err := conn.PingContext(context.Background()); err != nil {
		log.Fatalf("sqlx.PingContext: %v", err)
	}

	return conn
}
