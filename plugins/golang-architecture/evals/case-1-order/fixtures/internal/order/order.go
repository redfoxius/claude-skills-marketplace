package order

import "time"

type Order struct {
	ID         string
	CustomerID string
	Total      int64
	Status     string
	CreatedAt  time.Time
}
