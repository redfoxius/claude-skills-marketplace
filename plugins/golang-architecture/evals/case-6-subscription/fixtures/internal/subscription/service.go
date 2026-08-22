package subscription

import (
	"context"
	"errors"

	"subscriptionapi/internal/store"
)

type Store interface {
	FindActive(ctx context.Context, userID string) (string, error)
}

type Service struct {
	store Store
}

func NewService(s Store) *Service {
	return &Service{store: s}
}

func (s *Service) Activate(ctx context.Context, userID, plan string) error {
	_, err := s.store.FindActive(ctx, userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return err
	}
	return errors.New("subscription already active")
}
