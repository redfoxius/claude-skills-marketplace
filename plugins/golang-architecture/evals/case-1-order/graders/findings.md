---
type: llm
criteria: |
  The response should identify all three of the following architectural problems
  in the order-service fixture, each with a plausible file:line citation:

  1. internal/order/service.go imports "database/sql" and "github.com/lib/pq"
     directly and executes a raw SQL query inside a business-logic method
     (business logic depending on concrete infrastructure).
  2. internal/postgres/repository.go declares the OrderRepository interface in
     the infrastructure/producer package instead of the internal/order consumer
     package (producer-declared port instead of consumer-declared port).
  3. internal/order/service.go's NewService constructs a concrete
     *notify.SlackClient directly instead of receiving it as an injected
     dependency from a composition root (cmd/orderapi/main.go).

  Score 1.0 only if all three are identified with correct reasoning. Score 0 if
  none are identified. Otherwise give partial credit proportional to how many
  of the three are correctly identified (roughly 0.33 per finding).
focus: last_message
weight: 1
---

# Judge rubric

Award credit per finding identified above, out of 3. A finding counts even if
the file:line is slightly off, as long as the correct file and the correct
architectural problem are named. Do not award credit for unrelated nits
(naming, error wrapping, missing tests).
