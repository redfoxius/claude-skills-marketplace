# Sources & Rationale

## What this skill is built from

Go doesn't have one canonical "official" backend architecture the way some
frameworks do — the community converges on a handful of compatible ideas:
a de-facto standard project layout, Clean/Hexagonal Architecture adapted to
Go's interface idiom, and Ben Johnson's "package-oriented design" as the
most Go-native expression of the same dependency rule. This skill combines
the parts of each that agree, and names the one place they genuinely
diverge (see "Where Go's idiom differs" below).

- [Standard Go Project Layout (golang-standards/project-layout)](https://github.com/golang-standards/project-layout) — the de-facto `cmd/`+`internal/`+`pkg/` convention this skill's Package Layout section is built on, including its own explicit caveat: treat it as a guide, adopt only what you need.
- [Go Project Structure 2026: Clean Architecture and Best Practices](https://reintech.io/blog/go-project-structure-2026-clean-architecture-best-practices) — "simplicity first, complexity when necessary," and the warning against grab-bag `utils`/`models`/`controllers` packages this skill's anti-patterns list borrows directly.
- [Standard Package Layout — Ben Johnson](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) (mirror: [gobeyond.dev](https://www.gobeyond.dev/standard-package-layout/)) — the four-principle source for this skill's Package-Oriented Design section: root package = domain types, subpackages grouped by dependency, a shared mock subpackage, `main` ties concrete adapters together.
- [Go and a Package Focused Design — Gopher Academy](https://blog.gopheracademy.com/advent-2016/go-and-package-focused-design/) — companion piece on designing from the dependency graph outward rather than the layer diagram inward.
- [Hexagonal Architecture in Golang (Ports and Adapters)](https://medium.com/@sourav.ahmed5654/a-practical-guide-to-hexagonal-architecture-in-golang-0465f53eb2a5) — worked Go example of the Ports & Adapters vocabulary this skill's Dependency Graph section borrows ("port" = interface, "adapter" = implementation).
- [Hexagonal Architecture in Golang: Project Structure, Example & Best Practices — GoLinuxCloud](https://www.golinuxcloud.com/hexagonal-architecture-golang/) — a second full worked example, useful for cross-checking the `cmd`/`internal`/`pkg` mapping against Hexagonal's core/ports/adapters split.
- [How to implement Clean Architecture in Go (Golang) — Three Dots Labs](https://threedots.tech/post/introducing-clean-architecture/) and [Combining DDD, CQRS, and Clean Architecture in Go — Three Dots Labs](https://threedots.tech/post/ddd-cqrs-clean-architecture-combined/) — the source for treating error handling and testability as *consequences* of the dependency rule, not separate concerns.
- [Wild Workouts — Go DDD/Clean Architecture/CQRS example repo](https://github.com/ThreeDotsLabs/wild-workouts-go-ddd-example) — a full runnable reference implementation combining the ideas above; useful when a written explanation isn't enough and you want to read real code.
- [Ardan Labs — service starter-kit (Domain-Driven, Data-Oriented Go services)](https://github.com/ardanlabs/service) — production-grade precedent for "manual wiring in `main.go` is the default, not a compromise" and for structuring exported/unexported APIs within a package.
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md) — source for this skill's Error Handling section: `%w` vs `%v`, exporting sentinel errors for `errors.Is`/`errors.As`, adding context when wrapping.
- [Go Dependency Injection Approaches — Wire vs. Fx vs Manual](https://leapcell.io/blog/go-dependency-injection-approaches-wire-vs-fx-and-manual-best-practices) — source for this skill's Composition Root escalation ladder (manual → Wire → Fx) and the tradeoffs at each step.
- [Compile-time DI with Go Cloud's Wire — go.dev blog](https://go.dev/blog/wire) — official framing of Wire as "the same code you'd write by hand, generated," which is why it's the first escalation and not a departure from manual wiring's guarantees.
- [sqlc vs GORM vs sqlx: Go Database Libraries Compared 2026](https://reintech.io/blog/sqlc-vs-gorm-vs-sqlx-go-database-libraries-compared-2026) and [Comparing database/sql, GORM, sqlx, and sqlc — JetBrains Go blog](https://blog.jetbrains.com/go/2023/04/27/comparing-db-packages/) — background for why `examples.md` uses sqlc/pgx as the canonical DB-adapter example (type-safety, most Go-idiomatic) while naming GORM as a valid alternative for teams prioritizing velocity.
- [Go Chi vs Gin vs Echo: Web Framework Comparison 2026](https://reintech.io/blog/go-chi-vs-gin-vs-echo-web-framework-comparison-2026) and [Gin vs Echo vs Fiber 2026 — Encore](https://encore.dev/articles/gin-vs-echo-vs-fiber) — background for why `examples.md` uses stdlib `net/http`+chi (closest to idiomatic Go, easiest to map onto the dependency rule) while naming Gin/Echo/Fiber as swappable adapters — the rule doesn't change with the router choice.

## Where Go's idiom differs from the source material's usual framing

Most Onion/Hexagonal/Clean Architecture write-ups (including this repo's
own `onion-architecture` skill, written for TypeScript) implicitly assume
the *producer* declares a port interface — a `LLMProvider` interface lives
next to the concrete adapters that might implement it, and consumers import
that shared interface type.

Go's standard idiom inverts this: **"accept interfaces, return structs."**
The *consumer* package declares the interface it needs, sized to exactly
the methods it calls — there's no explicit `implements` keyword, so nothing
stops a `postgres.OrderRepository` from satisfying an interface it has
never heard of. This is why Package-Oriented Design puts the interface in
the *root* package (the consumer side) rather than in a shared
`ports`/`interfaces` package that every adapter imports — a shared ports
package is itself a small violation of "group by dependency": it becomes a
dependency hub that both the domain and every adapter import, which is
exactly the grab-bag-package anti-pattern this skill otherwise warns
against.

## How the pieces map onto each other

| Concept | Onion/Hexagonal/Clean naming | This skill's Go framing |
|---|---|---|
| Innermost layer | Domain / Entities | Root package of a bounded concern (`orders.Order`, `orders.Repository`) |
| Use-case orchestration | Application / Use Cases | A `Service` struct or plain functions in the same or a sibling file, calling ports |
| Boundary interfaces | Ports | Interfaces declared in the root package, by the consumer |
| Concrete implementations | Adapters | Subpackages named after what they depend on (`postgres/`, `httpapi/`, `redis/`) |
| Wiring | Composition Root | `cmd/<app>/main.go` (manual), escalating to `google/wire` or `uber-go/fx` |
| Framework choice | (implementation detail) | Router (Gin/Echo/Fiber/Chi/stdlib) and DB layer (GORM/sqlc/sqlx/pgx) — swappable adapters, never referenced by the rule itself |

## What this skill deliberately doesn't cover

If you pair this skill with router- or ORM-specific skills later, they'd
own:

- **Framework mechanics** — how to register a route, bind/validate a
  request body, configure middleware for whichever of Gin/Echo/Fiber/Chi/
  stdlib you picked. This skill only cares that the mechanics live in an
  adapter package, not what they look like inside it.
- **ORM/query mechanics** — schema definitions, migrations, transaction
  scoping, N+1 avoidance for whichever of GORM/sqlc/sqlx/pgx you picked.
  This skill only cares that the queries live in an adapter package
  satisfying a consumer-declared interface.
- **Concurrency patterns** beyond the one layering-relevant rule (propagate
  `context.Context` across port boundaries) — goroutine lifecycle,
  channel patterns, and worker-pool design are a separate concern from
  package boundaries.
