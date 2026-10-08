---
id: TH-0388
title: "I-0122 duplicates I-0123, which S-0348 fixed: close it"
anchor:
  path: design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md
status: open
participants: [orchestrator]
created: 2026-10-08T09:45:18Z
updated: 2026-10-08T09:45:18Z
---

# TH-0388 I-0122 duplicates I-0123, which S-0348 fixed: close it

On design/issues/I-0122-flai-check-finds-board-wip-limit-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T09:45:18Z orchestrator
Recommendation: close I-0122 as a duplicate of I-0123: `flai issue close I-0122 --reason "Duplicate of I-0123; fixed by S-0348 (ADR-0133): a check scoped to a story leaves board.wip-limit out"`.

I-0122 has the same title and cause as I-0123, with one instance from S-0291's close-out. S-0348, accepted at 917f01e1 and released in flai 1.39.15, fixed the cause and closed I-0123. I-0122 is still open in `design/issues/summary.md`. Closing issues is not among my permissions, so I leave it to you.
