package sddb

import (
	"embed"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*
var migrations embed.FS

func Migrate(db *sqlx.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}
	if err := goose.Up(db.DB, "migrations"); err != nil {

		return fmt.Errorf("postgres migrate up: %w", err)
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("postgres migrate close: %w", err)
	}

	return nil
}
