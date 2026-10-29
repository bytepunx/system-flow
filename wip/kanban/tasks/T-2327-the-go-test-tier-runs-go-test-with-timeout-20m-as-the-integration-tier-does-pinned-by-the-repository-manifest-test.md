---
id: T-2327
type: task
nature: remediation
title: The go-test tier runs go test with -timeout 20m, as the integration tier does, pinned by the repository-manifest test
status: backlog
parent: S-0540
owner: alex
created: 2026-10-29T07:37:05Z
updated: 2026-10-29T07:37:05Z
transitions: []
stream: S-0540
tags: [flai, testing]
touches: [system-flow.yaml, flai/internal/verify/select_test.go, design/conventions/tooling.md]
---
# T-2327 The go-test tier runs go test with -timeout 20m, as the integration tier does, pinned by the repository-manifest test

## Work

The go-test tier in `system-flow.yaml` passes no `-timeout`, so go test's 10-minute default applies to its `-race` run. `flai/internal/mcpserver` takes 9m20s to 9m58s under `-race` on a loaded host, and 10m1s to 10m21s at worst, so it fails on time alone. `scripts/integration.sh` passes `-timeout 20m` for the same package since I-0225. This is I-0227's first cause.

- Add `"-timeout", "20m"` to the go-test tier's `command`, after `-count=1`. Write the tier with `flai manifest set`, as S-0490 did, rather than by hand: `CLAUDE.md` gives the manifest's other keys to flai.
- Update the test in `flai/internal/verify/select_test.go` that S-0490 added to pin this repository's tier arguments, so that it pins the timeout too.
- In `design/conventions/tooling.md`, below the marker, in the project addition on `flai test`, say that the go-test tier runs with a 20-minute timeout a package, as the integration tier does.

It waits for no task: no other task of this story touches these files.

## Done when

- `system-flow.yaml`'s go-test tier passes `-timeout 20m`, and `flai check --strict` passes.
- The pinning test in `flai/internal/verify/select_test.go` fails without the timeout and passes with it.
- `flai test flai/internal/verify` passes.

## Notes
