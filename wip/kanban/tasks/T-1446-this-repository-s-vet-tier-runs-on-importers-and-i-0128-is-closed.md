---
id: T-1446
type: task
nature: improvement
title: "This repository's vet tier runs on {importers}, and I-0128 is closed"
status: done
parent: S-0364
owner: alex
created: 2026-10-12T14:40:50Z
updated: 2026-10-12T15:36:58Z
transitions:
  - to: ready
    at: 2026-10-12T15:35:38Z
    by: agent-S-0364
  - to: in-progress
    at: 2026-10-12T15:35:38Z
    by: agent-S-0364
  - to: done
    at: 2026-10-12T15:36:58Z
    by: agent-S-0364
stream: S-0364
tags: [cli]
touches: [system-flow.yaml, design/issues/I-0128-flai-test-vets-and-tests-only-the-changed-go-packages-so-a-change-that-breaks-a-package-importing-them-is-found-only-by-the-integration-tier.md, design/issues/summary.md]
after: [T-1444]
usage:
  source: log
  seconds: 80
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 31
      output: 9608
      cache_read: 1430850
      cache_write: 57209
      cost: 0.8502
---
# T-1446 This repository's vet tier runs on {importers}, and I-0128 is closed

## Work

In `system-flow.yaml`, change the `vet` tier's command from `{packages}` to `{importers}`. `go vet` type-checks every package it names, test files included, so a change that breaks an importing package now fails `flai test` and the close-out's `vet` tier. Leave `go-test` and `golangci-lint` on `{packages}`.

Check it on this tree with `flai test flai/internal/workitem --json`. The `vet` tier's command must name `./cmd` among its packages. Record the tier's duration in the story's notes. On this host `go vet ./...` over the module took 2.0s warm and 11.8s with an empty build cache.

Then close the issue from the worktree with `flai issue close I-0128 --reason "<what fixed it>"`. Name `{importers}` and the `vet` tier in the reason. The command updates `design/issues/summary.md`.

This task waits for T-1444: the manifest names `{importers}` only once flai fills it, and `flai check` refuses it before T-1443.

## Done when

- The `vet` tier in `system-flow.yaml` runs `go vet "{importers}"`.
- `flai test flai/internal/workitem` runs `vet` on `./cmd` too, and passes.
- I-0128 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it as open.

## Notes

Drafted by planner-S-0364.
