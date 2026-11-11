---
id: T-2968
type: task
nature: remediation
title: Say in the dashboard design and the users guide that a refusal on the new item page is scrolled to and focused
status: done
parent: S-0700
owner: alex
created: 2026-11-10T10:28:21Z
updated: 2026-11-11T08:23:56Z
transitions:
  - to: ready
    at: 2026-11-11T08:23:49Z
    by: agent-S-0700
  - to: in-progress
    at: 2026-11-11T08:23:49Z
    by: agent-S-0700
  - to: done
    at: 2026-11-11T08:23:56Z
    by: agent-S-0700
stream: S-0700
tags: [dashboard, docs]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md]
after: [T-2966]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 2
      output: 463
      cache_read: 240502
      cache_write: 6393
      cost: 0.1085
---
# T-2968 Say in the dashboard design and the users guide that a refusal on the new item page is scrolled to and focused

## Work

In `design/system/flaiover-dashboard.md`, extend the `/new` row's "A refusal shows the findings and keeps the text" with what T-2966 built: the alert is scrolled into view and focused, and the Create button says it was refused until the next submit (S-0700, I-0263). In `docs/users/flaiover.md`, in the new item paragraph, add one sentence saying the page scrolls to the findings and the button says it was refused. Bump each file's `updated`.

Waits for T-2966: it describes the behaviour that task settles, its words and its scroll target.

## Done when

- Both documents say what a refusal on `/new` now does, in the words the form uses.
- `flai test` on the two files passes the markdown lint.

## Notes
