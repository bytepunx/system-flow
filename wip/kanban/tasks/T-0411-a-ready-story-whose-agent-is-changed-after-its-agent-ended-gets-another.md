---
id: T-0411
type: task
nature: remediation
title: A ready story whose agent is changed after its agent ended gets another
status: in-progress
parent: S-0116
owner: alex
created: 2026-09-24T08:46:27Z
updated: 2026-09-24T08:46:44Z
transitions:
  - to: ready
    at: 2026-09-24T08:46:44Z
    by: agent-S-0116
  - to: in-progress
    at: 2026-09-24T08:46:44Z
    by: agent-S-0116
stream: S-0116
tags: []
touches: [flai/internal/serve, flai/internal/manifest, flai/internal/itemedit]
---
# T-0411 A ready story whose agent is changed after its agent ended gets another

## Work

The launcher gives a story one agent per entry into ready (ADR-0041): a story whose agent failed waits until it is moved to ready again. Record on each run the story's agent it was started with, and treat a ready story whose agent now says something else, with no agent running, as never started since it entered ready. Move the comparison `sameAgent` from `internal/itemedit` to `manifest.Agent` so both use one. A run recorded before this change, with no agent on it, counts as unchanged.

## Done when

- A behavior test in `flai/internal/serve` shows: a story whose agent failed is not started again on the next look; after its agent (harness, model, or a config value) is changed while it stays in ready, the next look starts it, within the in-progress limit; a change while its agent runs starts nothing.
- The reason logged for a story that had its agent says that changing its agent also starts another.
- `make test` and lint pass.

## Notes
