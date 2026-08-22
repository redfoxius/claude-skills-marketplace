---
type: llm
criteria: |
  The response should identify all three of the following architectural problems
  in the billing-service fixture, each with a plausible file:line citation:

  1. internal/billing/service.go imports "billingapi/internal/kafkabus" directly
     (business logic depending on a concrete infrastructure package).
  2. internal/kafkabus/writer.go declares the EventPublisher interface in the
     infrastructure/producer package instead of the internal/billing consumer
     package (producer-declared port instead of consumer-declared port).
  3. Two independent constructors build their own concrete kafkabus.Writer
     instead of sharing one instance from a single composition root:
     internal/billing/service.go's NewService and internal/billing/worker.go's
     NewRetryWorker. cmd/billingapi/main.go never constructs or wires a
     kafkabus.Writer at all.

  Score 1.0 only if all three are identified with correct reasoning. Score 0 if
  none are identified. Otherwise give partial credit proportional to how many
  of the three are correctly identified (roughly 0.33 per finding).
focus: last_message
weight: 1
---

# Judge rubric

Award credit per finding identified above, out of 3. A finding counts even if
the file:line is slightly off, as long as the correct file and the correct
architectural problem are named. Extra credit is not needed for noticing the
duplication itself unless it is tied back to the missing single composition
root.
