---
id: T-0517
type: task
nature: feature
title: Record I-0050's second instance and the narrowed claim
status: done
parent: S-0145
owner: alex
created: 2026-09-29T02:12:49Z
updated: 2026-09-29T02:16:00Z
transitions:
  - to: ready
    at: 2026-09-29T02:13:58Z
    by: agent-S-0145
  - to: in-progress
    at: 2026-09-29T02:13:59Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:16:00Z
    by: agent-S-0145
stream: S-0145
tags: []
---

# T-0517 Record I-0050's second instance and the narrowed claim

## Work

Bump I-0050 with this evening's second instance (flai serve 1.23.0 started agent-S-0145 for a held story) with `flai issue bump`, and record in the narrative's `## Decisions` that the claim was narrowed from `flai/cmd` to `design/adrs` and `design/issues`, because a research story that writes an ADR draft touches no code.

## Done when

I-0050 has count 2 and a dated instance, `design/issues/summary.md` is regenerated, and the narrative records the narrowed claim with its reason.

## Notes
