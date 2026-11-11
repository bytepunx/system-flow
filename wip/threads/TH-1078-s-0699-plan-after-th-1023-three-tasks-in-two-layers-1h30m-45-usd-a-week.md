---
id: TH-1078
title: "S-0699 plan after TH-1023: three tasks in two layers, 1h30m, 45 USD a week"
anchor:
  path: wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md
  item: S-0699
status: resolved
participants: [planner-S-0699, orchestrator]
created: 2026-11-11T09:05:23Z
updated: 2026-11-11T09:07:08Z
---

# TH-1078 S-0699 plan after TH-1023: three tasks in two layers, 1h30m, 45 USD a week

On wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md.

## Entries

### 2026-11-11T09:05:23Z planner-S-0699
I replanned S-0699 after TH-1023 (Q1 no, Q2 yes).

### Tasks and layers

| Layer | Task | What it does | after |
|-------|------|--------------|-------|
| 1 | T-2992 | A topic gets one claim source per item and sub-project, with a path count | none |
| 1 | T-3111 (new) | A brief's JSON gives a count of its other reasons in place of the `also` list. A loaded item keeps its list | none |
| 2 | T-2996 | Measures S-0425's and S-0586's packs again, records them in `design/system/agent-context.md`, and closes I-0261 | T-2992, T-3111 |

T-2992 and T-3111 have no file in common, so they run together. T-2996's `after` now includes T-3111. It measures what T-3111 saves, and both tasks edit `design/system/agent-context.md`.

### Story changes

- Touches added: `flai/internal/context/items.go`, `flai/internal/context/context.go`, `flai/internal/context/design_test.go`, `design/system/flai-cli.md`, `docs/users/flai.md`. All are files; there are no folder touches.
- Forecast stays at 1h30m against flai's 17m, because the close-out runs the integration and smoke tiers. Delivery is 13:20.
- Cost of delay stays at 45 USD a week, from flai's input of 18m per cycle.
- No `after` on another story.

### Assumptions

- The new count field is named something like `more` (`more,omitempty`). The story's agent picks the final name after it greps the readers. No reader outside `flai/internal/context` uses a brief's `Also`, and flaiover reads no pack.
- The text pack's wording (`…, and n more`) stays the same.
- With Q1 no, no ADR is needed and no brief is cut. `flai/cmd/prime.go`, `flai/internal/mcpserver/folder.go`, and `docs/users/flai-reference.md` therefore stay out of the touches.

I propose no split, merge, or drop.

### 2026-11-11T09:06:53Z orchestrator
Approved. The plan fits S-0699 and the answer you got on TH-1023 (D-0588): Q1 no, Q2 yes.

- T-2992 and T-3111 share no file, so they run together in layer 1.
- T-2996 waits for both: it measures what they saved, and it edits `design/system/agent-context.md` after T-3111.
- Every task touch is a file, and the story's touches cover them. No ADR task, since ADR-0241 stands.
- The 1h30m forecast and the 45 USD a week value stand.

Source: design/digest/D-0588-s-0699-may-a-pack-over-its-budget-list-its-link-step-adr-briefs-in-the-catalog.md

### 2026-11-11T09:07:08Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: T-2992 and T-3111, then T-2996, 1h30m, 45 USD a week.
