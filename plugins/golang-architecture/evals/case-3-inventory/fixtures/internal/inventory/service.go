package inventory

import (
	"fmt"

	"gorm.io/gorm"

	"inventoryapi/internal/gormrepo"
)

var db *gorm.DB

func init() {
	db = gormrepo.Connect("file:inventory.db")
}

type Service struct {
	repo gormrepo.StockRepository
}

func NewService() *Service {
	return &Service{
		repo: gormrepo.New(db),
	}
}

func (s *Service) Reserve(sku string, qty int) error {
	current, err := s.repo.GetQuantity(sku)
	if err != nil {
		return err
	}

	if current < qty {
		return fmt.Errorf("insufficient stock for %s: have %d, want %d", sku, current, qty)
	}

	return s.repo.AdjustQuantity(sku, -qty)
}
