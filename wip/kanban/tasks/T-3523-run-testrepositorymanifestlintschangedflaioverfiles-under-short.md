---
id: T-3523
type: task
nature: remediation
title: Run TestRepositoryManifestLintsChangedFlaioverFiles under -short
status: in-progress
parent: S-0827
owner: alex
created: 2026-11-17T08:28:22Z
updated: 2026-11-17T08:31:29Z
transitions:
  - to: ready
    at: 2026-11-17T08:31:29Z
    by: agent-S-0827
  - to: in-progress
    at: 2026-11-17T08:31:29Z
    by: agent-S-0827
stream: S-0827
tags: [flai, tests]
touches: [flai/internal/verify/select_test.go]
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 9
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 3
          output: 50
          cache_read: 217510
          cache_write: 2537
          cost: 0.0682
---
# T-3523 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short

## Work

In `flai/internal/verify/select_test.go`, remove the `testing.Short()` skip ("integration: reads the monorepo") from `TestRepositoryManifestLintsChangedFlaioverFiles`. Keep its "monorepo not present" skip. Give the test a comment as its neighbours have: cite I-0327, say that it reads only the manifest and the file tree through `os.DirFS` and needs no git, and that it is not gated under `-short` so that `make test` and the go-test tier catch a manifest change it must list.

This is S-0693's fix again: S-0693 made it and was cancelled before it landed, because its branch undid other work on main. Start from main's file, not from that branch.

Confirm first that the test needs no git and runs in under a second. If it needs git or takes longer, stop and ask on a thread on S-0827 rather than dropping the gate. If, run without the gate, it fails because the manifest now selects a tier its `want` lacks, or no longer selects one it lists, set `want` to what `system-flow.yaml` selects today and say so in the task's notes: the test pins the manifest, not the other way round.

Waits for nothing: it is the first layer.

## Done when

- `go test -short -run TestRepositoryManifestLintsChangedFlaioverFiles ./internal/verify/` in `flai/` runs the test, not skips it, and passes.
- With one line taken out of its `want` by hand, the same command fails; the line is restored after.
- `flai test flai/internal/verify/select_test.go` passes.

## Notes
