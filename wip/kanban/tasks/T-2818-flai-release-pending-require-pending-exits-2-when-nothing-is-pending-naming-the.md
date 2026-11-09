---
id: T-2818
type: task
nature: improvement
title: flai release --pending --require-pending exits 2 when nothing is pending, naming the release that covered it
status: backlog
parent: S-0643
owner: agent
created: 2026-11-09T22:12:26Z
updated: 2026-11-09T22:12:26Z
transitions: []
stream: S-0643
tags: [flai, release]
touches: [flai/cmd/release.go, flai/cmd/release_pending_test.go, flai/internal/release/release.go, flai/internal/release/release_test.go]
---
# T-2818 flai release --pending --require-pending exits 2 when nothing is pending, naming the release that covered it

## Work

Add `--require-pending` to `flai release --pending` in `flai/cmd/release.go`, as `--require-met` is to `--evaluate`. With nothing pending, the command prints what it prints today and then exits 2, reusing `exitReleaseNotMet`. Its output, text and `--json`, names each component's newest release tag, the commit the tag points at, its tagger (or the commit's author for a lightweight tag), and its time. That is the release that covered what was pending. Refuse the flag without `--pending`, as `--require-met` is refused without `--evaluate`. Decide whether `--dry-run` with nothing pending exits 2 too, and say why in the narrative. The recommended answer is yes, so a preview matches the step.

Add the lookup of a component's newest release tag with its commit, tagger, and time to `flai/internal/release/release.go`, beside `CurrentVersion`, with a unit test in `flai/internal/release/release_test.go`.

Reproduce I-0320 in `flai/cmd/release_pending_test.go`, beside `TestReleasePendingNothingPending`. Publish once, as the operator did by hand at 16:57:57Z. Then run `--pending --require-pending --json` as the instance's step did. It must exit 2 and name the first publish's tag and commit. Plain `--pending` must still exit 0 with `"reason": "nothing pending"`. This task waits for none: it is the first layer.

## Done when

- `flai release --pending --require-pending` exits 2 with nothing pending, and names each component's newest release tag, commit, tagger, and time, in text and in `--json`.
- Plain `flai release --pending` behaves as before, and the dashboard's Publish and `release_publish` are unchanged.
- The new test reproduces I-0320's sequence and passes, and `flai test flai/cmd/release.go flai/internal/release/release.go` passes.

## Notes
