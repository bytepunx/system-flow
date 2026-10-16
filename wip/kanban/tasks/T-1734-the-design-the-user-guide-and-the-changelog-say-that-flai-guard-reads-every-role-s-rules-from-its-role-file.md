---
id: T-1734
type: task
nature: feature
title: The design, the user guide, and the changelog say that flai guard reads every role's rules from its role file
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:33Z
updated: 2026-10-16T05:44:33Z
transitions: []
stream: S-0416
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md]
after: [T-1731, T-1732, T-1733]
---
# T-1734 The design, the user guide, and the changelog say that flai guard reads every role's rules from its role file

## Work

- In `design/system/strategic-agents.md`, under § The guard, § Its permissions and the guard, and § The analyzer's guard, say that each table is what the role's file states and that the named conditions stay in flai.
- Describe the guard's role path and its fail-closed refusal in `design/system/flai-cli.md` under `flai guard`, and in `docs/users/flai.md`.
- Name the settings change and the role path in `template/CHANGELOG.md`.

Last layer: it describes the whole story.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
