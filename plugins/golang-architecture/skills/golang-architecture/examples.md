# Examples

Canonical example stack used below (generic, swap freely): stdlib
`net/http` + `chi` for the HTTP adapter, `sqlc`-generated queries over `pgx`
for the DB adapter. The rule doesn't change if you use Gin/Echo/Fiber or
GORM/sqlx instead — only the adapter package's internals do.

## 1. A compliant package slice (`internal/orders`)

Root package — domain types and the interface the application layer needs.
Zero imports outside the standard library:

```go
// internal/orders/orders.go
package orders

import (
	"context"
	"errors"
	"time"
)

type Order struct {
	ID        string
	CustomerID string
	Total      int64 // cents
	PlacedAt   time.Time
}

var ErrNotFound = errors.New("order not found")

// Repository is declared by the consumer (this package), sized to exactly
// what the application layer calls — not pre-declared by whichever adapter
// happens to implement it.
type Repository interface {
	Get(ctx context.Context, id string) (Order, error)
	Insert(ctx context.Context, o Order) error
}
```

Application layer — orchestrates the use case via the port, never a
concrete driver:

```go
// internal/orders/service.go
package orders

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Place(ctx context.Context, customerID string, total int64) (Order, error) {
	if total <= 0 {
		return Order{}, errors.New("total must be positive")
	}
	o := Order{ID: newID(), CustomerID: customerID, Total: total, PlacedAt: time.Now()}
	if err := s.repo.Insert(ctx, o); err != nil {
		return Order{}, fmt.Errorf("insert order: %w", err)
	}
	return o, nil
}
```

Adapter — the *only* file that imports `pgx`/`sqlc`-generated code for this
concern; satisfies `orders.Repository` without ever importing `orders`
being aware of it:

```go
// internal/orders/postgres/repository.go
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"myapp/internal/orders"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}

func (r *OrderRepository) Get(ctx context.Context, id string) (orders.Order, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, customer_id, total_cents, placed_at FROM orders WHERE id = $1`, id)
	var o orders.Order
	if err := row.Scan(&o.ID, &o.CustomerID, &o.Total, &o.PlacedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return orders.Order{}, orders.ErrNotFound
		}
		return orders.Order{}, fmt.Errorf("scan order: %w", err)
	}
	return o, nil
}

func (r *OrderRepository) Insert(ctx context.Context, o orders.Order) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO orders (id, customer_id, total_cents, placed_at) VALUES ($1, $2, $3, $4)`,
		o.ID, o.CustomerID, o.Total, o.PlacedAt)
	return err
}
```

HTTP adapter — request in, service call, response out. No SQL, no business
rules:

```go
// internal/orders/httpapi/handler.go
package httpapi

import (
	"encoding/json"
	"net/http"

	"myapp/internal/orders"
)

type Handler struct {
	svc *orders.Service
}

func NewHandler(svc *orders.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Place(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string `json:"customer_id"`
		TotalCents int64  `json:"total_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	o, err := h.svc.Place(r.Context(), req.CustomerID, req.TotalCents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(o)
}
```

## 2. Composition root (`cmd/api/main.go`)

The only file that imports the domain package *and* every concrete adapter
package at once:

```go
// cmd/api/main.go
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"myapp/internal/orders"
	"myapp/internal/orders/httpapi"
	"myapp/internal/orders/postgres"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error("connect postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	orderRepo := postgres.NewOrderRepository(pool)   // concrete adapter
	orderSvc := orders.NewService(orderRepo)          // wired via the port
	orderHandler := httpapi.NewHandler(orderSvc)

	r := chi.NewRouter()
	r.Post("/orders", orderHandler.Place)

	log.Info("listening", "addr", ":8080")
	http.ListenAndServe(":8080", r)
}
```

`orders.NewService` never sees `postgres.OrderRepository` — only the
`orders.Repository` interface it satisfies. Swapping Postgres for another
store means writing a new adapter package and changing one line here.

## 3. Package-Oriented Design — directory shape

```
internal/orders/
  orders.go          # domain types + Repository/Service interfaces (root pkg)
  service.go         # application layer — orchestrates via ports
  postgres/
    repository.go    # implements orders.Repository against Postgres
  httpapi/
    handler.go        # implements the HTTP transport for the service
  mock/
    repository.go     # hand-written fake implementing orders.Repository, for tests
```

Grouped **by dependency** (`postgres/`, `httpapi/`) — not
`internal/controllers/orders_controller.go`,
`internal/services/orders_service.go`,
`internal/repositories/orders_repository.go` as siblings across unrelated
concerns. The layer split (domain → application → adapter) still exists —
it's just nested inside `orders/`, not spread across top-level
technical-layer folders.

## 4. Before / after — adding a new package

Adding a hypothetical `webhooks` package that verifies a signature and
persists an event.

**Bad** — signature check and persistence mixed into the HTTP handler, and
the handler imports the DB driver directly:

```go
// ❌ internal/webhooks/handler.go
package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool" // ❌ driver import in the transport layer
)

func Handle(pool *pgxpool.Pool, secret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		mac := hmac.New(sha256.New, secret)
		mac.Write(body)
		expected := mac.Sum(nil)
		if !hmac.Equal(expected, []byte(r.Header.Get("X-Signature"))) {
			http.Error(w, "bad signature", http.StatusUnauthorized)
			return
		}
		pool.Exec(r.Context(), `INSERT INTO webhook_events (payload) VALUES ($1)`, body) // ❌ SQL in the handler
		w.WriteHeader(http.StatusOK)
	}
}
```

**Good** — verification is a pure function in the domain package, storage
goes through a consumer-declared port, and the handler only shapes the
HTTP in/out:

```go
// ✅ internal/webhooks/webhooks.go
package webhooks

import "context"

type Event struct{ Payload []byte }

// Store is declared here, by the consumer, sized to exactly what Service needs.
type Store interface {
	Insert(ctx context.Context, e Event) error
}

// VerifySignature is pure — no I/O, easy to unit test on its own.
func VerifySignature(payload, sig, secret []byte) bool {
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	return hmac.Equal(mac.Sum(nil), sig)
}

type Service struct {
	store  Store
	secret []byte
}

func NewService(store Store, secret []byte) *Service { return &Service{store: store, secret: secret} }

func (s *Service) HandleGitHub(ctx context.Context, payload, sig []byte) error {
	if !VerifySignature(payload, sig, s.secret) {
		return errors.New("invalid signature")
	}
	return s.store.Insert(ctx, Event{Payload: payload})
}
```

```go
// ✅ internal/webhooks/httpapi/handler.go
package httpapi

func (h *Handler) GitHub(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	sig := []byte(r.Header.Get("X-Signature"))
	if err := h.svc.HandleGitHub(r.Context(), body, sig); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}
```

```go
// ✅ internal/webhooks/postgres/store.go
package postgres

type EventStore struct{ pool *pgxpool.Pool }

func (s *EventStore) Insert(ctx context.Context, e webhooks.Event) error {
	_, err := s.pool.Exec(ctx, `INSERT INTO webhook_events (payload) VALUES ($1)`, e.Payload)
	return err
}
```

`VerifySignature` and `Service.HandleGitHub` are now testable with a fake
`Store` and no database — the same testability-follows-from-the-rule
argument the skill's `SKILL.md` makes.
