---
id: T-0494
type: task
nature: feature
title: flai upgrade keeps a project's own topics on a convention and takes the template's otherwise
status: done
parent: S-0134
owner: alex
created: 2026-09-28T21:58:12Z
updated: 2026-09-28T22:04:46Z
transitions:
  - to: ready
    at: 2026-09-28T22:04:46Z
    by: agent-S-0134
  - to: in-progress
    at: 2026-09-28T22:04:46Z
    by: agent-S-0134
  - to: done
    at: 2026-09-28T22:04:46Z
    by: agent-S-0134
stream: S-0134
tags: []
---

# T-0494 flai upgrade keeps a project's own topics on a convention and takes the template's otherwise

## Work

- `system-flow.lock.yaml` records the topics the template gave each marker file (`flai new`, `flai upgrade`, `--relock`).
- The upgrade merge keeps a project's `topics` when they differ from what the lock records, or when nothing is recorded; otherwise it takes the template's.
- Docs: `design/system/conventions.md`, `template.md`, `flai-cli.md`, `project-manifest.md`, `docs/users/flai.md`, `docs/contributors/template.md`, and the conventions README in the template and here.

## Done when

- A behaviour test shows a narrowed convention keeps its topics through `flai upgrade`, and one left as the template gave it takes the template's new topics.
- `make flai-test` and `flai check --strict` pass.

## Notes

TH-0028, answer C.
