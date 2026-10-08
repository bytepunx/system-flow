---
id: T-1355
type: task
nature: improvement
title: "flai test and flai verify list a covered tier as skipped: covered by its tier, in the text and in --json"
status: backlog
parent: S-0342
owner: alex
created: 2026-10-08T08:06:11Z
updated: 2026-10-08T08:06:28Z
transitions: []
stream: S-0342
tags: [cli]
touches: [flai/internal/verify/text.go, flai/internal/verify/text_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go]
after: [T-1349]
---
# T-1355 flai test and flai verify list a covered tier as skipped: covered by its tier, in the text and in --json

## Work

- `Result.Text` in `text.go` prints a skipped tier as `skipped <name>: covered by <tier>`, with no duration; `text_test.go` shows the line among a passed and a failed tier.
- `verify.Step` in `story.go` gains `CoveredBy string` with `json:"covered_by,omitempty"`, and `tierStep` copies it from the `TierResult`, so the verify record and `flai verify --json` and `--last` carry `state: skipped` and `covered_by`.
- `verifyText` in `cmd/verify.go` prints a skipped step as `skipped <name>: covered by <tier>`, as `flai test` does, with no duration; check the other loop over `rep.Steps` in that file (near line 298) treats `skipped` as neither a failure nor a step not reached. `cmd/verify_test.go` shows the line.
- When a step before the tiers fails, a covered tier is still listed as `skipped` with its `covered_by`, not `not-reached`, since it would not have run either way.
- `story_test.go` covers a story whose changes select `go-test` and `integration`, with `integration` covering `go-test`: the steps hold `go-test` skipped, covered by `integration`, and `integration` run; and the same story with a failing check step.

It waits for T-1349, which adds the `skipped` state and `CoveredBy` to `TierResult`.

## Done when

- `flai test --all --json` and `flai verify --json` give a covered tier `"state": "skipped"` and `"covered_by": "<tier>"`; their text says `skipped <name>: covered by <tier>`.
- A covered tier is `skipped`, never `not-reached`, when an earlier step failed.
- `flai test ./internal/verify/ ./cmd/` passes.

## Notes

Layer 3 of S-0342's plan. `story.go`, `story_test.go`, `cmd/verify.go`, and `cmd/verify_test.go` are also in S-0341's touches.
