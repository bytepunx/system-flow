---
id: T-0812
type: task
nature: feature
title: flai release --evaluate says whether the release policy is met, with its figures
status: done
parent: S-0217
owner: alex
created: 2026-10-04T19:09:11Z
updated: 2026-10-05T06:14:35Z
transitions:
  - to: ready
    at: 2026-10-05T06:06:12Z
    by: agent-S-0217
  - to: in-progress
    at: 2026-10-05T06:06:12Z
    by: agent-S-0217
  - to: done
    at: 2026-10-05T06:14:35Z
    by: agent-S-0217
stream: S-0217
tags: [flai]
touches: [flai/internal/release/evaluate.go, flai/internal/release/evaluate_test.go, flai/cmd/release.go, flai/cmd/release_evaluate_test.go]
after: [T-0809]
usage:
  source: log
  seconds: 503
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 518
      cache_read: 2632221
      cache_write: 103698
      cost: 2.3572
---
# T-0812 flai release --evaluate says whether the release policy is met, with its figures

## Work

Add `flai/internal/release/evaluate.go`. It reads `orchestration.release` and returns the policy, whether it is met, its figures, and a sentence saying why:

- `threshold`: take the stories accepted and not yet released from `PendingIDs()` (`release.go`, line 655). Sum their cost of delay values and count them. It is met when either figure is at or over the manifest's. A story with no value counts toward `count` but adds nothing to the sum, and is named.
- `theme`: it is met when every story of the named epic, or every story with the named tag, is accepted, and at least one of them is unreleased. List the ones not yet accepted.
- `judgement`: it is never met by itself. Say that the call is the orchestrator's or the operator's.

Give `flai release` (`flai/cmd/release.go`) `--evaluate [--json]`. It refuses to run with `--apply`, `--pending`, or `--dry-run`, and tags, bumps, and pushes nothing.

This task waits for T-0809, for the `orchestration.release` setting. It runs with T-0811, whose paths it does not share.

## Done when

- fixture tests pin `threshold` met and not met on value and on count, `theme` by epic and by tag, met and not, and `judgement`
- `flai release --evaluate` prints the policy, the verdict, and the figures, and changes no file, tag, or remote
- `go test ./internal/release/ ./cmd/` passes

## Notes
