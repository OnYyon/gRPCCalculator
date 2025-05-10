package sqlite

import (
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(
	storagePath string,
	migrPath string,
) (*Storage, error) {
	conn, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, err
	}
	migrator, err := NewMigrator(conn, migrPath, storagePath)
	if err != nil {
		return nil, err
	}
	err = migrator.ApplyMigrations()
	if err != nil {
		return nil, err
	}
	return &Storage{db: conn}, nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Storage) SaveNewUser(
	ctx context.Context, login string,
	passHash []byte,
) error {
	// TODO:
	return nil
}

func (s *Storage) SaveExpression(
	ctx context.Context,
	expression string,
	expressionID string,
) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO Expressions (user_id, expression, expressionID, status) 
		VALUES(0, ?, ?, "pending")
	`, expression, expressionID)
	return err
}

func (s *Storage) GetExpressionByID(
	ctx context.Context,
	expressionID string,
) (float64, error) {
	// TODO:
	return 0, nil
}

// TODO:
