---
type: llm
criteria: |
  The subscription-service fixture has correct dependency direction (no
  domain-facing file imports database/sql or a driver package), a correctly
  consumer-declared port (Store is declared in internal/subscription, the
  consumer, not in internal/store, the producer), clean context propagation
  (context.Context is threaded from cmd/subscriptionapi/main.go through
  Service.Activate down to Store.FindActive and the concrete
  QueryRowContext call), and a clean single composition root
  (cmd/subscriptionapi/main.go is the only place PostgresStore is
  constructed). It has exactly one planted problem, with two symptoms:

  1. (Root cause) internal/store/postgres.go:8 declares its own sentinel
     error `ErrNotFound`, and internal/subscription/service.go:15 (the
     domain/consumer) imports internal/store (the adapter/producer) solely
     to check `errors.Is(err, store.ErrNotFound)`. This is an
     error-ownership inversion: the domain should declare its own sentinel
     error (owned by the consumer, same principle as port ownership) and the
     adapter should translate its infrastructure-specific errors
     (sql.ErrNoRows) into that domain-owned sentinel at the boundary,
     instead of exporting its own error vocabulary upward for the domain to
     depend on.

  2. (Symptom) internal/subscription/service.go:8 is forced to import
     internal/store for no reason other than the sentinel above — the
     domain package has a compile-time dependency on the infrastructure
     package it would not otherwise need.

  Score 1.0 if the response identifies the error-ownership-inversion problem
  (root cause 1, regardless of exact wording — "error handling as a layering
  concern", "domain checking an infra-owned sentinel", "error contract should
  be consumer-owned" etc. all count) with a plausible file:line citation, AND
  does not claim dependency-direction, context-propagation, or
  composition-root problems exist (they don't — inventing one of those is a
  false positive). Symptom 2 alone (just noting the import, without
  connecting it to error ownership) is worth partial credit (~0.5), not full
  credit. Score 0 if neither is found, or if the response says the code has
  no significant issues.
focus: last_message
weight: 1
---

# Judge rubric

This case was validated (2026-08-23, in isolation) NOT to discriminate skill
vs. no-skill — both runs caught the error-ownership inversion precisely and
independently, with near-identical reasoning. It plays the same role as
case-4: a precision/no-false-positive sanity check (does the skill still
find the real issue without inventing problems in the parts that are
genuinely clean), not a with/without detection-delta test.
