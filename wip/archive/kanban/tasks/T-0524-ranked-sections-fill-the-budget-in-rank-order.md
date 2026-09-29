---
id: T-0524
type: task
nature: feature
title: Ranked sections fill the budget in rank order
status: done
parent: S-0146
owner: alex
created: 2026-09-29T03:19:38Z
updated: 2026-09-29T03:55:58Z
transitions:
  - to: ready
    at: 2026-09-29T03:19:45Z
    by: agent-S-0146
  - to: in-progress
    at: 2026-09-29T03:33:54Z
    by: agent-S-0146
  - to: done
    at: 2026-09-29T03:55:58Z
    by: agent-S-0146
stream: S-0146
tags: []
touches: [flai/internal/context]
---
# T-0524 Ranked sections fill the budget in rank order

## Work

The ranked step takes the sections ranked by BM25 in rank order, each cut at its own heading, and keeps each that fits what the budget leaves; an ADR loads whole when it fits, else its best section.

## Done when

- A test shows a tight budget keeps the top-ranked sections that fit and leaves the rest in the catalog.

## Notes
