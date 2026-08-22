---
type: llm
criteria: |
  The shipping-service fixture has clean dependency direction (business logic
  never imports a concrete adapter type) and a single composition root
  (cmd/shippingapi/main.go constructs the only concrete adapter, FedExClient,
  and injects it). Its one architectural problem is that the Carrier interface
  is declared in internal/carrier/carrier.go (the producer/infrastructure
  package) instead of in internal/shipping (the consumer/domain package) that
  actually uses it — inverting the Go-idiomatic rule that ports are declared by
  the consumer, not pre-declared by the producer. A secondary supporting signal
  is that Carrier is a fat interface (Ship, Cancel, RateQuote) shaped around
  everything FedExClient offers, even though shipping.Service only ever calls
  Ship — evidence the interface was authored from the producer's side rather
  than carved out to the consumer's actual need.

  Score 1.0 if the response correctly identifies the port-ownership-inversion
  problem (regardless of exact wording) with a plausible file:line citation
  pointing at internal/carrier/carrier.go and/or internal/shipping/service.go.

  Score 0 if the response says the code is well-architected / has no
  significant issues, or only raises unrelated nits (naming, error-wrapping
  style, missing tests) without ever flagging where the Carrier interface
  should live.
focus: last_message
weight: 1
---

# Judge rubric

This is the discriminating case: a generic senior-engineer review (no
Go-specific idiom awareness) is expected to often praise this code as
well-designed DI and miss the port-location problem, since "interface next to
its implementation" is conventional in many other architectural traditions.
Score strictly — a response that only says "this is clean, no issues found"
must score 0, even though the rest of its reasoning about DI/composition-root
being correct is itself accurate.
