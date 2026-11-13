---
id: T-3198
type: task
nature: feature
title: The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off
status: done
parent: S-0759
owner: alex
created: 2026-11-12T16:05:58Z
updated: 2026-11-13T17:31:39Z
transitions:
  - to: ready
    at: 2026-11-13T17:19:36Z
    by: system-flow
  - to: in-progress
    at: 2026-11-13T17:19:36Z
    by: system-flow
  - to: done
    at: 2026-11-13T17:31:39Z
    by: agent-S-0759
stream: S-0759
tags: [cli]
touches: [docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/runbooks/agent-container.md, design/system/flai-cli.md, flai/internal/check/image.go, flai/internal/check/image_test.go]
after: [T-3191, T-3194]
usage:
  source: log
  seconds: 723
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 101
      output: 38676
      cache_read: 7675419
      cache_write: 277479
      cost: 3.9374
---
# T-3198 The settings index, the guides, the runbook, and flai-cli.md describe the operator's build command and the switch that turns the automatic image job off

## Work

- `docs/operators/settings.md`: a row under § flai configuration for `agent.projects.<root>.image.command` and `agent.projects.<root>.image.auto`, near `agent.projects`, and the commands' flags under § Flags. `TestSettingsIndex` fails on a key or flag with no row.
- `docs/users/flai.md` § In a container: both settings, `flai serve agent image command` and `flai serve agent image auto`.
- `docs/users/flai-reference.md`: `make flai-reference`, never by hand.
- `docs/operators/runbooks/agent-container.md`: a section, "Your own images", on naming a build command, its `FLAI_IMAGE_*` variables, its log, and turning the automatic job off, beside `--image`.
- `design/system/flai-cli.md`: mark § Where the agent images come from (ADR-0295) › The operator's own machinery built (S-0759), and bump `updated`.

It waits for T-3191 and T-3194, whose commands and behaviour it describes.

## Done when

- The five files describe both settings as built.
- `flai test` on the five files passes: `TestSettingsIndex`, `TestFlagIndexIsCurrent`, `TestReferenceIsCurrent`, and the markdown lint.

## Notes
