package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type Expression struct {
	ID         string
	Expression string
	Status     string
	Result     string
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
	expressionID string,
	expression string,
	user_id string,
) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO Expressions (user_id, expression, expressionID, status) 
		VALUES(?, ?, ?, ?)
	`, user_id, expression, expressionID, "processing")
	return err
}

func (s *Storage) UpdateExpression(
	ctx context.Context,
	expressionID string,
	result float64,
) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE Expressions
		SET result = ?, status = ? WHERE expressionID = ?`,
		result, "completed", expressionID)
	return err
}

func (s *Storage) GetExpressionByID(
	ctx context.Context,
	expressionID string,
) (map[string]string, error) {
	var (
		id         string
		expression string
		status     string
		result     sql.NullString
	)
	err := s.db.QueryRowContext(ctx, `
		SELECT expressionID, expression, status, result FROM Expressions WHERE expressionID = ?`, expressionID,
	).Scan(&id, &expression, &status, &result)
	o := make(map[string]string)
	o["id"] = id
	o["expression"] = expression
	o["status"] = status
	if result.Valid {
		o["result"] = result.String
	} else {
		o["result"] = ""
	}
	return o, err
}

func (s *Storage) GetExpressionList(
	ctx context.Context,
	userid string,
) ([]Expression, error) {
	rows, err := s.db.Query("SELECT expressionID, expression, status, result FROM Expressions WHERE user_id = ?", userid)
	if err != nil {
		return nil, fmt.Errorf("failed to query expressions: %v", err)
	}
	defer rows.Close()

	var expressions []Expression
	var res sql.NullString
	for rows.Next() {
		var expr Expression
		err := rows.Scan(&expr.ID, &expr.Expression, &expr.Status, &res)
		if !res.Valid {
			expr.Result = ""
		} else {
			expr.Result = res.String
		}
		if err != nil {
			return nil, fmt.Errorf("failed to scan expression: %v", err)
		}
		expressions = append(expressions, expr)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error after iterating rows: %v", err)
	}

	return expressions, nil
}

func (s *Storage) RegisterUser(
	ctx context.Context,
	login string,
	passwordHash []byte,
) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO Users (username, password_hash)
		VALUES (?, ?)`, login, passwordHash,
	)
	return err
}

func (s *Storage) GetUser(
	ctx context.Context,
	login string,
) ([]byte, error) {
	var row []byte
	err := s.db.QueryRowContext(ctx, `
		SELECT password_hash FROM Users WHERE username = ?`, login,
	).Scan(&row)
	return row, err
}
