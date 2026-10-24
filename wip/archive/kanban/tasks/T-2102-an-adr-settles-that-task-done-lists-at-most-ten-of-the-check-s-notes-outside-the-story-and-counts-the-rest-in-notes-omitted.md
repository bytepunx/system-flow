---
id: T-2102
type: task
nature: feature
title: An ADR settles that task_done lists at most ten of the check's notes outside the story and counts the rest in notes_omitted
status: done
parent: S-0492
owner: alex
created: 2026-10-23T23:47:57Z
updated: 2026-10-23T23:49:17Z
transitions:
  - to: ready
    at: 2026-10-23T23:48:15Z
    by: agent-S-0492
  - to: in-progress
    at: 2026-10-23T23:48:16Z
    by: agent-S-0492
  - to: done
    at: 2026-10-23T23:49:17Z
    by: agent-S-0492
stream: S-0492
tags: []
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 61
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3233
      cache_read: 982941
      cache_write: 22105
      cost: 0.4381
---

# T-2102 An ADR settles that task_done lists at most ten of the check's notes outside the story and counts the rest in notes_omitted

## Work

Record with `flai adr new --commit` that the check step of `flai task done` lists at most ten of the check's findings outside the story as `notes`, in the order the check gives them, and counts the rest in `notes_omitted`, so the answer keeps the 40,000 bytes ADR-0179 promises. Refines ADR-0179 and ADR-0085. Waits for nothing.

## Done when

The ADR is committed on `story/S-0492` with its index row, and names I-0202's instance as its evidence.

## Notes
