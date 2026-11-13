---
id: T-3256
type: task
nature: feature
title: The touches, activity, task plan, inbox, and palette helpers read a state's kind
status: backlog
parent: S-0758
owner: alex
created: 2026-11-12T16:10:00Z
updated: 2026-11-12T16:10:00Z
transitions: []
stream: S-0758
tags: [dashboard]
touches: [flaiover/src/lib/touches.ts, flaiover/src/lib/activity.ts, flaiover/src/lib/taskplan.ts, flaiover/src/lib/inbox.svelte.ts, flaiover/src/lib/server/inbox.ts, flaiover/src/lib/viz/palette.ts]
after: [T-3252]
---
# T-3256 The touches, activity, task plan, inbox, and palette helpers read a state's kind

## Work

Waits for T-3252: these helpers read the module that task adds. Replace each comparison with a state's name in these files with the module's kind or colour: open from closed, active work, the wait for review, and accepted.

## Done when

None of these files compares a status with a state's name, and their tests pass.

## Notes
