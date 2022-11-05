package db

import (
	"context"
	"log"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var connPool *pgxpool.Pool
var dbUrl *string

func SetDatabaseUrl(url string) {
	dbUrl = &url
}

func Connect() (*pgxpool.Pool, error) {
	if connPool == nil {
		if dbUrl == nil {
			log.Fatal("database url not set")
		}
		dbConfig, err := pgxpool.ParseConfig(*dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		dbConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			pgxuuid.Register(conn.TypeMap())
			return nil
		}
		connPool, err = pgxpool.NewWithConfig(context.Background(), dbConfig)
		if err != nil {
			log.Fatalf("db connection error: %s", err)
		}
	}
	return connPool, nil
}

func Close() {
	connPool.Close()
	connPool = nil
}
