---
id: TH-1267
title: "Plan for S-0817: merge by words a line both sides change"
anchor:
  path: wip/kanban/stories/S-0817-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md
  item: S-0817
status: open
participants: [planner-S-0817]
created: 2026-11-17T02:32:19Z
updated: 2026-11-17T02:32:19Z
---

# TH-1267 Plan for S-0817: merge by words a line both sides change

On wip/kanban/stories/S-0817-stories-of-one-epic-rewrite-the-same-section-of-one-design-document-so.md.

## Entries

### 2026-11-17T02:32:19Z planner-S-0817
The proposed fix: where the section merge finds both sides changing the same line, and the sentence merge (ADR-0287) cannot settle it or the line is a table row, it merges that line again by its words. It keeps both changes when they touch different words. It still stops on the same word, an insertion at the same point, a heading, or a fenced line.

## Why this fix

- E-0033's stops on `design/system/server-performance.md`, I-0387's first instance, ended with ADR-0287 (S-0719, 2026-11-11). No acceptance has stopped on that file since.
- Every acceptance stop since then is on `design/system/flai-cli.md`: S-0689, S-0761, and S-0797. In S-0797's, both sides inserted words at different places in one table row and in one sentence. Today's merge leaves table rows to the line merge and stops when both change one sentence.

## Tasks and layers

| Layer | Task | Waits for |
|---|---|---|
| 1 | T-3495 ADR refining ADR-0287 | none |
| 2 | T-3496 word merge in `flai/internal/storygit/mdmerge.go`, with unit tests | T-3495 |
| 2 | T-3497 design and user docs: `workflow.md`, `agent-coordination.md`, `flai-cli.md`, `docs/users/flai.md` | T-3495 |
| 3 | T-3498 real-git reproductions in `flai/cmd/stream_sync_test.go` and `accept_conflict_test.go` | T-3496 |
| 4 | T-3499 close I-0387 and tick the criteria | T-3498, T-3497 |

T-3496 and T-3497 share no path and run together.

## Figures

- Forecast 1h, delivery 2026-11-17T11:44:00Z. `flai forecast` gave 2m before any touch was declared. S-0719, the same shape, took 42m.
- Cost of delay 600 USD a week, `flai cod`'s figure from the 4h input, kept.
- Twelve file touches, no folder touch, no `after`.

## Assumptions

- A same-point insertion of words stops rather than concatenating, since no order of the two insertions is right. Line insertions keep concatenating, main's first.
- The ADR's number is predicted as 0314. The agent renames the touch if `flai adr new` gives another.
- The trial merge still counts these files as a conflict, as ADR-0276 leaves it. That is out of scope here.

I propose no split, merge, or drop.
