---
id: T-0330
type: task
nature: feature
title: "Tried end to end: a real run, a real cancel that kills a spawned grandchild, a real timeout"
status: done
parent: S-0082
owner: alex
created: 2026-09-22T22:37:18Z
updated: 2026-09-22T23:45:37Z
transitions:
  - to: ready
    at: 2026-09-22T23:29:16Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T23:29:18Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:45:37Z
    by: system-flow
stream: S-0082
tags: []
---
# T-0330 Tried end to end: a real run, a real cancel that kills a spawned grandchild, a real timeout

## Work
In a scratch project with its own worktree for a throwaway story, own `flai serve`, own config with the `checks` action enabled and two named checks configured (one that passes, one that spawns a background child of its own and then sleeps, to prove group-kill reaches it): from a Playwright session against the review page (not just the CLI, since T-0329's polling/tail loop is the part most likely to be wrong in ways a component test's mocks cannot show):
1. Run checks; watch the output arrive live as each command runs.
2. Cancel mid-run; confirm on the host, by PID, that both the check process and the background child it spawned are gone, not orphaned; confirm the page shows `cancelled`.
3. Run again and let it finish (both checks passing); confirm outcome and duration are shown, and still shown after leaving the page and coming back.
4. A check configured to run longer than the time limit; confirm it is stopped at the limit, outcome `timed-out`, nothing left running.
5. Start a run, then try to start a second one for the same story while the first is active; confirm it is refused, not queued or replaced.

Tick the acceptance criteria against what was actually seen. This does not touch the operator's own dashboard container or `flai serve` at all — a scratch `flai serve` and a scratch project worktree are enough; no Docker, no shared container name to collide with.

## Done when
All five checks pass for real; the story's three acceptance criteria are ticked from what was verified; scratch resources (the scratch `flai serve` by exact PID, the scratch project directory) removed afterward.

## Notes
Depends on T-0326–T-0329. Lower risk than S-0081's equivalent task: nothing here is shared, host-wide, or hard to isolate.
