package sqlite

import (
	"database/sql"

	"github.com/OnYyon/gRPCCalculator/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

// type Storage struct {
// 	db *sql.DB
// }

func MustRunNewStorage(cfg *config.Config) {
	conn, err := sql.Open("sqlite3", cfg.Database.DBPath)
	if err != nil {
		panic(err)
	}
	migrator, err := NewMigrator(conn, cfg)
	if err != nil {
		panic(err)
	}
	err = migrator.ApplyMigrations()
	if err != nil {
		panic(err)
	}
}
