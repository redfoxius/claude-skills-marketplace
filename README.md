# claude-skills-marketplace

A [Claude Code](https://claude.com/claude-code) plugin marketplace with
skills and agents I use in my own projects.

## Install

```
/plugin marketplace add redfoxius/claude-skills-marketplace
/plugin install golang-architecture@claude-skills-marketplace
/plugin install planner@claude-skills-marketplace
```

## Plugins

| Plugin | Description | Release |
|--------|-------------|---------|
| [golang-architecture](plugins/golang-architecture) | Forces Go-idiomatic package boundaries and dependency direction: domain/application packages never import a concrete infra package directly, ports are interfaces declared by the consumer, one composition root wires every adapter. | [v1.1.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/golang-architecture-v1.1.0) |
| [frontend-ui-architecture](plugins/frontend-ui-architecture) | Where frontend code lives and how it's layered — folder/feature structure, component-folder anatomy, business-logic placement, types organization, barrel-file conventions for React and Next.js App Router. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/frontend-ui-architecture-v1.0.0) |
| [marketplace-release](plugins/marketplace-release) | Publishes a skill to this repo — scaffolding a new skill, or running `scripts/release_skill.py` to zip, tag, publish a GitHub Release, and update this table's release link for an existing one. | [v1.1.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/marketplace-release-v1.1.0) |
| [planner](plugins/planner) | Turns a feature/bugfix request into a structured Development Plan, reading each touched module's architectural constraints and gotchas first. Never writes code. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/planner-v1.0.0) |
| [implementer](plugins/implementer) | Executes an already-written Development Plan, applying the relevant project skills per file and self-verifying with the package's own test/typecheck commands. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/implementer-v1.0.0) |
| [architecture-reviewer](plugins/architecture-reviewer) | Read-only review of a diff/PR/branch/directory for architectural-boundary violations, routing each file to the right architecture skill and citing every finding to file:line. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/architecture-reviewer-v1.0.0) |
| [plan-verifier](plugins/plan-verifier) | Checks finished code against every point of a Development Plan and its Implementation Report, with a MET/NOT MET/UNVERIFIABLE verdict per criterion backed by real evidence. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/plan-verifier-v1.0.0) |
| [test-writer](plugins/test-writer) | Writes or extends tests for UI and backend code following each package's own conventions, always self-verifying with the real test command. Test files only. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/test-writer-v1.0.0) |
| [doc-writer](plugins/doc-writer) | Turns an already-implemented feature into feature-facing reference/explanation documentation with Mermaid diagrams, linking out to existing docs instead of duplicating them. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/doc-writer-v1.0.0) |
| [researcher](plugins/researcher) | Read-only research agent — repository search or external documentation lookup, every finding backed by a file:line or source citation. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/researcher-v1.0.0) |

## Adding a new skill

1. Create `plugins/<name>/.claude-plugin/plugin.json` and
   `plugins/<name>/skills/<name>/SKILL.md`.
2. Add an entry to `plugins` in `.claude-plugin/marketplace.json` with
   `"source": "./plugins/<name>"`.
3. Add a row to the `## Plugins` table above, with `—` in the Release
   column — `scripts/release_skill.py` fills it in on the first release.
4. Commit and push — no build step, no publish action.

## Adding a new agent

Same shape as a skill, but Claude Code discovers agents from an `agents/`
directory instead of `skills/`:

1. Create `plugins/<name>/.claude-plugin/plugin.json` and
   `plugins/<name>/agents/<name>.md` (the subagent file itself — `name`,
   `description`, `tools`, `model` frontmatter, then its prompt body).
2. Add an entry to `plugins` in `.claude-plugin/marketplace.json` with
   `"source": "./plugins/<name>"`.
3. Add a row to the `## Plugins` table above, with `—` in the Release
   column — `scripts/release_skill.py` fills it in on the first release.
4. Commit and push — no build step, no publish action.

## Releasing a plugin version (skill or agent)

`scripts/release_skill.py` releases both shapes — it auto-detects which
one `<name>` is by checking for `plugins/<name>/skills/<name>/SKILL.md`
first, then `plugins/<name>/agents/<name>.md`.

Every **skill** carries its version in three places that must stay in
sync: `plugins/<name>/.claude-plugin/plugin.json`, the matching entry in
`.claude-plugin/marketplace.json`, and the `version` frontmatter field in
`plugins/<name>/skills/<name>/SKILL.md`. Every **agent** carries it in two
— `plugin.json` and `marketplace.json` only, since an agent's `.md` file
(`name`/`description`/`tools`/`model` frontmatter, then its prompt body)
has no `version` field of its own. Bump whichever apply together, then:

```bash
python3 scripts/release_skill.py <plugin-name> "What changed in this version."
```

The script verifies those version fields agree (aborting otherwise, and
naming exactly which file disagreed), refuses to re-tag a version that's
already released, zips the component directory
(`plugins/<name>/skills/<name>/` or `plugins/<name>/agents/`), publishes
it as a GitHub Release tagged `<name>-v<version>` (namespaced per plugin —
one release covers one plugin's zip, not the whole marketplace), and
rewrites that plugin's row in the Plugins table above with a link to the
new release. The README change is left unstaged for review:

```bash
git add README.md && git commit -m "docs: <name> vX.Y.Z release link"
git push origin main
```

The published release asset gives a stable download URL
(`.../releases/download/<name>-v<version>/<name>-v<version>.zip`) — useful
for anything that consumes a skill or agent by URL rather than through
`/plugin install` (e.g. importing a single one into another tool's own
skill-import feature).
