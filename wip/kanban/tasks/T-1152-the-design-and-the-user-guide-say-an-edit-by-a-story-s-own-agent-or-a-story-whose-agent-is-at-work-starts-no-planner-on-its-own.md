---
id: T-1152
type: task
nature: remediation
title: The design and the user guide say an edit by a story's own agent, or a story whose agent is at work, starts no planner on its own
status: backlog
parent: S-0305
owner: alex
created: 2026-10-07T02:19:47Z
updated: 2026-10-07T02:19:58Z
transitions: []
stream: S-0305
tags: [docs, planner]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-1151]
---
# T-1152 The design and the user guide say an edit by a story's own agent, or a story whose agent is at work, starts no planner on its own

## Work

Say what T-1151 built, in its words, in three places:

- `design/system/strategic-agents.md`, section The planner › Planning again: in the triggers table's "An edit" row and the paragraph below it, an edit by the story's own agent never triggers, as a planner's does not. In The queue, an entry whose story is in progress or in review, or whose story's agent runs, is dropped and logged as `queued planner dropped`. Bump `updated`.
- `design/system/flai-cli.md`, the `flai plan` row: the replanner starts a queued run after `planCheck` and the check that the story's agent is not at work. Bump `updated`.
- `docs/users/flai.md`, the table under "While `plan` is on, `flai serve` also plans again on its own": the edit row says that the story's own agent's edits start nothing, and neither does any trigger while the story's agent works it.

It waits for T-1151 so that the words match the code as built. It shares no path with T-1154 and runs beside it in layer 2.

S-0245, S-0269, and S-0277 are in progress and claim `design/system/flai-cli.md` and `docs/users/flai.md`. If `flai stream sync` stops on them, keep both sides.

## Done when

- [ ] The three documents say that the story's own agent's edits, and a story whose agent is at work, start no planner on their own, and that the operator's `flai plan` still does.
- [ ] `flai check --strict` and the markdown lint pass.

## Notes

Drafted by planner-S-0305. The rows were found by searching for "plans again" and "Planning again" in `design/` and `docs/`.
