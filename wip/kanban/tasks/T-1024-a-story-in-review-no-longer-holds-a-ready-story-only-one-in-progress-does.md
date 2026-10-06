---
id: T-1024
type: task
nature: improvement
title: "A story in review no longer holds a ready story; only one in progress does"
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:14:38Z
updated: 2026-10-06T12:15:34Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/workitem/hold.go, flai/internal/workitem/hold_test.go, flai/internal/serve/hold_test.go, flai/internal/mcpserver/work_test.go]
after: [T-1023]
---
# T-1024 A story in review no longer holds a ready story; only one in progress does

## Work

In `flai/internal/workitem/hold.go`, `NewHolds` opens only stories in progress. A story in review no longer enters `h.open`, so it holds no ready story, and the reason text no longer says "in review". `Claim` still works for a story in review, because `flai check`, `flai stream sync`'s trial merge, and `flai accept`'s overlap notice use it. Leave the `Open` call flai serve's launcher makes before an agent moves its story as it is.

Update the tests that expect a hold by a story in review: `hold_test.go`, `flai/internal/serve/hold_test.go`, and `flai/internal/mcpserver/work_test.go`. Add one that reproduces S-0220's four hours from I-0087: a story in review whose claim overlaps every ready story, and `FirstClear` offers the first of them.

This task waits for T-1023, whose ADR decides the rule. T-1027 also changes `hold.go`, so it waits for this task.

## Done when

- A ready story that overlaps only a story in review is not held, and `wait_for_work` offers it. A test shows both.
- A ready story that overlaps a story in progress is still held, with the same reason as before.
- `scripts/flai-test.sh` passes, the race detector included.

## Notes
