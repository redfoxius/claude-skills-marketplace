---
name: golang-architecture
description: "Forces Go-idiomatic package boundaries and dependency direction for backend modules: domain/business-logic packages never import a concrete infrastructure package (a DB driver, an ORM, a specific HTTP framework, a broker client) directly, ports are interfaces declared by the consumer package rather than pre-declared by the producer (the Go-idiomatic inversion of classic Onion/Hexagonal ports), and every concrete adapter is constructed in one composition root (cmd/<app>/main.go by default, google/wire for large static graphs, uber-go/fx only for lifecycle-heavy modular apps). Use when starting a new Go service, adding a new package or external integration, reviewing a Go PR for layering violations, or deciding which package new logic belongs in. Architecture and dependency direction only — NOT framework routing mechanics (Gin/Echo/Fiber/Chi handler syntax) and NOT ORM/query mechanics (GORM/sqlc/sqlx/pgx query-writing)."
version: "1.0.0"
---

# Golang Architecture (Backend)

Where Go package dependencies are allowed to point, and who's allowed to
import what. This skill is stack-agnostic — it doesn't assume Gin vs. Echo
vs. Chi, or sqlc vs. GORM vs. sqlx — because it owns *the dependency rule*,
not any framework's mechanics:

- **A router/framework skill** (if you have one) would own HTTP *mechanics*
  — middleware, binding, validation.
- **A DB/ORM skill** (if you have one) would own query *mechanics* — schema,
  migrations, transactions.
- **This skill** owns *the dependency rule* — which package a piece of code
  belongs in, and what it's forbidden to import.

For code examples, see [examples.md](examples.md). For sources and how Go's
idiom compares to canonical Onion/Hexagonal/Clean Architecture, see
[README.md](README.md).

## The Dependency Graph

```
   ┌─────────────────────────────────────────────────┐
   │ Infrastructure   (postgres/, redis/, http/, …)   │
   │  ┌─────────────────────────────────────────────┐ │
   │  │ Interface adapters (HTTP handlers,           │ │
   │  │  repository structs implementing a port)      │ │
   │  │  ┌───────────────────────────────────────┐   │ │
   │  │  │ Application (use-case functions /      │   │ │
   │  │  │  services orchestrating ports)          │   │ │
   │  │  │  ┌─────────────────────────────────┐  │   │ │
   │  │  │  │ Domain (root package: types,      │  │   │ │
   │  │  │  │  business rules — zero deps)       │  │   │ │
   │  │  │  └─────────────────────────────────┘  │   │ │
   │  │  └───────────────────────────────────────┘   │ │
   │  └─────────────────────────────────────────────┘ │
   └─────────────────────────────────────────────────┘
        dependencies point INWARD only  →
   ports are interfaces DECLARED by domain/application,
   SATISFIED by adapters — never the other way around
```

- **Domain** — the package's root types and pure business rules. No
  `database/sql`, no HTTP, no broker client — its only "side effect" is
  calling an interface it declares itself.
- **Application** — use-case functions or a thin `Service` struct that
  orchestrate the domain and call ports (interfaces) it or the domain
  package declares. Depends on interface *types*, never on a concrete
  adapter package.
- **Interface adapters** — HTTP handlers (request in / response out) and
  repository structs (the concrete type that *implements* a port, e.g. a
  `postgres.OrderRepository` implementing an `orders.Repository` interface).
- **Infrastructure** — the concrete driver/client packages (`pgx`, `gorm`,
  `redis`, a specific router) and the composition root — the only place a
  concrete adapter gets constructed and handed to the application layer.

## Severity Levels

- **CRITICAL** — breaks the dependency rule: a domain/application package
  imports a concrete infra package directly, or an adapter is constructed
  outside the composition root.
- **HIGH** — wrong package for the logic, a producer-declared port instead
  of a consumer-declared one, or a composition root that does real business
  logic instead of wiring.
- **MEDIUM** — organizational polish (package layout, premature interfaces,
  context-propagation gaps) — wrong today, cheap to fix later.

---

## Package Layout (MEDIUM)

- Default shape: `cmd/<app>/main.go` (composition root) + `internal/<name>/`
  per bounded concern + an optional top-level `pkg/` only for code meant to
  be imported by *other* repositories. Treat this as a guide, not a mandate
  — a small service can be one `internal/` package; add structure when a
  real seam appears, not speculatively.
- Ban grab-bag packages: `utils`, `common`, `helpers`, `models`, `types` as
  a dumping ground for unrelated code. If a function only makes sense next
  to the type it operates on, it belongs in that type's package, not a
  shared bucket every other package ends up importing.
- Prefer grouping `internal/` subpackages **by bounded concern**
  (`internal/orders/`, `internal/billing/`), not by technical layer
  (`internal/controllers/`, `internal/services/`, `internal/repositories/`)
  — the layer split still exists *inside* each concern's package (see
  Package-Oriented Design below), it just isn't the top-level axis.

## The Dependency Rule — Ports Declared by the Consumer (CRITICAL)

- Dependencies point inward only. A domain or application package must
  never import a concrete driver/ORM/framework package directly (`pgx`,
  `gorm.io/gorm`, `redis`, `gin`, a specific message-broker client) — only
  through an interface.
- **Go inverts the classic OOP framing here**: in Onion/Hexagonal write-ups
  for other languages, the producer (the infrastructure side) usually
  pre-declares the port interface. In Go, the idiom is "accept interfaces,
  return structs" — **the consumer package declares the interface it
  needs**, sized to exactly the methods it calls, and any adapter that
  happens to satisfy that method set works, with no explicit `implements`
  relationship. A `postgres.OrderRepository` doesn't need to know an
  `orders.Repository` interface exists; it just needs to have the right
  method signatures.
- If you catch yourself writing `import "gorm.io/gorm"` or a driver-specific
  import inside a use-case/service file, that's the rule breaking — push
  the concern one layer outward into the adapter package, don't special-case
  it.

## Package-Oriented Design (HIGH)

Ben Johnson's "Standard Package Layout" is the concrete Go-native shape of
the dependency rule — group subpackages **by dependency**, not by technical
layer:

- The **root package** of a bounded concern holds domain types and the
  interfaces the application layer needs (e.g. `orders.Order`,
  `orders.Repository`, `orders.Service`) — zero imports of anything outside
  the standard library.
- **Subpackages named after what they depend on**, not what they do:
  `postgres/` (implements `orders.Repository` against Postgres), `http/`
  (exposes the service over HTTP), `redis/` (a cache adapter) — never
  `controllers/`, `services/`, `repositories/` as top-level siblings, which
  groups by *role* instead of by *what would force a rewrite if swapped*.
- A shared `mock/` (or `mocks/`) subpackage holds test doubles for the root
  package's interfaces, used by every consumer's tests.
- `cmd/<app>/main.go` is the only place that imports both the root package
  and every concrete subpackage at once, to wire them together.

## Composition Root (CRITICAL)

- Default: **explicit manual wiring in `cmd/<app>/main.go`** — construct
  every concrete adapter, pass it into the service/use-case constructor,
  start the server. This is the Go-idiomatic choice: explicit, no
  reflection, compiler-checked, and it's what most production Go codebases
  (Ardan Labs' service starter-kit included) actually do.
- Escalate only when the manual version gets unwieldy:
  - **`google/wire`** for a large but *static* dependency graph — it
    generates the same explicit wiring code you'd hand-write, at compile
    time, with zero runtime cost. Reach for it when `main.go` has grown into
    a genuinely hard-to-follow chain of constructors, not preemptively.
  - **`uber-go/fx`** only for large, highly modular applications that need
    real lifecycle management (start/stop hooks across many independently
    developed modules) — it's a runtime DI framework with real overhead and
    "magic," which cuts against Go's explicitness bias, so it needs to earn
    its place.
- Nothing outside the composition root should ever call `postgres.New(...)`,
  construct an HTTP client, or open a Redis connection. A package needing a
  new capability takes it as a constructor parameter, typed as the
  interface it declares — the composition root decides what satisfies it.

## Error Handling as a Layering Concern (HIGH)

- Use sentinel errors (`var ErrNotFound = errors.New(...)`) or custom error
  types at a package's boundary when a caller needs to branch on *which*
  error occurred; plain `fmt.Errorf` otherwise.
- Wrap with `%w` (not `%v`) when a caller one layer up should be able to
  match the underlying cause with `errors.Is`/`errors.As`; add context to
  the message as you wrap (`fmt.Errorf("get order %s: %w", id, err)`) so a
  bare `"connection refused"` becomes traceable to the call that produced
  it.
- Don't let a panic cross a package boundary as your error-handling
  strategy — recover only at a single well-known boundary (typically the
  HTTP middleware layer, to turn a panic into a 500), never scattered
  through business logic.

## Context Propagation Across Ports (MEDIUM)

- Any method that crosses a port boundary — a repository call, an HTTP
  client call, a queue publish — takes `context.Context` as its first
  parameter, propagated from the request/job that triggered it, not
  `context.Background()` conjured mid-stack.
- This is the one concurrency idiom that's actually a layering rule: it's
  what lets an adapter enforce request-scoped cancellation/timeouts/tracing
  without the application layer knowing *how* — the domain and application
  layers just pass the context through, they don't inspect it.

## Testability Follows From the Rule (HIGH)

- Because the application layer depends on interfaces it declares itself,
  unit tests construct the service with a hand-written fake (or one from the
  package's `mock/` subpackage) satisfying that interface — no real
  Postgres, no real HTTP server, no `testcontainers` for a business-logic
  test.
- Reserve real-dependency integration tests (real Postgres via
  `testcontainers-go`, a real Redis) for the adapter packages themselves —
  testing that `postgres.OrderRepository` actually round-trips an `Order`
  correctly, not for testing business rules that live one layer in.
- If a "unit" test needs a running database to pass, that's a signal the
  code under test reached past its port, not a signal to accept slower unit
  tests.

## Anti-Patterns / Red Flags (CRITICAL)

Treat any of these as a layering violation to fix, not a style nit:

- `import "gorm.io/gorm"`, a specific SQL driver, or a router package
  (`gin`, `echo`) inside a domain or use-case/service file.
- A concrete adapter (`postgres.New(...)`, `redis.NewClient(...)`, an HTTP
  client) constructed anywhere outside the composition root.
- An interface pre-declared by the producer package and imported by every
  consumer, instead of each consumer declaring the (smaller) interface it
  actually needs.
- A `cmd/<app>/main.go` with real business logic in it (validation rules,
  calculations, conditional workflow) instead of pure wiring.
- A grab-bag `utils`/`common`/`models` package that every other package
  ends up importing, creating an accidental dependency hub.
- Package-by-layer top-level structure (`internal/controllers`,
  `internal/services`, `internal/repositories` as siblings) instead of
  package-by-dependency or package-by-concern.
