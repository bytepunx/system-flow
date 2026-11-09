---
id: T-2804
type: task
nature: improvement
title: Where a board call still takes 300 ms or more, file the shortfall as an issue and ask the operator on a thread how the criterion stands
status: backlog
parent: S-0636
owner: alex
created: 2026-11-09T22:07:26Z
updated: 2026-11-09T22:07:26Z
transitions: []
stream: S-0636
tags: [cli]
touches: [design/issues/summary.md]
after: [T-2803]
---
# T-2804 Where a board call still takes 300 ms or more, file the shortfall as an issue and ask the operator on a thread how the criterion stands

## Work

Waits for T-2803: it reads the shortfall that task records.

When T-2803 finds every call of every project under 300 ms, tick the story's criterion and cancel this task, saying so.

Otherwise, follow D-0532 and D-0535:

- Record the shortfall as an issue with `flai issue new`, citing the phase T-2803 named, its cause in `design/system/server-performance.md`, and this story. When an open issue already records that cause, bump it instead. `release.pending` working out its plans from every item is one candidate (§ Cause 14 fixed).
- Open a thread on S-0636 with the figures and a recommendation: reword the criterion to the fall measured against `main` on the same host, with the 300 ms target left to the issue. Rewording the criterion is the operator's call.
- Name the issue in the story's Notes.

## Done when

- [ ] Every board call answered under 300 ms and the criterion is ticked, or an issue records the shortfall, the story's Notes name it, and a thread on S-0636 asks the operator how the criterion stands.

## Notes
