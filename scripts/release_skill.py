#!/usr/bin/env python3
"""Release a skill in this marketplace: verify plugin.json/marketplace.json/
SKILL.md agree on a version, zip the skill directory, publish it as a
GitHub Release tagged <name>-v<version>, and update that skill's row in the
README Plugins table with a link to the new release.

Usage: scripts/release_skill.py <skill-name> [release notes...]

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


def main() -> None:
    if len(sys.argv) < 2:
        fail("usage: release_skill.py <skill-name> [release notes...]")
    name = sys.argv[1]
    notes = " ".join(sys.argv[2:]) or f"Release {name}."

    repo_root = Path(__file__).resolve().parent.parent
    plugin_json_path = repo_root / "plugins" / name / ".claude-plugin" / "plugin.json"
    skill_md_path = repo_root / "plugins" / name / "skills" / name / "SKILL.md"
    marketplace_json_path = repo_root / ".claude-plugin" / "marketplace.json"
    readme_path = repo_root / "README.md"

    for p in (plugin_json_path, skill_md_path, marketplace_json_path, readme_path):
        if not p.exists():
            fail(f"missing {p.relative_to(repo_root)}")

    plugin_version = json.loads(plugin_json_path.read_text()).get("version")

    marketplace_json = json.loads(marketplace_json_path.read_text())
    entry = next((p for p in marketplace_json.get("plugins", []) if p.get("name") == name), None)
    if entry is None:
        fail(f"no marketplace.json entry for {name}")
    marketplace_version = entry.get("version")

    skill_md = skill_md_path.read_text()
    m = re.search(r'^version:\s*"([^"]+)"', skill_md, re.MULTILINE)
    if not m:
        fail(f"no version frontmatter in {skill_md_path.relative_to(repo_root)}")
    skill_version = m.group(1)

    if len({plugin_version, marketplace_version, skill_version}) != 1:
        fail(
            "version mismatch — bump all three before releasing: "
            f"plugin.json={plugin_version} marketplace.json={marketplace_version} "
            f"SKILL.md={skill_version}"
        )
    version = plugin_version
    tag = f"{name}-v{version}"

    existing = subprocess.run(
        ["gh", "release", "list", "--repo", REPO, "--json", "tagName"],
        capture_output=True, text=True, check=True,
    )
    if tag in [t["tagName"] for t in json.loads(existing.stdout)]:
        fail(f"release {tag} already exists — bump the version first")

    skill_dir = skill_md_path.parent
    zip_path = repo_root / f"{tag}.zip"
    with zipfile.ZipFile(zip_path, "w", zipfile.ZIP_DEFLATED) as zf:
        for f in sorted(skill_dir.rglob("*")):
            if f.is_file():
                zf.write(f, f.relative_to(skill_dir))

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
            "(see 'Adding a new skill')"
        )
    readme_path.write_text("".join(lines))

    print(f"Released {tag}: {release_url}")
    print("README updated — review and commit:")
    print(f"  git add README.md && git commit -m 'docs: {name} v{version} release link'")


if __name__ == "__main__":
    main()
