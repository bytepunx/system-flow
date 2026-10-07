---
id: T-1164
type: task
nature: improvement
title: An ADR refining ADR-0085 and ADR-0115 records the remedy for I-0078, proposed from its instances
status: backlog
parent: S-0280
owner: alex
created: 2026-10-07T15:03:25Z
updated: 2026-10-07T15:03:25Z
transitions: []
stream: S-0280
tags: [flai]
touches: [design/adrs]
---
# T-1164 An ADR refining ADR-0085 and ADR-0115 records the remedy for I-0078, proposed from its instances

## Work

Write an ADR with `flai adr new` that records the remedy for I-0078, proposed from its 38 instances. It refines ADR-0085 and ADR-0115.

What the instances show:

- All 38 name one file, S-0250's, cancelled by the operator from backlog on 2026-10-05 and left in `wip/kanban/stories` until commit 1f7b7296 archived it by hand.
- `item.archive` (`flai/internal/check/check.go`) warns on an epic or story that is done or cancelled and not archived. A story in progress is never closed, so at close-out the finding never names the closing story. It is the main checkout's housekeeping: no story branch makes it, and no story's agent can fix it in its worktree.
- `ScopeToStory` (`flai/internal/check/scope.go`) marks it outside, and `recordOutside` (`flai/cmd/check.go`) records it, so every close-out bumps I-0078 while a cancelled item lingers.

The proposed remedy, as ADR-0115 did for `wip.overlap`:

- A close-out records no `item.archive` in an issue.
- Its scoped check leaves out an `item.archive` that does not name the story.
- `flai check` in the main checkout still warns `item.archive`, so the operator still sees what to archive.

Record the alternative considered and why it was not taken: cancelling archives what it cancels. ADR-0055 lets a cancelled item move back to backlog, and ADR-0028 leaves a cancelled story's narrative where it is, so archiving on cancel would reverse both.

This task waits for nothing. Every other task builds what it decides.

## Done when

- The ADR is accepted under `design/adrs` and listed in `design/adrs/README.md`.
- It names I-0078, its instances, the remedy, and the alternative not taken.
- It refines ADR-0085 and ADR-0115, and `flai check --strict` passes.

## Notes

`design/adrs` is a folder touch: `flai adr new` picks the file's number and slug. It lies inside `claims.shared`, so it holds no story.
