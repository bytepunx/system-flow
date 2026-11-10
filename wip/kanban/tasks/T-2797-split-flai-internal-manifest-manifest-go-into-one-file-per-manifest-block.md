---
id: T-2797
type: task
nature: improvement
title: Split flai/internal/manifest/manifest.go into one file per manifest block
status: done
parent: S-0631
owner: alex
created: 2026-11-09T22:04:48Z
updated: 2026-11-10T06:38:47Z
transitions:
  - to: ready
    at: 2026-11-10T06:16:47Z
    by: agent-S-0631
  - to: in-progress
    at: 2026-11-10T06:16:48Z
    by: agent-S-0631
  - to: done
    at: 2026-11-10T06:38:47Z
    by: agent-S-0631
stream: S-0631
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_strategic.go, flai/internal/manifest/manifest_roles.go, flai/internal/manifest/manifest_accept.go, flai/internal/manifest/manifest_kluster.go, flai/internal/manifest/manifest_tests.go]
usage:
  source: log
  seconds: 1281
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 19
      output: 154
      cache_read: 2728664
      cache_write: 30935
      cost: 1.0149
---
# T-2797 Split flai/internal/manifest/manifest.go into one file per manifest block

## Work

`manifest.go` is 2,449 lines. Move each block of the manifest's types and validation into its own file of package `manifest`. What any of them does stays the same. The line ranges below are approximate, from the plan's read on 2026-11-09:

| Lines | Declarations | File |
|-------|--------------|------|
| 1–208, 2061–2449 | `Manifest`, `Requirement`, `Design`, `Load`, `Find`, `Dir`, `Layer`, the project entries (`Project`, `Install`), `Components`, `Dashboard` | `manifest.go`, which keeps them |
| 209–645 | `Prime`, `Issues`, `Planning`, `Analysis`, `Digest`, `Threads`, `Orchestration` | `manifest_strategic.go` |
| 646–963 | `Permissions`, `RoleSettings`, `RoleErrors`, `RoleSchedule`, `RoleAgent`, `OrchestrationAgent` | `manifest_roles.go` |
| 964–1301 | `Release`, `Acceptance`, `MergeQueue`, `BatchMax` | `manifest_accept.go` |
| 1302–1667 | `Kluster`, `KlusterComponent`, `DNSLabel`, `ReadKlusterConfig`, `KlusterErrors` | `manifest_kluster.go` |
| 1668–2060 | `NamedCommand`, `TestTier`, `TestTiers`, `TestErrors`, `TestWarnings`, `coverCycles`, `commandProblems` | `manifest_tests.go` |

`projects.go` (512 lines) could take the project entries instead. Choose the place that reads better, and keep every file under 1,200.

The manifest's settings catalog is in `settings.go`, which this task leaves alone, so `docs/operators/settings.md` and `docs/users/flai-reference.md` do not change.

Change no test file. `manifest_test.go` stays as it is.

This task waits for nothing and shares no path with the other splits, so all of them run together.

## Done when

- No non-test file in `flai/internal/manifest` is over 2,000 lines, and none of the new files is over 1,200.
- `git diff --stat` shows no `_test.go` file changed.
- `flai test flai/internal/manifest` passes, gofmt, goimports, and golangci-lint included.

## Notes
