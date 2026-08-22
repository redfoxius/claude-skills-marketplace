package store

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotFound = errors.New("store: not found")

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) FindActive(ctx context.Context, userID string) (string, error) {
	row := s.db.QueryRowContext(ctx, "SELECT plan FROM subscriptions WHERE user_id = $1 AND active", userID)

	var plan string
	if err := row.Scan(&plan); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	return plan, nil
}
