# claude-skills-marketplace

A [Claude Code](https://claude.com/claude-code) plugin marketplace with
skills I use in my own projects.

## Install

```
/plugin marketplace add redfoxius/claude-skills-marketplace
/plugin install golang-architecture@claude-skills-marketplace
```

## Plugins

| Plugin | Description | Release |
|--------|-------------|---------|
| [golang-architecture](plugins/golang-architecture) | Forces Go-idiomatic package boundaries and dependency direction: domain/application packages never import a concrete infra package directly, ports are interfaces declared by the consumer, one composition root wires every adapter. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/golang-architecture-v1.0.0) |
| [frontend-ui-architecture](plugins/frontend-ui-architecture) | Where frontend code lives and how it's layered — folder/feature structure, component-folder anatomy, business-logic placement, types organization, barrel-file conventions for React and Next.js App Router. | [v1.0.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/frontend-ui-architecture-v1.0.0) |
| [marketplace-release](plugins/marketplace-release) | Publishes a skill to this repo — scaffolding a new skill, or running `scripts/release_skill.py` to zip, tag, publish a GitHub Release, and update this table's release link for an existing one. | [v1.1.0](https://github.com/redfoxius/claude-skills-marketplace/releases/tag/marketplace-release-v1.1.0) |

## Adding a new skill

1. Create `plugins/<name>/.claude-plugin/plugin.json` and
   `plugins/<name>/skills/<name>/SKILL.md`.
2. Add an entry to `plugins` in `.claude-plugin/marketplace.json` with
   `"source": "./plugins/<name>"`.
3. Add a row to the `## Plugins` table above, with `—` in the Release
   column — `scripts/release_skill.py` fills it in on the first release.
4. Commit and push — no build step, no publish action.

## Releasing a skill version

Every skill carries its version in three places that must stay in sync:
`plugins/<name>/.claude-plugin/plugin.json`, the matching entry in
`.claude-plugin/marketplace.json`, and the `version` frontmatter field in
`plugins/<name>/skills/<name>/SKILL.md`. Bump all three together, then:

```bash
python3 scripts/release_skill.py <skill-name> "What changed in this version."
```

The script verifies the three version fields agree (aborting otherwise),
refuses to re-tag a version that's already released, zips
`plugins/<name>/skills/<name>/`, publishes it as a GitHub Release tagged
`<name>-v<version>` (namespaced per skill — one release covers one skill's
zip, not the whole marketplace), and rewrites that skill's row in the
Plugins table above with a link to the new release. The README change is
left unstaged for review:

```bash
git add README.md && git commit -m "docs: <name> vX.Y.Z release link"
git push origin main
```

The published release asset gives a stable download URL
(`.../releases/download/<name>-v<version>/<name>-v<version>.zip`) — useful
for anything that consumes a skill by URL rather than through
`/plugin install` (e.g. importing a single skill into another tool's own
skill-import feature).
