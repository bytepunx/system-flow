---
id: T-3006
type: task
nature: improvement
title: Close I-0310 with the fix and the figures T-3001 recorded
status: done
parent: S-0676
owner: alex
created: 2026-11-10T10:31:21Z
updated: 2026-11-11T02:22:26Z
transitions:
  - to: ready
    at: 2026-11-11T02:22:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T02:22:18Z
    by: system-flow
  - to: done
    at: 2026-11-11T02:22:26Z
    by: agent-S-0676
stream: S-0676
tags: [flai, issues]
touches: [design/issues/I-0310-a-fresh-flai-list-decodes-10-5-mb-of-kept-items-json-so-it-takes-0-3-to-1-4-s.md, design/issues/summary.md]
after: [T-3001]
usage:
  source: log
  seconds: 6
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 4
      output: 46
      cache_read: 641589
      cache_write: 1097
      cost: 0.2371
---
# T-3006 Close I-0310 with the fix and the figures T-3001 recorded

## Work

Close I-0310 with `flai issue close I-0310 --reason "<reason>"`. The reason names the fix T-2998 built and the fresh list's median before and after, from the section T-3001 wrote in `design/system/server-performance.md`. This task waits for T-3001, whose figures the reason quotes.

`flai issue close` rewrites `design/issues/summary.md`; commit both files.

## Done when

- `flai issue show I-0310` reports it closed, with the reason.
- `design/issues/summary.md` no longer lists I-0310.
- `flai check --strict` reports nothing on either file.

## Notes
