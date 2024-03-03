package postgres

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/trenchesdeveloper/tweeter/config"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, config config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(config.Database.Url)
	if err != nil {
		log.Fatalf("can't parse postgres config: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, dbConf)

	if err != nil {
		log.Fatalf("can't connect to postgres: %v", err)
	}

	db := &DB{Pool: pool}

	db.Ping(ctx)

	return db
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Ping(ctx context.Context) {
	if err := db.Pool.Ping(ctx); err != nil {
		log.Fatalf("can't ping postgres: %v", err)
	}

	log.Println("postgres connected")
}
