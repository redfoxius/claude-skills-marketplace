package order

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"orderapi/internal/notify"
	"orderapi/internal/postgres"
)

type Service struct {
	repo     postgres.OrderRepository
	notifier *notify.SlackClient
	db       *sql.DB
}

func NewService(repo postgres.OrderRepository, dsn string) (*Service, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	notifier := notify.NewSlackClient("https://hooks.slack.com/services/T000/B000/XXXX")

	return &Service{
		repo:     repo,
		notifier: notifier,
		db:       db,
	}, nil
}

func (s *Service) MarkShipped(orderID string) error {
	row := s.db.QueryRow("SELECT status FROM orders WHERE id = $1", orderID)
	var status string
	if err := row.Scan(&status); err != nil {
		return err
	}

	if status == "shipped" {
		return fmt.Errorf("order %s already shipped", orderID)
	}

	o, err := s.repo.FindByID(orderID)
	if err != nil {
		return err
	}
	o.Status = "shipped"

	if err := s.repo.Save(o); err != nil {
		return err
	}

	return s.notifier.Notify(fmt.Sprintf("order %s shipped", orderID))
}
