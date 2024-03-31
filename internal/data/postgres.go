package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	pgmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/jackc/pgx/v5"
	"log/slog"
)

type PGStore[T interface{}] struct {
	db  *sql.DB
	log slog.Logger
}

func NewPGStore[T interface{}](ctx context.Context, dbUrl string, log slog.Logger) (*PGStore[T], error) {
	db, err := sql.Open("pgx", dbUrl)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	store := &PGStore[T]{
		db:  db,
		log: log,
	}

	err = store.runMigrations()
	if err != nil {
		log.Error("failed to run migrations", "error", err)
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return store, nil
}

func (s *PGStore[T]) runMigrations() error {
	driver, err := pgmigrate.WithInstance(s.db, &pgmigrate.Config{})
	if err != nil {
		s.log.Error("failed to create migration driver", "error", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		s.log.Error("failed to create migration instance", "error", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Up()
	if err == nil || errors.Is(err, migrate.ErrNoChange) {
		s.log.Info("no new migrations to apply")
		return nil
	}

	s.log.Error("failed to apply migrations", "error", err)
	return fmt.Errorf("failed to apply migrations: %w", err)
}
