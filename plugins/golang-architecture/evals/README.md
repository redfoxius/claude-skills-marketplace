# golang-architecture evals

Six independent eval cases, each a tiny standalone Go project under
`case-N-*/fixtures/` with a small number of architectural violations planted,
with no comments in the fixtures revealing what's wrong.

| Case | Domain | Violations planted |
|---|---|---|
| case-1-order | order service (postgres + slack) | dependency rule, port inversion, composition root (adapter built inline in constructor) |
| case-2-billing | billing service (kafka) | dependency rule, port inversion, composition root (adapter built twice in two constructors, `main.go` wires nothing) |
| case-3-inventory | inventory service (gorm) | dependency rule, port inversion, composition root (connection opened via package-level `init()`) |
| case-4-shipping | shipping service (carrier client) | **sanity check** (originally designed as a discriminating "trap"; validated not to discriminate — see below) — dependency direction and composition root are both clean; the only problem is the `Carrier` port declared in the producer package instead of the consumer. |
| case-5-pricing | pricing service (FX rate gateway) | **validated discriminator** — dependency direction, port ownership, and composition root are all clean. Scored finding: no `context.Context` propagated across the `RateProvider` port boundary down to the HTTP call — validated 2026-08-23 in isolation (no-skill: missed entirely; with-skill: caught precisely). A second planted problem (adapter package named `gateway`, a generic role, instead of `openexchangerates`, its concrete dependency) did NOT discriminate — the with-skill run explicitly cleared it — so it's kept unscored as a bonus signal only. See case-5's case.yaml for detail. |
| case-6-subscription | subscription service (postgres store) | **sanity check** (originally designed as a candidate discriminator; validated not to discriminate — see below) — dependency direction, port ownership, context propagation, and composition root are all clean. The one planted problem is "error handling as a layering concern": the adapter (`internal/store`) correctly translates `sql.ErrNoRows` into its own `ErrNotFound` sentinel, but the domain (`internal/subscription`) then imports the adapter package solely to check that sentinel via `errors.Is`, instead of owning its own domain-level sentinel — an error-ownership inversion structurally parallel to case-4's port-ownership inversion. |

### Score so far: 1 of 3 candidate discriminators actually discriminate

Three "does this need the skill to catch it" hypotheses have been validated
in full isolation so far:

| Candidate | Rule | Result |
|---|---|---|
| case-4-shipping | port declared by producer, not consumer | **Did not discriminate** — both runs caught it, citing "Go idiom" / "hexagonal inversion" |
| case-5-pricing | missing `context.Context` across a port boundary | **Discriminates** — no-skill run: "no significant problems"; with-skill run: caught precisely |
| case-6-subscription | domain branching on an infra-owned error sentinel | **Did not discriminate** — both runs caught it, citing "error contract should be consumer-owned" |

Takeaway: interface-shaped rules (port ownership, error-sentinel ownership —
"the consumer should own its own contract") turn out to already be
well-internalized general Go/software-design knowledge among capable
reviewers, even without this skill. Context propagation across a boundary is
the one validated exception so far — it's the kind of thing reviewers know
*abstractly* but reliably forget to actually check for on a first pass.
case-4 and case-6 are kept as precision/no-false-positive sanity checks
instead (does the skill still find the real issue without inventing problems
in the parts that are genuinely clean) rather than detection-delta tests.

Each validation followed the same protocol (see workflow below): first a
batched run (contaminated by priming — seeing the same pattern repeat across
cases teaches the reviewer to look for it again, regardless of skill), then a
rerun in full isolation as the actual measurement.

If designing a future case-N candidate discriminator, favor rules that are
narrow, mechanical, and easy to forget to check rather than ones that restate
a well-known design principle (DIP/ISP) in a new location — "consumer owns
the X" pattern-matches too easily once a reviewer has seen it once, in any
form.

Each case follows the native `claude plugin eval` layout:

```
case-N-*/
  case.yaml        # metadata + prompt_file + context.add_dirs: [fixtures/]
  prompt.md         # the review instruction given to the agent under test
  fixtures/          # the Go project being reviewed
  graders/
    findings.md      # llm grader scoring coverage of the planted findings
    skill-fired.md    # with_only tool_used grader confirming the skill fired
```

## Status: blocked on early access

`claude plugin eval` is currently gated behind Anthropic early access and is
**not enabled** in this environment. The `case.yaml`/grader schema here is
authored from `claude plugin eval --help` plus best-effort documentation and
has **not been validated against a real run** — expect to adjust field names
once early access lands and `claude plugin eval plugins/golang-architecture`
can actually execute.

Once enabled, the intended CI command is:

```bash
claude plugin eval plugins/golang-architecture \
  --ablation with-without \
  --threshold 0.8 \
  --json evals/results/latest.json
```

`--ablation with-without` runs each case's `prompt.md` once with the skill
available and once without, and reports the score delta — this is the
with/without comparison this suite was built to automate.

## Running these evals manually today (no early access needed)

Until early access lands, evaluate a case by hand:

1. Spawn a fresh agent, give it the case's `prompt.md` content plus the
   absolute path to its `fixtures/` directory, and explicitly tell it not to
   invoke the Skill tool.
2. Spawn a second fresh agent with the same prompt, but tell it to first
   invoke the `golang-architecture` skill via the Skill tool, then review.
3. Compare each run's findings against `graders/findings.md`'s criteria by
   hand.

This is exactly how case-1 through case-3 were validated before this eval
suite existed (see the session that authored this directory).
