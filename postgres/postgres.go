package postgres

import (
	"context"
	"fmt"
	"log"
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/trenchesdeveloper/tweeter/config"
)

type DB struct {
	Pool *pgxpool.Pool
	config config.Config
}

func New(ctx context.Context, config *config.Config) *DB {
	dbConf, err := pgxpool.ParseConfig(config.Database.Url)
	if err != nil {
		log.Fatalf("can't parse postgres config: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, dbConf)

	if err != nil {
		log.Fatalf("can't connect to postgres: %v", err)
	}

	db := &DB{Pool: pool, config: *config}

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

func (db *DB) Migrate() error {
	_, b, _, _ := runtime.Caller(0)
	migrationPath := fmt.Sprintf("file:///%s/migrations", path.Dir(b))
	log.Println("db url``	", db.config.Database.Url)
	m, err := migrate.New(migrationPath, db.config.Database.Url)

	if err != nil {
		return fmt.Errorf("can't create migration: %v", err)
	}

	if err := m.Up(); err != nil  && err != migrate.ErrNoChange {
		return fmt.Errorf("can't migrate: %v", err)
	}

	log.Println("postgres migrated")

	return nil
}
