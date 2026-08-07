---
name: marketplace-release
description: "Publishes a skill to the redfoxius/claude-skills-marketplace Claude Code plugin repo: scaffolding plugins/<name>/.claude-plugin/plugin.json + skills/<name>/SKILL.md + a matching entry in .claude-plugin/marketplace.json for a brand-new skill, or running scripts/release_skill.py to package, tag, and publish a GitHub Release for a version bump of an existing skill. That script enforces plugin.json/marketplace.json/SKILL.md version agreement, refuses to re-tag an existing release, and rewrites the skill's README table row with a link to the new release — nothing about a release is a manual doc-editing step. Use when asked to publish, release, or version-bump a skill in this specific marketplace repo — not for general Claude Code skill authoring (structuring a SKILL.md's content) and not for consuming a plugin (see this repo's README Install section)."
version: "1.1.0"
---

# Marketplace Release

Two workflows in `redfoxius/claude-skills-marketplace`: adding a brand-new
skill, and releasing a new version of one that already exists.

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
4. Add a row to the `## Plugins` table in the repo root `README.md`, with
   `—` in the Release column — `scripts/release_skill.py` fills that cell
   in on the first release, matching the row by
   `| [<name>](plugins/<name>)` at line start.
5. Commit and push — no build step, no publish action for this step.

## Releasing a version (new or bumped)

Bump `version` in `plugin.json`, the `marketplace.json` entry, and the
`SKILL.md` frontmatter together — the release script refuses to run if
they disagree, so there's no "forgot one of the three" failure mode. Then:

```bash
python3 scripts/release_skill.py <skill-name> "What changed in this version."
```

This single command does everything a release needs:

- **Preflight** — reads the version out of `plugin.json`, the matching
  `marketplace.json` entry, and `SKILL.md`'s frontmatter; aborts if they
  don't all match, and aborts if `<name>-v<version>` is already tagged
  (no silent re-release).
- **Package** — zips `plugins/<name>/skills/<name>/` in full (whatever
  files it contains — `SKILL.md` alone or with `examples.md`/`README.md`)
  via Python's `zipfile`, not a shell `cd`/relative-path recipe — the kind
  of thing that silently landed a zip one directory off in an earlier
  version of this process.
- **Publish** — `gh release create <name>-v<version>` with that zip as the
  asset, tag namespaced per skill since one release covers one skill's zip,
  not the whole marketplace.
- **Update docs** — rewrites the Release cell of that skill's row in the
  README Plugins table to `[vX.Y.Z](<release-url>)`. This is the step that
  used to be a separate manual edit, easy to forget right after tagging a
  release — now it's not a separate step at all.

The script leaves the README diff **unstaged** — review it, then:

```bash
git add README.md && git commit -m "docs: <name> vX.Y.Z release link"
git push origin main
```

The published release asset gives a stable download URL
(`.../releases/download/<name>-v<version>/<name>-v<version>.zip`) for
anything that consumes a skill by URL rather than through
`/plugin install` — e.g. importing a single skill into another tool's own
skill-import feature.
