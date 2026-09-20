---
id: S-0070
type: story
nature: improvement
title: Cancelling an epic cancels its open stories and their tasks, and cancelling a story cancels its open tasks
status: done
parent: E-0006
owner: alex
created: 2026-09-20T06:14:18Z
updated: 2026-09-20T06:34:48Z
transitions:
  - to: ready
    at: 2026-09-20T06:14:34Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T06:16:44Z
    by: system-flow
  - to: review
    at: 2026-09-20T06:33:56Z
    by: system-flow
  - to: done
    at: 2026-09-20T06:34:48Z
    by: alex
tags: [cli, dashboard]
touches: [flai/internal, flai/cmd, flaiover/src, design/system, template]
---
# S-0070 Cancelling an epic cancels its open stories and their tasks, and cancelling a story cancels its open tasks

## Goal
Cancelling an epic cancels everything under it that is still open: its stories, and their tasks. Today `flai move E-nnnn cancelled` changes the epic alone and leaves its stories in backlog, ready, or in progress under a parent that will never be delivered; the operator has to find and cancel each one, and an agent can still pull one. The same holds one level down: cancelling a story leaves its tasks open. After this story, one cancellation, from the CLI, the board, or MCP, closes the whole subtree, says what it closed, and leaves finished work alone.

## Acceptance criteria
- [x] `flai move E-nnnn cancelled --reason "<why>"` cancels the epic, every story under it that is not done or cancelled, and every task under those stories that is not done or cancelled. Each cancelled child gets its own `cancelled` transition with the same actor and the same timestamp as the epic's, and a reason under its `## Notes` that names the item whose cancellation caused it: `<id> cancelled: <why>`
- [x] Cancelling a story cascades to its open tasks in the same way, whether the story is cancelled directly or through its epic
- [x] Items that are done, already cancelled, or archived are not touched and keep their own history
- [x] It is all or nothing: every transition is validated before anything is written, and when one child cannot be cancelled nothing changes and the error names it
- [x] A story in `review` under a cancelled epic is handled deliberately, not by accident of the rules (today `review` cannot go to `cancelled`): the story decides whether the cascade cancels it, refuses until the operator has accepted or sent it back, or asks; the choice is shown in the preview, and if it changes the allowed transitions an ADR refines ADR-0004
- [x] A story with an open stream is closed the way cancelling that story alone closes it, and the output lists what is left behind for a person to remove or keep: the narrative, the story branch, and the worktree. Unmerged commits on a story branch are never deleted by the cascade
- [x] The operator sees what will happen first: on a terminal `flai move` lists the items it will cancel, by type and current state, and asks for confirmation unless `--yes` is given; `--dry-run` prints the list and changes nothing; `--json` reports every item cancelled. An epic or story with no open children cancels as it does today, with no extra question
- [x] The board does the same: cancelling an epic or a story in flaiover shows the list of children that go with it before it happens, and runs the same flow through the CLI; like every move it is written to the files and committed with whatever is committed next. MCP `item_move` runs the same cascade wherever it may cancel at all
- [x] `flai check` reports an open story under a cancelled epic, and an open task under a cancelled story, so a tree edited by hand or cancelled by an older flai is caught; the message says how to fix it
- [x] An agent working on a story that is cancelled under it finds out: the change appears in `inbox` and `wait_for_events` as a cancellation with its cause, and the conventions say what the agent does then (stop, log, leave the branch alone). The baseline text changes in `template/` first and here in the same story
- [x] `flai stats` and the charts count the cascaded items as cancelled at the time of the cascade, with no change to the definitions in `design/system/metrics.md`; tests cover an epic with stories in every state and tasks under them
- [x] `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, and the users' documentation (`docs/users/flai.md`, `docs/users/flaiover.md`) say what a cancellation now does

## Tasks
- T-0249 workitem: plan and apply a cancellation with everything open under it, all or nothing, and the ADR for stories in review
- T-0250 flai move: show what a cancellation takes with it, confirm on a terminal, --dry-run, --json, and what is left behind; MCP item_move reports the same
- T-0251 flai check: an open story under a cancelled epic, and an open task under a cancelled story
- T-0252 flaiover: cancelling an epic or a story shows what goes with it before it happens
- T-0253 Conventions, design, and documentation: what a cancellation does now, and what an agent does when its story is cancelled

## Notes
Raised by the operator on 2026-09-20, in the middle of S-0069: "I would like to add a new story that supersedes the current WIP - when an epic is cancelled, all of it's stories and tasks should also be cancelled."

What exists. `Repo.Move` in `flai/internal/workitem/rules.go` validates one item; children are consulted only for `done`. The allowed transitions are backlog, ready, and in-progress to cancelled; `review` goes only to `done` or back to `in-progress` (ADR-0004). `design/system/workflow.md` already says a story is not ready while its parent epic is cancelled, which is the only trace of the relationship. When S-0056 was cancelled on 2026-09-19 its tasks were cancelled one by one by hand.

Undoing a cascade is out of scope: cancelled is terminal, and a story wanted after all is written again or moved by hand with a recorded reason.
