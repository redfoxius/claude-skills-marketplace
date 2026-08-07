---
name: frontend-ui-architecture
description: "Defines where frontend code lives and how it's layered: folder/feature structure, component-folder anatomy, component splitting, constants placement, utils vs helpers vs services, business-logic placement (custom hooks, API/service boundary), types organization, and barrel-file (index.ts) conventions — for plain React and Next.js App Router apps. Use when scaffolding a new feature/component, deciding where a file or piece of logic belongs, or reviewing/refactoring project structure. Architecture and code organization only — NOT component behavior/anti-patterns (see react-best-practices) and NOT Next.js routing mechanics, data-fetching primitives, or performance (see next-best-practices)."
version: "1.0.0"
---

# Frontend UI Architecture

Where frontend code lives and how it's layered. Complements two sibling
skills instead of overlapping them:

- **`react-best-practices`** owns component *behavior* — purity, hooks
  rules, memoization, anti-patterns.
- **`next-best-practices`** owns App Router *mechanics* — which special
  file, which data-fetching primitive, caching, performance.
- **This skill** owns *where things live and how they're layered* — for
  plain React feature work and for Next.js App Router apps.

These are default conventions distilled from widely-used React/Next.js
project-structure practice (see [README.md](README.md) for sources). If the
project you're working in already has an established, consistent structure
of its own, follow that project's own convention first — this skill fills
the gaps where no local convention exists yet, or explains the reasoning
behind adopting one. For code examples, see [examples.md](examples.md).

## Severity Levels

- **CRITICAL** — wrong layer for logic/data-access; will cause coupling,
  untestability, or duplicated business rules
- **HIGH** — structure that causes real scaling pain (merge conflicts,
  unclear ownership, cross-feature coupling)
- **MEDIUM** — organizational polish; wrong today, cheap to fix later

---

## Folder & Feature Structure (HIGH)

- Prefer route-colocation over a separate top-level `features/` directory:
  put a component under the route that owns it, in a private `_components/`
  folder (the underscore prefix excludes it from Next.js routing). Reserve
  top-level `components/` for UI reused across **2+ routes**.
  - Typical shape: `<route>/_components/<Name>/` for route-local
    components; cross-route pieces (a shared shell, a shared data-table,
    etc.) live in a top-level `src/components/`.
- Nest `_components/` again inside a component folder once it grows
  internal subcomponents, instead of flattening everything into one folder
  (e.g. `OrderDetailsDrawer/_components/LineItemRow/`).
- Don't introduce a parallel `features/` directory outside `app/` unless
  the project is genuinely multi-team or has 20+ features — that's
  Feature-Sliced-Design territory (see README) and adds indirection most
  apps don't pay back until well past that size.
- Unidirectional dependency flow: `components/` / `lib/` (shared) →
  route-local `_components/` → the route itself. A shared component must
  never import from a specific route's `_components/`.
- If the project hand-copies or vendors cross-package contracts/types (a
  shared-types directory synced between a frontend and backend, for
  example), don't restructure them under this skill's rules — they follow
  whatever sync convention the project already documents for them.

## Component Folder Anatomy (HIGH)

A consistent shape to default to for every component beyond trivial size:

```
ComponentName/
  ComponentName.tsx      # the component
  ComponentName.test.tsx
  index.ts                # single re-export — the folder's public API
  constants.ts             # colocated, component-scoped constants (optional)
  helpers.ts               # colocated pure helper functions (optional)
  styles.ts                 # colocated style objects, if not pure Tailwind (optional)
  hooks/                     # colocated hooks — only once there's >1, or one reused within the folder
```

- Skip the folder entirely for a trivial one-off component — a single
  `.tsx` file is fine until it grows a sibling file.
- `index.ts` re-exports exactly one thing
  (`export { X, X as default } from "./X"`) — see **Barrel Files** below
  for why this specific shape is the safe kind.

## Component Splitting & Composition (HIGH)

Sizing/purity rules (max lines, max props, pure render) live in
`react-best-practices` — this skill only adds *where the split lands on
disk*.

- When a component splits into a data-fetching half and a rendering half,
  the data-fetching half is a hook (colocated `hooks/`, or `lib/hooks/`
  for cross-component reuse) — not a wrapper "container" component. The
  presentational half stays the folder's default export.
- Don't over-split: a subcomponent or hook with exactly one caller and no
  independent meaning is easier to read inlined back into its parent.
- Atomic design (atoms/molecules/organisms) is a design-system pattern,
  not a feature-code pattern — skip it for route/`_components/` work. It
  earns its place mainly once a shared UI-kit package grows into a real
  multi-app design system.

## Constants (MEDIUM)

- Colocate: a `constants.ts` inside the component folder that uses it.
  There is usually no need for a global `constants/` directory, and none
  should be added speculatively.
- Promote a constant to a shared location (`lib/`) only once it's actually
  reused across **2+ unrelated** component folders.
- A one-off value used exactly once (a single magic padding number)
  doesn't earn a constants file — inline it.

## Utils vs Helpers vs Services (MEDIUM)

- `helpers.ts` colocated in a component folder = logic specific to that
  component only (e.g. a card component's own formatting helper). Not
  meant to be imported by anything else.
- Cross-component pure functions (formatting, URL building) go in `lib/`
  as a file named for what it does (`lib/format.ts`, `lib/urls.ts`) —
  never a generic `utils.ts` junk drawer that accumulates unrelated
  functions.
- For apps that call a backend API rather than touching a database
  directly, the `services/` boundary is usually a single `lib/api.ts` (or
  `lib/http-client.ts`) — the **only** module allowed to know the API base
  URL or call `fetch`/`axios` against it — plus one hook file per API
  domain (`lib/hooks/<domain>.ts`), wrapping it in TanStack Query (or
  equivalent) hooks. A new API domain gets a new `lib/hooks/<domain>.ts`,
  never a scattered `fetch()` inside a component.

## Business Logic Placement (CRITICAL)

- All data fetching/mutation logic lives in a `lib/hooks/<domain>.ts`
  data-fetching hook, never inline in a component body — this is the
  concrete file target for `react-best-practices`' "Data Fetching" rule.
- Component-local business logic (derived values, formatting specific to
  one component) → colocated `helpers.ts`.
- Logic reused by 2+ components but not generic enough for a shared
  `utils` module → promote to `lib/`, named by domain (what it's for), not
  by type (what kind of file it is).

## Data Access Boundary — Next.js Specifics (CRITICAL)

- The common industry pattern for App Router is a server-only Data Access
  Layer (`import 'server-only'`, authorization checks inside, DTOs out) —
  for apps where Server Components read a database directly.
- If the app instead always goes through a separate backend's REST/GraphQL
  API rather than touching the database directly, the equivalent boundary
  is the app's API client module (e.g. `lib/api.ts`): treat it like a DAL.
  Nothing outside that module and its per-domain hooks should construct a
  request URL or call `fetch` against the API — that's the enforcement
  point, playing the same role a real DAL plays for direct DB access.
- If a Server Component ever needs to read a database directly (bypassing
  a backend API), reach for the real DAL pattern instead — `server-only`,
  auth-in-layer, DTOs out — see README sources for the canonical reference.

## Types (MEDIUM)

- Component-specific types stay inline in the component file — avoid
  creating a `types.ts` at the single-component scope.
- Cross-cutting types live centrally in a shared location
  (`lib/types.ts`, `lib/models.ts`, or similar) — promote a type there
  only once 2+ components need it.
- Contracts shared with a backend (or another package) should have exactly
  one source of truth — if the project hand-copies/vendors shared types,
  follow its documented sync convention rather than inventing a third copy.

## Barrel Files — `index.ts` (MEDIUM)

Two different things get called "barrel files" — treat them very
differently:

- **Safe, and the recommended default**: a per-component `index.ts`
  re-exporting exactly one component (`export { X, X as default } from
  "./X"`). This is a stable public-API seam for a single module, not a
  wildcard re-export, so it doesn't carry the bundle-bloat /
  circular-dependency cost that "avoid barrel files" advice targets.
- **Use sparingly**: a directory-wide `export *` barrel aggregating many
  modules (e.g. a `lib/hooks/index.ts` re-exporting every domain's hooks).
  If you add one, document inline *why* both import paths (the barrel and
  the specific file) need to resolve, and keep it to at most one such
  barrel per app rather than letting every directory grow one — prefer
  importing a specific `lib/hooks/<domain>` file directly.

## Naming (MEDIUM)

- PascalCase for component folders/files; camelCase for everything else
  (hooks, helpers, constants exports) — matches the Airbnb convention.
- Hook files and functions start with `use`.
