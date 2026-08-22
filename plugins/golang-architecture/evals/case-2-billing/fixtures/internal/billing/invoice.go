package billing

import "time"

type Invoice struct {
	ID          string
	CustomerID  string
	AmountCents int64
	Status      string
	IssuedAt    time.Time
}
