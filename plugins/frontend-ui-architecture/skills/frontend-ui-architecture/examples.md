# Frontend UI Architecture — Examples

Companion to [SKILL.md](SKILL.md). Generic, illustrative examples — swap
the domain names (`Agent`/`Order`/`user`/etc.) for whatever your app
actually models; the shape is what matters.

## Component Folder Anatomy

Canonical shape: `app/orders/_components/OrderCard/`.

```
OrderCard/
├── OrderCard.tsx        # component
├── OrderCard.test.tsx
├── index.ts              # single re-export
├── constants.ts           # e.g. a STATUS_COLOR map
└── helpers.ts               # e.g. statusColor() — uses constants.ts
```

`index.ts` — the safe barrel shape (one named re-export, doubles as default):

```ts
export { OrderCard, OrderCard as default } from "./OrderCard";
```

`helpers.ts` — colocated, component-scoped, not meant to be imported
elsewhere:

```ts
import { STATUS_COLOR } from "./constants";

export function statusColor(status: string): string {
  return STATUS_COLOR[status] ?? "var(--text-secondary)";
}
```

A larger component nests its own `_components/` instead of flattening —
e.g. `OrderDetailsDrawer/_components/{LineItemRow,ShippingSummary,PaymentInfo}/`.

## Constants: Colocate First, Promote Once Reused

**Wrong — pre-creating a global layer before it's earned:**

```
src/
└── constants/
    ├── orders.ts       # only OrderCard uses these
    ├── customers.ts     # only CustomerCard uses these
    └── invoices.ts     # only InvoiceBanner uses these
```

Nothing here is actually shared — it's one file per component, just
moved away from the component that owns it. Now editing `OrderCard`
means opening two folders instead of one.

**Right — colocate, promote only when 2+ consumers need it:**

```
app/orders/_components/OrderCard/constants.ts        # STATUS_COLOR — used only here
app/customers/_components/CustomerCard/constants.ts   # used only here
```

If a value like a shared status-color scale later needs both `OrderCard`
and `InvoiceBanner`, *then* it moves to `lib/` — named for what it is
(`lib/status-colors.ts`), not dumped into a generic `constants.ts`.

## Utils vs Helpers vs Services

**Wrong — a junk-drawer `utils.ts` growing unrelated functions:**

```ts
// lib/utils.ts
export function formatDate(d: Date) { /* ... */ }
export function buildInvoiceUrl(id: string) { /* ... */ }
export function statusColor(status: string) { /* ... */ }   // only OrderCard uses this!
export function debounce(fn: Function, ms: number) { /* ... */ }
```

Mixed scope (generic + component-specific), mixed domain (dates, URLs,
UI color), one file everyone touches → constant merge conflicts.

**Right — split by scope and domain:**

```ts
// lib/format.ts — generic, reusable, named for what it does
export function formatDate(d: Date) { /* ... */ }

// lib/invoice-urls.ts — generic, reusable, named for what it does
export function buildInvoiceUrl(id: string) { /* ... */ }

// app/orders/_components/OrderCard/helpers.ts — component-scoped, not exported beyond the folder
export function statusColor(status: string) { /* ... */ }
```

No generic `services/` folder — instead, a single API-client module plus
one hooks file per domain:

```ts
// lib/api.ts — the ONLY module that knows the API base URL
export async function apiFetch(path: string, init?: RequestInit) { /* ... */ }
```

```ts
// lib/hooks/orders.ts — one file per API domain, wraps lib/api.ts in TanStack Query
import { useQuery } from "@tanstack/react-query";
import { apiFetch } from "../api";

export function useOrders() {
  return useQuery({
    queryKey: ["orders"],
    queryFn: () => apiFetch("/orders"),
  });
}
```

## Business Logic: Hook, Not Component Body

**Wrong — fetch and business logic inline in the component:**

```tsx
function OrdersListView() {
  const [orders, setOrders] = useState([]);
  useEffect(() => {
    fetch(`${process.env.NEXT_PUBLIC_API_BASE}/orders`)
      .then((r) => r.json())
      .then(setOrders);
  }, []);
  // ...render
}
```

Two violations at once: `fetch` outside `lib/api.ts`, and data-fetching
logic inline instead of in a hook.

**Right:**

```tsx
import { useOrders } from "@/lib/hooks/orders";

function OrdersListView() {
  const { data: orders } = useOrders();
  // ...render
}
```

## Barrel Files: Two Different Risk Profiles

**Safe — single re-export, the recommended default:**

```ts
// app/orders/_components/OrderCard/index.ts
export { OrderCard, OrderCard as default } from "./OrderCard";
```

One module, one export. No wildcard, nothing to tree-shake incorrectly,
nothing to create a circular-dependency hotspot.

**Use sparingly — wildcard aggregation, at most one deliberate instance
per app:**

```ts
// lib/hooks/index.ts
/* hooks/ barrel — every data-fetching hook across domains.
   Import from "@/lib/hooks" for convenience, or from a domain file
   directly (e.g. "@/lib/hooks/orders") — both resolve here. */
export * from "./orders";
export * from "./customers";
export * from "./invoices";
```

This is the pattern TkDodo's article warns about (module-graph fan-out,
harder tree-shaking) — it's only tolerable when there's exactly one such
barrel, it stays small, and a comment documents *why* both import paths
need to work. Don't add a second wildcard barrel for a new domain —
prefer `import { useX } from "@/lib/hooks/x"` directly.

## Next.js: An API-Client Module as the DAL-Equivalent Boundary

A real Next.js Data Access Layer (Server Component reading a DB
directly) looks like this:

```ts
// lib/dal/orders.ts (direct-DB pattern — shown for contrast)
import "server-only";

export async function getOrder(id: string) {
  const session = await verifySession();      // auth check inside the DAL
  if (!session) throw new Error("Unauthorized");

  const order = await db.order.findUnique({ where: { id } });
  return { id: order.id, total: order.total };  // DTO, not the raw ORM row
}
```

When the app never reads a database directly and instead calls a separate
backend's API, the equivalent enforcement point is narrower but plays the
same role:

```ts
// lib/api.ts — the one module allowed to know the API base URL
const API_BASE = process.env.NEXT_PUBLIC_API_BASE;

export async function apiFetch(path: string, init?: RequestInit) {
  const res = await fetch(`${API_BASE}${path}`, init);
  if (!res.ok) throw new Error(`API error: ${res.status}`);
  return res.json();
}
```

Everything else — components, hooks, Server Actions — calls through
`apiFetch`, never `fetch` directly against the API base.
