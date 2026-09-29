---
id: T-0518
type: task
nature: feature
title: Measure where the context pack's size comes from
status: done
parent: S-0145
owner: alex
created: 2026-09-29T02:12:50Z
updated: 2026-09-29T02:16:46Z
transitions:
  - to: ready
    at: 2026-09-29T02:16:00Z
    by: agent-S-0145
  - to: in-progress
    at: 2026-09-29T02:16:00Z
    by: agent-S-0145
  - to: done
    at: 2026-09-29T02:16:46Z
    by: agent-S-0145
stream: S-0145
tags: []
---

# T-0518 Measure where the context pack's size comes from

## Work

Run `flai prime --story --json` on a cli story and break the pack down by step (topics, linked, ranked), by folder, and by document, so the ADR draft's context states where the bytes come from and what each candidate approach would remove. Estimate what heading topics on the largest files and a budget would each save.

## Done when

The narrative and the ADR draft carry a table of the pack's size by step and by the largest documents, measured on this repository, with the command that produced it.

## Notes
