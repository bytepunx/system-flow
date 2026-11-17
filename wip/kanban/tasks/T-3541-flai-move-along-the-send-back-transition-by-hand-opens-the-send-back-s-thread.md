---
id: T-3541
type: task
nature: remediation
title: flai move along the send-back transition by hand opens the send back's thread with the move's reason
status: in-progress
parent: S-0823
owner: alex
created: 2026-11-17T08:30:49Z
updated: 2026-11-17T10:14:03Z
transitions:
  - to: ready
    at: 2026-11-17T10:14:00Z
    by: agent-S-0823
  - to: in-progress
    at: 2026-11-17T10:14:03Z
    by: agent-S-0823
stream: S-0823
tags: [flai]
touches: [flai/internal/sendback/sendback.go, flai/internal/sendback/sendback_test.go, flai/cmd/move.go, flai/cmd/move_sendback_test.go, docs/users/flai-reference.md]
---
# T-3541 flai move along the send-back transition by hand opens the send back's thread with the move's reason

## Work

Add to `flai/internal/sendback/sendback.go` the send back by hand: given the story as it was before the move, the state it moved to, who moved it, the reason, and now, open the send back's thread with `Thread` when the move followed the workflow's transition with the `send-back` effect out of the story's state (on the baseline, review to in-progress) and the mover is neither `storygit.SentBackAuthor` (`flai`, whose own send backs open their thread first) nor the story's own agent (`agent-S-nnnn`, which is working it). Its title begins with `storygit.SentBackPrefix`, so that flai serve's `startAcceptSentBack` starts the story's agent on it and `resolveAcceptSentBack` resolves it once the story is back in review, and says who sent the story back by hand; its text names who, the reason, and that the story is to be worked and moved to review again. With no reason, as a workflow without `reason-given` allows, the text says none was given and to ask the mover on the thread.

Call it from `flai/cmd/move.go` after `repo.TransitionAll`, which is also what the dashboard's `item.move` runs, with `--by` the operator. Say in the command's help that the move opens the thread, and regenerate `docs/users/flai-reference.md` with `make flai-reference`.

Test it in `flai/internal/sendback/sendback_test.go` and in a new `flai/cmd/move_sendback_test.go` that reproduces I-0397: a story in review moved to in-progress by `alex` with a reason gets one open thread by `flai`, titled with the prefix, holding the reason; the same move by the story's agent or by `flai` opens none; a second send back by hand while the first thread is open adds to it.

This task waits for nothing: it is the first layer.

## Done when

- `flai move S-nnnn in-progress --reason "<why>"` from review, by anyone but `flai` and the story's agent, leaves one open `Acceptance sent back:` thread by `flai` on the story holding the reason.
- The tests above pass, and `TestReferenceIsCurrent` passes with the regenerated reference.

## Notes
