---
name: marketplace-release
description: "Publishes a skill to the redfoxius/claude-skills-marketplace Claude Code plugin repo: scaffolding plugins/<name>/.claude-plugin/plugin.json + skills/<name>/SKILL.md + a matching entry in .claude-plugin/marketplace.json for a brand-new skill, or packaging and tagging a GitHub Release (versioned zip, <name>-v<version> tag) for a version bump of an existing skill. Enforces that plugin.json, marketplace.json, and SKILL.md's version frontmatter stay numerically identical before any release is tagged. Use when asked to publish, release, or version-bump a skill in this specific marketplace repo — not for general Claude Code skill authoring (structuring a SKILL.md's content) and not for consuming a plugin (see this repo's README Install section)."
version: "1.0.1"
---

# Marketplace Release

Two workflows in `redfoxius/claude-skills-marketplace`: adding a brand-new
skill, and releasing a new version of one that already exists. Both end
with a `git push` to `origin/main`; a version release additionally tags
and publishes a GitHub Release.

## Adding a brand-new skill

1. `plugins/<name>/.claude-plugin/plugin.json` — `name`, `displayName`,
   `version: "1.0.0"`, `description`, `author`, `homepage`, `repository`,
   `license`, `keywords`. Copy an existing plugin's `plugin.json` as the
   template rather than writing one from scratch.
2. `plugins/<name>/skills/<name>/SKILL.md` — frontmatter `name`,
   `description` (the single-line trigger Claude Code matches against —
   front-load what the skill governs and when to use it, and name what it
   deliberately excludes), `version: "1.0.0"` matching step 1.
3. Add an entry to the `plugins` array in `.claude-plugin/marketplace.json`
   with `"source": "./plugins/<name>"`, mirroring `plugin.json`'s
   `description`, `version`, and metadata.
4. Add a row to the `## Plugins` table in the repo root `README.md`.
5. Commit and push — no build step, no publish action for this step.

## Releasing a version (new or bumped)

**Preflight — the three version fields must match exactly before tagging:**
`plugins/<name>/.claude-plugin/plugin.json` → `version`,
`.claude-plugin/marketplace.json` → that plugin's `version` entry,
`plugins/<name>/skills/<name>/SKILL.md` → frontmatter `version`. If you
bumped one for this release, bump all three in the same commit — a
mismatch across these files is a bug, not three independent numbers.

```bash
name=<skill-name>
version=<x.y.z>   # must equal all three fields checked above

(cd "plugins/$name/skills/$name" && zip -r - .) > "${name}-v${version}.zip"

gh release create "${name}-v${version}" "${name}-v${version}.zip" \
  --title "${name} v${version}" \
  --notes "What changed in this version."
rm "${name}-v${version}.zip"
```

Tags are namespaced per skill (`<name>-v<version>`), not per repo — one
release covers one skill's zip, since the marketplace hosts several. The
zip archives the skill's whole directory (`SKILL.md` plus any
`examples.md`, `README.md`, or other files it ships) via `zip -r ... .`
from inside `skills/<name>/`, so it stays correct whether the skill is a
single file or several — don't hardcode filenames into the zip command.

The published release asset gives a stable download URL
(`.../releases/download/<name>-v<version>/<name>-v<version>.zip`) for
anything that consumes a skill by URL rather than through
`/plugin install` — e.g. importing a single skill into another tool's own
skill-import feature.
