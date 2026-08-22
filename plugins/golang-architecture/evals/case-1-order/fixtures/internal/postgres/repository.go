package postgres

import "orderapi/internal/order"

// OrderRepository defines the persistence contract for orders.
type OrderRepository interface {
	Save(o order.Order) error
	FindByID(id string) (order.Order, error)
}

type PGOrderRepository struct {
	dsn string
}

func NewPGOrderRepository(dsn string) *PGOrderRepository {
	return &PGOrderRepository{dsn: dsn}
}

func (r *PGOrderRepository) Save(o order.Order) error {
	return nil
}

func (r *PGOrderRepository) FindByID(id string) (order.Order, error) {
	return order.Order{}, nil
}
