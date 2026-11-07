---
id: T-2506
type: task
nature: remediation
title: File each fault the publish run found as an issue, or as a fix task of S-0562, and say in the Notes when it found none
status: done
parent: S-0562
owner: alex
created: 2026-11-07T03:45:32Z
updated: 2026-11-07T09:48:04Z
transitions:
  - to: ready
    at: 2026-11-07T09:47:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T09:47:58Z
    by: system-flow
  - to: done
    at: 2026-11-07T09:48:04Z
    by: agent-S-0562
stream: S-0562
tags: [template]
touches: [design/issues/summary.md]
after: [T-2505]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 466
      cache_read: 214535
      cache_write: 2260
      cost: 0.0703
---
# T-2506 File each fault the publish run found as an issue, or as a fix task of S-0562, and say in the Notes when it found none

## Work

Go through what T-2504 and T-2505 recorded: a precondition that did not hold, a step that failed or stopped where it should not have, a refusal, and anything in the formula's steps or its body that did not match what ran.

For each fault, choose one:

- A fix small enough for this story: write it as a new task of S-0562 with `flai task new`, naming the files it changes, such as `design/workflows/formulas/publish.md` and its template copy, and work it.
- Anything larger, or a change to `flai` itself: record it with `flai issue new`, or bump the open issue that already records it with `flai issue bump`.

When the run found nothing wrong, say so in S-0562's Notes with `flai edit`.

This task waits for T-2505, because the faults are known only once the instance has run to acceptance.

## Done when

- Each fault the run found is fixed in a done task of S-0562 or named by an issue, and S-0562's Notes list them with the task or issue ID.
- Or S-0562's Notes say the run found nothing wrong.

## Notes
