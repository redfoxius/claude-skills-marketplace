package billing

import (
	"encoding/json"

	"billingapi/internal/kafkabus"
)

// RetryWorker republishes invoices that failed to send on the first attempt.
type RetryWorker struct {
	publisher kafkabus.EventPublisher
	pending   []Invoice
}

func NewRetryWorker(brokers []string) *RetryWorker {
	return &RetryWorker{
		publisher: kafkabus.NewWriter(brokers),
		pending:   make([]Invoice, 0),
	}
}

func (w *RetryWorker) Enqueue(inv Invoice) {
	w.pending = append(w.pending, inv)
}

func (w *RetryWorker) FlushOnce() error {
	for _, inv := range w.pending {
		payload, err := json.Marshal(inv)
		if err != nil {
			return err
		}
		if err := w.publisher.Publish("invoice.retry", payload); err != nil {
			return err
		}
	}
	w.pending = w.pending[:0]
	return nil
}
