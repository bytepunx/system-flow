---
id: T-1740
type: task
nature: feature
title: The design, the manifest's design, and the user and operator docs say how a defined role is started, as whom, and where its activity is logged
status: in-progress
parent: S-0417
owner: alex
created: 2026-10-16T05:45:15Z
updated: 2026-10-29T05:14:19Z
transitions:
  - to: ready
    at: 2026-10-29T05:14:17Z
    by: agent-S-0417
  - to: in-progress
    at: 2026-10-29T05:14:19Z
    by: agent-S-0417
stream: S-0417
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/project-manifest.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/agent-narrative.md, docs/users/flai-reference.md]
after: [T-1739]
---
# T-1740 The design, the manifest's design, and the user and operator docs say how a defined role is started, as whom, and where its activity is logged

## Work

- Add a section to `design/system/strategic-agents.md` on starting a defined role: the entry points, the refusals, how it runs, its record, and its activity document.
- Add `roles.<name>.agent` and `roles.<name>.permissions` to `design/system/project-manifest.md`, and add their rows to `docs/operators/settings.md`. The settings test fails on a key with no row.
- Describe `flai role run` and `role_run` in `docs/users/flai.md`.
- Add how to start a defined role to `docs/operators/index.md`.
- Name the commands in `template/CHANGELOG.md`.

Last layer.

## Done when

The markdown lint, the settings-index test, and `flai check --strict` pass.

## Notes
