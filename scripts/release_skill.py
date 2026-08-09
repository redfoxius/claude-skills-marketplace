#!/usr/bin/env python3
"""Release a plugin (skill or agent) in this marketplace: verify
plugin.json/marketplace.json agree on a version (and SKILL.md's frontmatter
too, for skills — agent files carry no version of their own), zip the
component directory, publish it as a GitHub Release tagged <name>-v<version>,
and update that plugin's row in the README Plugins table with a link to the
new release.

Usage: scripts/release_skill.py <plugin-name> [release notes...]

Leaves the README change unstaged — review the diff, then:
  git add README.md && git commit -m "docs: <name> vX.Y.Z release link"
"""
import json
import re
import subprocess
import sys
import zipfile
from pathlib import Path

REPO = "redfoxius/claude-skills-marketplace"


def fail(msg: str) -> None:
    print(f"error: {msg}", file=sys.stderr)
    sys.exit(1)


def find_component(repo_root: Path, name: str):
    """Locate the plugin's skill or agent file and return
    (component_type, component_dir, component_version). Skills carry their
    own `version` frontmatter in SKILL.md; agent files don't define a
    version at all, so component_version is None for agents — the version
    check falls back to plugin.json vs. marketplace.json agreement only."""
    skill_md = repo_root / "plugins" / name / "skills" / name / "SKILL.md"
    agent_md = repo_root / "plugins" / name / "agents" / f"{name}.md"

    if skill_md.exists():
        text = skill_md.read_text()
        m = re.search(r'^version:\s*"([^"]+)"', text, re.MULTILINE)
        if not m:
            fail(f"no version frontmatter in {skill_md.relative_to(repo_root)}")
        return "skill", skill_md.parent, m.group(1)

    if agent_md.exists():
        return "agent", agent_md.parent, None

    fail(
        f"no plugins/{name}/skills/{name}/SKILL.md or "
        f"plugins/{name}/agents/{name}.md found for '{name}'"
    )


def main() -> None:
    if len(sys.argv) < 2:
        fail("usage: release_skill.py <plugin-name> [release notes...]")
    name = sys.argv[1]
    notes = " ".join(sys.argv[2:]) or f"Release {name}."

    repo_root = Path(__file__).resolve().parent.parent
    plugin_json_path = repo_root / "plugins" / name / ".claude-plugin" / "plugin.json"
    marketplace_json_path = repo_root / ".claude-plugin" / "marketplace.json"
    readme_path = repo_root / "README.md"

    for p in (plugin_json_path, marketplace_json_path, readme_path):
        if not p.exists():
            fail(f"missing {p.relative_to(repo_root)}")

    component_type, component_dir, component_version = find_component(repo_root, name)

    plugin_version = json.loads(plugin_json_path.read_text()).get("version")

    marketplace_json = json.loads(marketplace_json_path.read_text())
    entry = next((p for p in marketplace_json.get("plugins", []) if p.get("name") == name), None)
    if entry is None:
        fail(f"no marketplace.json entry for {name}")
    marketplace_version = entry.get("version")

    versions = {"plugin.json": plugin_version, "marketplace.json": marketplace_version}
    if component_version is not None:
        versions["SKILL.md"] = component_version

    if len(set(versions.values())) != 1:
        mismatch = ", ".join(f"{k}={v}" for k, v in versions.items())
        fail(f"version mismatch — bump all together before releasing: {mismatch}")
    version = plugin_version
    tag = f"{name}-v{version}"

    existing = subprocess.run(
        ["gh", "release", "list", "--repo", REPO, "--json", "tagName"],
        capture_output=True, text=True, check=True,
    )
    if tag in [t["tagName"] for t in json.loads(existing.stdout)]:
        fail(f"release {tag} already exists — bump the version first")

    zip_path = repo_root / f"{tag}.zip"
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zf:
        for f in sorted(component_dir.rglob("*")):
            if f.is_file():
                zf.write(f, f.relative_to(component_dir))

    try:
        subprocess.run(
            ["gh", "release", "create", tag, str(zip_path),
             "--title", f"{name} v{version}", "--notes", notes, "--repo", REPO],
            check=True,
        )
    finally:
        zip_path.unlink(missing_ok=True)

    release_url = f"https://github.com/{REPO}/releases/tag/{tag}"

    readme = readme_path.read_text()
    row_prefix = f"| [{name}](plugins/{name})"
    lines = readme.splitlines(keepends=True)
    updated = False
    for i, line in enumerate(lines):
        if line.startswith(row_prefix):
            cells = line.rstrip("\n").split("|")
            if len(cells) < 3:
                fail(f"README row for {name} doesn't look like a 3-column table row")
            cells[-2] = f" [v{version}]({release_url}) "
            lines[i] = "|".join(cells) + "\n"
            updated = True
            break
    if not updated:
        fail(
            f"no README table row found for {name} — add one first "
            "(see 'Adding a new skill' / 'Adding a new agent')"
        )
    readme_path.write_text("".join(lines))

    print(f"Released {tag} ({component_type}): {release_url}")
    print("README updated — review and commit:")
    print(f"  git add README.md && git commit -m 'docs: {name} v{version} release link'")


if __name__ == "__main__":
    main()
