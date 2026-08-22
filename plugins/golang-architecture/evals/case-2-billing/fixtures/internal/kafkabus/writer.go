package kafkabus

import "fmt"

// EventPublisher describes anything capable of publishing a domain event
// onto the message bus.
type EventPublisher interface {
	Publish(topic string, payload []byte) error
}

type Writer struct {
	brokers []string
}

func NewWriter(brokers []string) *Writer {
	return &Writer{brokers: brokers}
}

func (w *Writer) Publish(topic string, payload []byte) error {
	fmt.Printf("publishing to %s on %v\n", topic, w.brokers)
	return nil
}
