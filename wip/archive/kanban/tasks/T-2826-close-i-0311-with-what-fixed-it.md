---
id: T-2826
type: task
nature: remediation
title: Close I-0311 with what fixed it
status: done
parent: S-0647
owner: alex
created: 2026-11-10T06:24:12Z
updated: 2026-11-10T07:07:29Z
transitions:
  - to: ready
    at: 2026-11-10T07:07:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T07:07:19Z
    by: system-flow
  - to: done
    at: 2026-11-10T07:07:29Z
    by: agent-S-0647
stream: S-0647
tags: [issues]
touches: [design/issues/I-0311-teststartdevserverstopsvitewhenitscontextends-fails-under-load-the-fake-vite.md, design/issues/summary.md]
after: [T-2825]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 440
      cache_read: 381508
      cache_write: 4336
      cost: 0.1198
---
# T-2826 Close I-0311 with what fixed it

## Work

Second layer: it waits for T-2825, because the reason names the fix T-2825 makes and the runs that show it holds.

Run `flai issue close I-0311 --reason "<reason>"` in the story's worktree. The reason says the fake vite was stopped before it wrote its pid, that the tests now check the pid flai started, and names the reproducing run and the passing runs from T-2825's Notes.

## Done when

- I-0311's status is closed, with the reason.
- `design/issues/summary.md` no longer lists I-0311.
- `flai check --strict` is clean.

## Notes
