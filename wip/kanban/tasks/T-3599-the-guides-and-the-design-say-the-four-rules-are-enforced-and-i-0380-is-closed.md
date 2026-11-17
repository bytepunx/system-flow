---
id: T-3599
type: task
nature: remediation
title: The guides and the design say the four rules are enforced, and I-0380 is closed naming what fixed each
status: backlog
parent: S-0844
owner: alex
created: 2026-11-17T12:37:14Z
updated: 2026-11-17T12:37:14Z
transitions: []
stream: S-0844
tags: [docs, issues]
touches: [docs/operators/runbooks/change-the-board.md, design/system/workflow.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0380-adr-0301-rules-for-a-project-s-own-workflow-file-that-no-code-enforces-yet.md, design/issues/summary.md]
after: [T-3595, T-3596, T-3597, T-3598]
---
# T-3599 The guides and the design say the four rules are enforced, and I-0380 is closed naming what fixed each

## Work

This task waits for T-3595, T-3596, T-3597, and T-3598. It documents the blocker names, the queue kind, the refusal text, and the check finding those tasks choose, and it closes the issue only once all four are done.

- `docs/operators/runbooks/change-the-board.md`: step 4 of changing the board states the operator-only rule and the dashboard's refusal. Say that flai now enforces both, and name the orchestrator's blocker and the merge queue's kind. In step 3 of the WIP limit section, remove the note that the help of `flai board limit` names only the baseline's states. Say that `flai check` reports `board.wip-state` for a limit on a state without `limit: true`.
- `design/system/workflow.md`, § The workflow file: say how the operator-only acceptance and the save refusal are enforced, and name `board.wip-state`.
- `design/system/flai-cli.md`: add the new blocker to the acceptance gate, and `board.wip-state` to `flai check`'s rules.
- `docs/users/flai.md`: add a row for the new blocker beside `protected_paths` in the orchestrator's blocker table, and add the queue kind and `board.wip-state` where the merge queue's reasons and the check's rules are listed.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`: where the document save is described, say that the editor saves `design/workflows/board.md` only on a story's branch.
- In the story's worktree, run `flai issue close I-0380 --reason "<reason>"`. The reason names, for each of the four instances, the task and the code that fixed it. It also says that instance 4's `SetLimit` and `Save` already followed the workflow, and that this story removed what was left. flai rewrites `design/issues/summary.md`. Commit both. Tick criterion 2 with `flai criteria tick S-0844 2`, and criterion 1 once T-3595 to T-3598 are done and their tests pass.

## Done when

- Each document above describes the enforced rules and names the new blocker, the queue kind, and `board.wip-state` as the code spells them.
- No guide still calls one of the four rules a present limitation.
- I-0380 reads `status: closed` with a reason naming the fix of each instance, `design/issues/summary.md` no longer lists it, and both are committed on `story/S-0844`.
- Criteria 1 and 2 of S-0844 are ticked, and `flai test` on the changed paths passes.

## Notes
