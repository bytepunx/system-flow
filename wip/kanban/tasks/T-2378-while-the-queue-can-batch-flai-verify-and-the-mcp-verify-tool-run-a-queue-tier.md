---
id: T-2378
type: task
nature: feature
title: While the queue can batch, flai verify and the MCP verify tool run a queue tier whose when matches the story's touches or changed paths, and leave the rest to the batch
status: backlog
parent: S-0547
owner: alex
created: 2026-10-31T07:40:02Z
updated: 2026-10-31T07:40:02Z
transitions: []
stream: S-0547
tags: [flai]
touches: [flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, docs/users/flai-reference.md]
after: [T-2376]
---
# T-2378 While the queue can batch, flai verify and the MCP verify tool run a queue tier whose when matches the story's touches or changed paths, and leave the rest to the batch

## Work

Criterion 2. `SelectStoryLeavingQueue` (`flai/internal/verify/select.go`) answers every selected `queue` tier apart as left to the batch (ADR-0174). Give it the story's touches, and keep in `run` each selected `queue` tier whose `When` matches one of those touches or one of the branch's changed paths, by the `paths` pattern language; a touch that is a folder or a component name matches as the claim reads it. A tier kept to run covers as any other does, so `integration` run here skips `go-test` again.

Carry the touches through the verification's options in `flai/internal/verify/story.go`: the story's claim, its touches and its open tasks', read where `flai/cmd/verify.go` and `flai/internal/mcpserver/verify.go` set `LeaveToQueue`. Name in the report why a `queue` tier ran, its `when` pattern and the touch or path it matched, beside the `queued` tiers' "left to the merge queue" line, in the text and the JSON. Keep ADR-0221: a `queue` tier the close-out ran at the head counts for a hand acceptance, and `flai verify --all` at the same patch reuses it. The batch's verification is unchanged: it never leaves a tier out.

Update `flai verify`'s help where it says which tiers it leaves to the queue, and regenerate `docs/users/flai-reference.md` with `make flai-reference`, which I-0203's stale reference shows a test checks.

Waits for T-2376, whose `Tier.When` it reads.

## Done when

- Tests in `select_test.go` and `story_test.go`: a `queue` tier whose `when` matches a touch runs, one whose `when` matches a changed path runs, one whose `when` matches neither is `queued`, and one without `when` is `queued` as today; with the queue unable to batch every selected tier runs as before.
- `flai verify` and the MCP `verify` tool name why a `queue` tier ran, in the text and the JSON, in tests in `flai/cmd/verify_test.go` and `flai/internal/mcpserver/verify_test.go`.
- `docs/users/flai-reference.md` is regenerated, and `flai test` on the changed paths passes.

## Notes

Drafted by planner-S-0547. Matching against the touches and the changed paths, not the touches alone, is an assumption on the plan's thread: a story whose touches are narrower than its change still runs the tier.
