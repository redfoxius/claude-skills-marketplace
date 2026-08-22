# golang-architecture evals

Four independent eval cases, each a tiny standalone Go project under
`case-N-*/fixtures/` with exactly three (case-4: one) architectural
violations planted, with no comments in the fixtures revealing what's wrong.

| Case | Domain | Violations planted |
|---|---|---|
| case-1-order | order service (postgres + slack) | dependency rule, port inversion, composition root (adapter built inline in constructor) |
| case-2-billing | billing service (kafka) | dependency rule, port inversion, composition root (adapter built twice in two constructors, `main.go` wires nothing) |
| case-3-inventory | inventory service (gorm) | dependency rule, port inversion, composition root (connection opened via package-level `init()`) |
| case-4-shipping | shipping service (carrier client) | **trap case** — dependency direction and composition root are both clean; the only problem is the `Carrier` port declared in the producer package instead of the consumer. Designed to catch reviewers who don't apply the Go-idiomatic consumer-owns-the-port rule specifically — a generic review is expected to call this code "well-architected" and miss it. |

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
