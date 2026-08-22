package billing

import (
	"encoding/json"
	"fmt"

	"billingapi/internal/kafkabus"
)

type Service struct {
	publisher kafkabus.EventPublisher
}

func NewService(brokers []string) *Service {
	return &Service{
		publisher: kafkabus.NewWriter(brokers),
	}
}

func (s *Service) IssueInvoice(inv Invoice) error {
	inv.Status = "issued"

	payload, err := json.Marshal(inv)
	if err != nil {
		return err
	}

	if err := s.publisher.Publish("invoice.issued", payload); err != nil {
		return fmt.Errorf("publish invoice.issued: %w", err)
	}

	return nil
}
