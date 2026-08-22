---
type: llm
criteria: |
  The pricing-service fixture has correct dependency direction (internal/pricing
  never imports a concrete adapter type), a correctly consumer-declared port
  (RateProvider is declared in internal/pricing, the consumer, not in
  internal/gateway, the producer), and a clean single composition root
  (cmd/pricingapi/main.go is the only place OpenExchangeRatesClient is
  constructed). It has exactly two planted problems:

  1. No context.Context propagation across the port boundary: neither
     RateProvider.FetchRate (internal/pricing/service.go) nor
     Service.Convert (internal/pricing/service.go) accept a context.Context,
     and the concrete implementation
     (internal/gateway/openexchangerates.go's FetchRate) makes its HTTP call
     via c.http.Get(url) instead of building a context-aware request
     (http.NewRequestWithContext + c.http.Do), so a caller has no way to
     cancel, time out, or trace a call that crosses a network boundary. This
     is the scored, validated finding (2026-08-23: no-skill run missed it
     entirely, with-skill run caught it precisely).

  2. (Unscored, bonus signal only) internal/gateway is named after a generic
     architectural role ("gateway") rather than the concrete external
     dependency it wraps — a package-oriented-design naming violation. This
     was NOT caught by either the no-skill or the with-skill run in
     validation, so do not require it for a full score; note it only if the
     response happens to raise it.

  Score 1.0 if finding 1 is identified with correct reasoning and a
  plausible file:line citation, AND the response does not claim
  dependency-direction, port-ownership, or composition-root problems exist
  (they don't — inventing one of those is a false positive and should reduce
  the score even if finding 1 is also present). Score 0 if finding 1 is
  missed, or if the response says the code has no significant issues.
focus: last_message
weight: 1
---

# Judge rubric

This case is a validated with/without-skill discriminator on context
propagation across a port boundary: the no-skill run returned a clean
verdict, the with-skill run correctly named the issue. The package-naming
violation also planted here (internal/gateway) did NOT discriminate — even
the with-skill run explicitly cleared it — so it is intentionally excluded
from scoring. If SKILL.md's package-oriented-design guidance is later
strengthened to catch this pattern, re-validate and fold it back into the
scored criteria.
