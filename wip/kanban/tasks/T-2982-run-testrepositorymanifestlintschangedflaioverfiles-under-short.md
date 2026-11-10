---
id: T-2982
type: task
nature: remediation
title: Run TestRepositoryManifestLintsChangedFlaioverFiles under -short
status: backlog
parent: S-0693
owner: alex
created: 2026-11-10T10:29:27Z
updated: 2026-11-10T10:29:27Z
transitions: []
stream: S-0693
tags: [flai, tests]
touches: [flai/internal/verify/select_test.go]
---
# T-2982 Run TestRepositoryManifestLintsChangedFlaioverFiles under -short

## Work

In `flai/internal/verify/select_test.go`, remove the `testing.Short()` skip ("integration: reads the monorepo") from `TestRepositoryManifestLintsChangedFlaioverFiles`. Keep its "monorepo not present" skip. Give the test a comment as its neighbours have: cite I-0327, say that it reads only the manifest and the file tree and needs no git, and that it is not gated under `-short` so that `make test` and the go-test tier catch a manifest change it must list.

Confirm first that it needs no git and runs fast: `Select` reads `os.DirFS(root)` only. If it turns out to need git or to take more than a second, stop and ask on a thread on S-0693 rather than dropping the gate.

Waits for nothing: it is the first layer.

## Done when

- `go test -short -run TestRepositoryManifestLintsChangedFlaioverFiles ./internal/verify/` in `flai/` runs the test, not skips it, and passes.
- With one line taken out of its `want` by hand, the same command fails; the line is restored after.
- `flai test flai/internal/verify/select_test.go` passes.

## Notes
