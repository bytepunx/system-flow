---
id: T-1836
type: task
nature: feature
title: The design and the user and operator docs describe the curator's triggers, digest.schedule, digest.compact_after, absorption, and compaction
status: backlog
parent: S-0432
owner: alex
created: 2026-10-16T06:07:18Z
updated: 2026-10-16T06:07:18Z
transitions: []
stream: S-0432
tags: [cli, template]
touches: [docs/operators/settings.md, design/system/project-manifest.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md, template/template.yaml]
after: [T-1834, T-1835]
---
# T-1836 The design and the user and operator docs describe the curator's triggers, digest.schedule, digest.compact_after, absorption, and compaction

## Work

- `docs/operators/settings.md` and `design/system/project-manifest.md`: `digest.schedule` and `digest.compact_after`, with their defaults.
- `design/system/strategic-agents.md`, "The curator": its triggers (after an acceptance, coalesced; the schedule), and how each run absorbs and compacts.
- `docs/users/flai.md` and `design/system/flai-cli.md`: `flai digest absorb`, `cited`, and `compact`.
- `template/CHANGELOG.md` and the version in `template/template.yaml`.

Waits for the triggers and the prompt tasks, so the docs describe what was built.

## Done when

The documents say what the code does, the settings test finds a row for each new setting, and `flai check --strict` and the markdown lint report nothing on them.

## Notes

`digest/README.md` already says nothing is deleted (S-0428); link it rather than repeat it.
