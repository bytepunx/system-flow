---
id: T-2395
type: task
nature: improvement
title: flai serve starts a ready story that overlaps an in-progress one on a clean trial merge, and journals the conflicting paths when it holds
status: backlog
parent: S-0551
owner: alex
created: 2026-10-31T07:43:25Z
updated: 2026-10-31T07:43:25Z
transitions: []
stream: S-0551
tags: [flai, go]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/hold_test.go, flai/internal/serve/restart.go, flai/internal/serve/start.go]
after: [T-2389, T-2391]
---
# T-2395 flai serve starts a ready story that overlaps an in-progress one on a clean trial merge, and journals the conflicting paths when it holds

## Work

Give the holds that `flai serve` builds the trial-merge judge, backed by the storygit function, in `readyStories`, the ready loop, `askHolds`, `held`, `room` in `start.go`, and `heldStart` in `restart.go`. A ready story that overlaps an in-progress story starts when the trial merge is clean. When it conflicts, the journal's `agent not started` line names the conflicting paths. Run the judge once per pair and look, not once per path, so that a look stays cheap. Leave `askHolds` asking only about holds that still stand.

It waits for the storygit trial merge and the hold verdict's judge, which it wires together.

## Done when

- A test in `hold_test.go` or `agents_test.go`, on a temporary repository, shows an overlapping ready story started on a clean trial merge, and one held on a conflicting trial merge with the conflicting path in the journal line.
- `flai test flai/internal/serve` passes.

## Notes

S-0548 in progress touches `flai/internal/serve/agents_test.go`; agree the order in the two stories' conversation before changing it.
