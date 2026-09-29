---
id: S-0153
type: story
nature: remediation
title: Threads in the story pages are supposed to be paged
status: done
parent: E-0013
owner: alex
created: 2026-09-29T05:52:49Z
updated: 2026-09-29T07:00:18Z
transitions:
  - to: ready
    at: 2026-09-29T05:57:53Z
    by: alex
  - to: in-progress
    at: 2026-09-29T06:44:20Z
    by: agent-S-0153
  - to: review
    at: 2026-09-29T06:59:53Z
    by: agent-S-0153
  - to: done
    at: 2026-09-29T07:00:18Z
    by: alex
tags: [dashboard]
topics: [client-side]
touches: [flaiover/src, design/system/flaiover-dashboard.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 641
  models:
    - model: claude-opus-5-5
      input: 170
      output: 40088
      cache_read: 10445211
      cache_write: 169923
      cost: 4.2509
---
# S-0153 Threads in the story pages are supposed to be paged

## Goal

Threads were supposed to only be displayed one at a time on story pages with controls to allow the operator to "page" through them instead of them needing to vertically scroll to the correct one.

## Acceptance criteria
- [x] Threads are shown one at a time on the story page
- [x] Operators can page through threads using arrow controls
- [x] A label next to or in between the arrows identifies the cardinal index (1 based) out of the count of threads

## Tasks
- T-0548 The thread pager puts the count between larger arrows beside the heading, and the arrow keys page
- T-0549 A long thread shows its last two entries, with the earlier ones behind a show-earlier control

## Notes

S-0133 already showed one thread at a time. S-0153 moved the count between the arrows (`← n of m →`), put the pager beside the Threads heading as well as under the thread, and made Left and Right page (T-0548). On TH-0035 the operator placed the long scroll on S-0149, whose one thread had eleven entries, and chose to fold a long thread's earlier entries behind a toggle (T-0549). Checked by component tests only: the live dashboard had no host flai connected, and this branch is not built into it.
