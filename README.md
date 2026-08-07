# claude-skills-marketplace

A [Claude Code](https://claude.com/claude-code) plugin marketplace with
skills I use in my own projects.

## Install

```
/plugin marketplace add redfoxius/claude-skills-marketplace
/plugin install golang-architecture@claude-skills-marketplace
```

## Plugins

| Plugin | Description |
|--------|-------------|
| [golang-architecture](plugins/golang-architecture) | Forces Go-idiomatic package boundaries and dependency direction: domain/application packages never import a concrete infra package directly, ports are interfaces declared by the consumer, one composition root wires every adapter. |
| [frontend-ui-architecture](plugins/frontend-ui-architecture) | Where frontend code lives and how it's layered — folder/feature structure, component-folder anatomy, business-logic placement, types organization, barrel-file conventions for React and Next.js App Router. |
| [marketplace-release](plugins/marketplace-release) | Publishes a skill to this repo — scaffolding a new skill or packaging and tagging a versioned GitHub Release for an existing one, keeping `plugin.json`/`marketplace.json`/`SKILL.md` version fields in sync. |

## Adding a new skill

1. Create `plugins/<name>/.claude-plugin/plugin.json` and
   `plugins/<name>/skills/<name>/SKILL.md`.
2. Add an entry to `plugins` in `.claude-plugin/marketplace.json` with
   `"source": "./plugins/<name>"`.
3. Commit and push — no build step, no publish action.

## Releasing a skill version

Every skill carries its version in three places that must stay in sync:
`plugins/<name>/.claude-plugin/plugin.json`, the matching entry in
`.claude-plugin/marketplace.json`, and the `version` frontmatter field in
`plugins/<name>/skills/<name>/SKILL.md`.

A GitHub Release publishes that version as a standalone, versioned zip
archive of the skill's contents — useful for anything that consumes a
skill by URL rather than through `/plugin install` (e.g. importing a
single skill into another tool's own skill-import feature). Tags are
namespaced per skill (`<name>-v<version>`) since one release covers one
skill, not the whole marketplace.

```bash
name=golang-architecture
version=1.0.0

(cd "plugins/$name/skills/$name" && zip -r - .) > "${name}-v${version}.zip"

gh release create "${name}-v${version}" "${name}-v${version}.zip" \
  --title "${name} v${version}" \
  --notes "Release notes for this version."
rm "${name}-v${version}.zip"
```

Bump all three version fields together before tagging a new release —
mismatched versions across `plugin.json`, `marketplace.json`, and
`SKILL.md` are a bug, not three independent version numbers.
