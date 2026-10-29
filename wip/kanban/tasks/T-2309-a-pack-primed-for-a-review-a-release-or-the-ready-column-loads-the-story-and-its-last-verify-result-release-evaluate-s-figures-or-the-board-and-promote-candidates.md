---
id: T-2309
type: task
nature: improvement
title: A pack primed for a review, a release, or the ready column loads the story and its last verify result, release_evaluate's figures, or the board and promote_candidates
status: done
parent: S-0537
owner: alex
created: 2026-10-28T22:06:55Z
updated: 2026-10-29T15:51:57Z
transitions:
  - to: ready
    at: 2026-10-29T15:42:20Z
    by: agent-S-0537
  - to: in-progress
    at: 2026-10-29T15:42:21Z
    by: agent-S-0537
  - to: done
    at: 2026-10-29T15:51:57Z
    by: agent-S-0537
stream: S-0537
tags: [flai, prime, orchestrator]
touches: [flai/internal/context/reasons_figures.go, flai/internal/context/reasons_figures_test.go, flai/internal/context/reasons.go, flai/internal/context/reasons_test.go]
after: [T-2308]
usage:
  source: log
  seconds: 576
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 92
      output: 413
      cache_read: 5056221
      cache_write: 136700
      cost: 1.7842
---
# T-2309 A pack primed for a review, a release, or the ready column loads the story and its last verify result, release_evaluate's figures, or the board and promote_candidates

## Work

ADR-0209 §3, criterion 1, for the reasons whose content is a figure flai works out rather than a document. In a new `reasons_figures.go`, load:

- for `review:S-nnnn`, the story whole and its last `flai verify` result, from `verify.LastReport` (`flai/internal/verify/story.go`), or a line saying it has none;
- for `release`, the figures `release_evaluate` gives, from `release.EvaluateRepo` (`flai/internal/release/evaluate.go`);
- for `ready`, the board's columns, limits, and holds, from `LoadBoard` and `NewBoardView` (`flai/internal/workitem`), and the promotion candidates, from `PromotionCandidates` (`flai/internal/workitem/promote.go`).

Each loads as an item with its own step name and the reason that selected it. When importing `verify` or `release` from `context` makes a cycle, move the call behind a function the caller passes in, as `awaitingOperator` in `strategic.go` stands in for `inbox.Summarize`, and say so in the narrative.

Waits for T-2308, which parses the reasons and builds the pack these items go into.

## Done when

- `reasons_figures_test.go` has a test per reason, `review`, `release`, and `ready`, asserting what it loads, and one for a story in review with no verify result.
- `go build ./...` shows no import cycle.
- `flai test` passes on the changed paths.

## Notes
