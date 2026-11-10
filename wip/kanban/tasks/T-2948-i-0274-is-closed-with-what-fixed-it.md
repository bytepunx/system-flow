---
id: T-2948
type: task
nature: improvement
title: I-0274 is closed with what fixed it
status: backlog
parent: S-0665
owner: alex
created: 2026-11-10T08:50:50Z
updated: 2026-11-10T08:50:56Z
transitions: []
stream: S-0665
tags: [issues]
touches: [design/issues/I-0274-flai-serve-records-no-image-container-name-or-mounts-for-a-container-harness.md, design/issues/summary.md]
after: [T-2943, T-2945]
---
# T-2948 I-0274 is closed with what fixed it

## Work

Close I-0274 with `flai issue close I-0274 --reason "<reason>"`. The reason names the fields `agents.json` now records for a container run (`image`, `container`, `mounts`), the journal's start line, and the test in `flai/internal/serve/container_test.go` that reproduces the issue. The command rewrites `design/issues/summary.md`.

It waits for T-2943 and T-2945, the fix its reason names. It runs beside T-2946, the docs, which shares no path with it.

## Done when

- I-0274's status is closed, with a reason naming the fix and its test.
- `design/issues/summary.md` no longer lists I-0274 as open.
- `flai check --strict` passes.

## Notes
