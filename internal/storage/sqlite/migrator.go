package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	"github.com/golang-migrate/migrate/v4"

	// NOTE: драйвера
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator struct {
	migrator *migrate.Migrate
}

func NewMigrator(db *sql.DB, cfg *config.Config) (*Migrator, error) {
	m, err := migrate.New("file://"+cfg.Database.MigrationsPath,
		"sqlite3://"+cfg.Database.DBPath)
	if err != nil {
		return nil, err
	}
	return &Migrator{m}, nil
}

func (m *Migrator) ApplyMigrations() error {
	if err := m.migrator.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

func (m *Migrator) Down() error {
	if err := m.migrator.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

func (m *Migrator) Close() error {
	sourceErr, databaseErr := m.migrator.Close()
	if sourceErr != nil || databaseErr != nil {
		return fmt.Errorf("source error: %v, database error: %v", sourceErr, databaseErr)
	}
	return nil
}
