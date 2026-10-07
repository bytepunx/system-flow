---
id: T-1154
type: task
nature: remediation
title: Close I-0098 with what fixed it
status: backlog
parent: S-0305
owner: alex
created: 2026-10-07T02:19:51Z
updated: 2026-10-07T02:19:51Z
transitions: []
stream: S-0305
tags: [issues]
touches: [design/issues/I-0098-flai-serve-replans-a-story-on-its-own-agent-s-touches-edit-while-the-agent-works-it-and-the-planner-s-new-touches-hold-other-ready-stories.md, design/issues/summary.md]
after: [T-1151]
---
# T-1154 Close I-0098 with what fixed it

## Work

Run `flai issue close I-0098 --reason "..."` in the story's worktree. The reason names S-0305 and says what T-1151 changed: the replanner ignores an edit by the story's own agent, and drops a queued planner for a story in progress or in review, or for a story whose agent runs. The command updates the issue and regenerates `design/issues/summary.md`.

It waits for T-1151 so that the reason states what was built. It shares no path with T-1152 and runs beside it in layer 2.

## Done when

- [ ] I-0098 is closed and its reason names S-0305 and the fix.
- [ ] `design/issues/summary.md` no longer lists I-0098 as open.
- [ ] `flai check --strict` passes.

## Notes

Drafted by planner-S-0305 for the story's second criterion.
