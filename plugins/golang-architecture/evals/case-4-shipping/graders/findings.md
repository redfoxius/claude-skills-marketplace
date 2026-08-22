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
  pointing at internal/carrier/carrier.go and/or internal/shipping/service.go,
  AND does not claim dependency direction or composition-root placement are
  broken (they are not — cmd/shippingapi/main.go is the sole construction
  site for FedExClient, and shipping never imports a concrete infra type).

  Score 0 if the response says the code is well-architected / has no
  significant issues, only raises unrelated nits (naming, error-wrapping
  style, missing tests) without ever flagging where the Carrier interface
  should live, or invents a dependency-direction/composition-root problem
  that isn't actually present.
focus: last_message
weight: 1
---

# Judge rubric

This case is a precision / no-false-positive sanity check, not a
discriminating trap: validated on 2026-08-23 that a general senior-engineer
review (no golang-architecture skill) reliably catches the port-location
problem here too, both when reviewed alongside cases 1-3 and in full
isolation — "interface declared next to its implementation" turned out to be
common-enough Go knowledge on its own. Use this case to confirm the skill (a)
still finds the one real issue and (b) doesn't over-flag the parts of the
code that are genuinely clean, rather than to measure with/without-skill
detection delta.
