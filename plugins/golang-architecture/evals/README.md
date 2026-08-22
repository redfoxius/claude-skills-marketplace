# golang-architecture evals

Five independent eval cases, each a tiny standalone Go project under
`case-N-*/fixtures/` with a small number of architectural violations planted,
with no comments in the fixtures revealing what's wrong.

| Case | Domain | Violations planted |
|---|---|---|
| case-1-order | order service (postgres + slack) | dependency rule, port inversion, composition root (adapter built inline in constructor) |
| case-2-billing | billing service (kafka) | dependency rule, port inversion, composition root (adapter built twice in two constructors, `main.go` wires nothing) |
| case-3-inventory | inventory service (gorm) | dependency rule, port inversion, composition root (connection opened via package-level `init()`) |
| case-4-shipping | shipping service (carrier client) | **sanity check** (originally designed as a discriminating "trap"; validated not to discriminate — see below) — dependency direction and composition root are both clean; the only problem is the `Carrier` port declared in the producer package instead of the consumer. |
| case-5-pricing | pricing service (FX rate gateway) | **validated discriminator** — dependency direction, port ownership, and composition root are all clean. Scored finding: no `context.Context` propagated across the `RateProvider` port boundary down to the HTTP call — validated 2026-08-23 in isolation (no-skill: missed entirely; with-skill: caught precisely). A second planted problem (adapter package named `gateway`, a generic role, instead of `openexchangerates`, its concrete dependency) did NOT discriminate — the with-skill run explicitly cleared it — so it's kept unscored as a bonus signal only. See case-5's case.yaml for detail. |

### Why case-4 is a sanity check, not a trap

It was originally built to test whether a generic review (no golang-architecture
skill) would miss the port-location problem, since "interface next to its
implementation" is conventional in many non-Go architectural traditions. On
2026-08-23 it was validated manually (see workflow below) two ways:

1. Reviewed together with cases 1-3 in one batch — both the with-skill and
   no-skill runs caught it, but this run is contaminated: having just spotted
   the same producer-declared-port pattern three times in a row primes the
   reviewer to look for it a fourth time, regardless of skill.
2. Re-reviewed case-4 **in full isolation** (no other cases in context) — both
   runs still caught it, with explicit "Go idiom" / "hexagonal inversion"
   reasoning even without the skill.

Conclusion: this specific pattern is common-enough general Go knowledge that
it doesn't discriminate skill vs. no-skill, so case-4 was repurposed as a
precision/no-false-positive check instead (does the skill still find the real
issue without inventing problems in the parts that are genuinely clean).
A real discriminator would need to target something more skill-specific and
less "folk-known" — the isolated with-skill run for this fixture surfaced two
rules that go beyond the three-axis assumption these eval cases were built
around (context propagation across a port boundary; package-oriented naming —
grouping a port with its adapter under a generic name like `carrier` instead
of naming the adapter package after what it depends on, e.g. `fedex`). Those
would be a better starting point for a genuinely discriminating case-5, since
they're less likely to already be general reviewer folk wisdom — sourced from
the skill's actual SKILL.md rather than its one-line catalog description.

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
