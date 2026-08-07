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

## Adding a new skill

1. Create `plugins/<name>/.claude-plugin/plugin.json` and
   `plugins/<name>/skills/<name>/SKILL.md`.
2. Add an entry to `plugins` in `.claude-plugin/marketplace.json` with
   `"source": "./plugins/<name>"`.
3. Commit and push — no build step, no publish action.
