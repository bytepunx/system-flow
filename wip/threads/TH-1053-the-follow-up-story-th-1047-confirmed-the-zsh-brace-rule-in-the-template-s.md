---
id: TH-1053
title: The follow-up story TH-1047 confirmed, the zsh brace rule in the template's baseline tooling convention, does not exist yet
anchor:
  path: design/conventions/tooling.md
  heading: Rules
status: open
participants: [orchestrator]
created: 2026-11-10T19:29:35Z
updated: 2026-11-10T19:29:35Z
---

# TH-1053 The follow-up story TH-1047 confirmed, the zsh brace rule in the template's baseline tooling convention, does not exist yet

On design/conventions/tooling.md § Rules.

## Entries

### 2026-11-10T19:29:35Z orchestrator
**Recommended:** create the follow-up story now, for example:

```sh
flai story new --epic <epic> "The baseline zsh rule says to write \${name} before a colon, since quoting does not stop zsh's history modifiers"
```

- On TH-1047 you confirmed answer 1: the rule goes into `template/root/design/conventions/tooling.md` through a follow-up story, then gets copied here above the marker, with a `template/CHANGELOG.md` line (D-0578).
- No story in `wip/kanban/stories` covers it yet. Creating a story is yours, and I cannot do it.
- Until it lands, this repository is covered by S-0667's project addition below the marker.

**Alternative:** tell agent-S-0667, or another story's agent, to create it.
