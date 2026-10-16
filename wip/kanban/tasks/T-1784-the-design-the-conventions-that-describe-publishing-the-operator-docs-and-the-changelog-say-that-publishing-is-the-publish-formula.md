---
id: T-1784
type: task
nature: improvement
title: The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula
status: backlog
parent: S-0423
owner: alex
created: 2026-10-16T05:48:51Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0423
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, docs/operators/index.md, template/CHANGELOG.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1782]
---
# T-1784 The design, the conventions that describe publishing, the operator docs, and the changelog say that publishing is the publish formula

## Work

- In `design/system/strategic-agents.md` § Publishing a release, say that under `publish` the orchestrator runs the `publish` formula, and keep `release_publish` for a project with none.
- Where `work-management.md` and `strategic-agents.md`, template and here, describe publishing, name the formula. A baseline rule is changed in the template's baseline and in this repository's copy alike, as the template's release would bring it.
- Add publishing through the formula to `docs/operators/index.md`, and name it in `template/CHANGELOG.md`.

Waits for the orchestrator's rule. Runs beside the run, since they share no path.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
