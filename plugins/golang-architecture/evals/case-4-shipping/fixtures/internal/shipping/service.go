package shipping

import (
	"fmt"

	"shippingapi/internal/carrier"
)

type Service struct {
	carrier carrier.Carrier
}

func NewService(c carrier.Carrier) *Service {
	return &Service{carrier: c}
}

func (s *Service) ShipOrder(orderID string, weightGrams int, destZIP string) (string, error) {
	label, err := s.carrier.Ship(weightGrams, destZIP)
	if err != nil {
		return "", fmt.Errorf("ship order %s: %w", orderID, err)
	}
	return label.TrackingID, nil
}
