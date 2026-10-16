---
id: T-1761
type: task
nature: improvement
title: The design, the convention, and the changelog say the three roles are role files run through the path every role takes
status: backlog
parent: S-0419
owner: alex
created: 2026-10-16T05:46:30Z
updated: 2026-10-16T05:46:30Z
transitions: []
stream: S-0419
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [T-1754, T-1757, T-1759]
---
# T-1761 The design, the convention, and the changelog say the three roles are role files run through the path every role takes

## Work

- In `design/system/strategic-agents.md`, under each role's How it runs and What it is told, say that its prompt is its role file's body and that it starts through the defined-role path. Say what the embedded baseline does for a project without role files.
- In the convention, template and here, say that a project changes a built-in role below its file's marker.
- Name the change and `roles.baseline` in `template/CHANGELOG.md`.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
