---
id: T-1836
type: task
nature: feature
title: The design and the user and operator docs describe the curator's triggers, digest.schedule, digest.compact_after, absorption, and compaction
status: done
parent: S-0432
owner: alex
created: 2026-10-16T06:07:18Z
updated: 2026-10-24T14:46:48Z
transitions:
  - to: ready
    at: 2026-10-24T14:31:28Z
    by: agent-S-0432
  - to: in-progress
    at: 2026-10-24T14:31:29Z
    by: agent-S-0432
  - to: done
    at: 2026-10-24T14:46:48Z
    by: agent-S-0432
stream: S-0432
tags: [cli, template]
touches: [docs/operators/settings.md, docs/operators/index.md, design/system/project-manifest.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/guard.go, template/CHANGELOG.md, template/template.yaml, design/system/agent-adapters.md]
after: [T-1834]
usage:
  source: log
  seconds: 919
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 99
      output: 39457
      cache_read: 6511058
      cache_write: 269922
      cost: 3.6069
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
