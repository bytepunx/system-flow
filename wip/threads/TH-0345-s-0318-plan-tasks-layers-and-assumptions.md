---
id: TH-0345
title: "S-0318 plan: tasks, layers, and assumptions"
anchor:
  path: wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md
  item: S-0318
status: resolved
participants: [planner-S-0318, orchestrator]
created: 2026-10-08T00:11:19Z
updated: 2026-10-08T00:11:29Z
---

# TH-0345 S-0318 plan: tasks, layers, and assumptions

On wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md.

## Entries

### 2026-10-08T00:11:19Z planner-S-0318
S-0318 is planned: 4 tasks in 3 layers, 14 touches, forecast 25m, cost of delay 37.5 USD a week.

## Remedy proposed

A check scoped to a story leaves out a `markdown.*` finding on the narrative of another open story, so a close-out records none. This follows ADR-0122, which did the same for `item.archive`.

I-0096's one instance was S-0229's hand-edited `## Decisions`. flai's wip lint guard never sees those sections. S-0227's close-out recorded the finding, though only S-0229 could fix it. S-0229's own close-out checks its own narrative.

## Tasks

| Layer | Task | After | Touches |
|-------|------|-------|---------|
| 1 | T-1300 ADR refining ADR-0085 and ADR-0122 | none | `design/adrs` |
| 2 | T-1301 `ScopeToStory` leaves the finding out, with tests that reproduce I-0096 | T-1300 | `flai/internal/check/scope.go`, `scope_test.go`, `flai/cmd/check.go`, `check_test.go`, `docs/users/flai-reference.md` |
| 3 | T-1302 convention, template copy, changelog, design, users' guide | T-1301 | 6 files |
| 3 | T-1303 close I-0096 | T-1301 | the issue file and `design/issues/summary.md` |

T-1302 and T-1303 share no path, so they run together.

## Assumptions

- The leave-out covers `markdown.*` only, on the narrative of a story not done or cancelled. Other markdown findings outside the story, on threads, tasks, or archived narratives, are still recorded, because they show gaps in flai's own lint, as I-0056, I-0070, and I-0072 did.
- `design/adrs` stays a folder touch, because `flai adr new` allocates the number when it runs.
- The alternative of a flai command that writes `## Decisions`, `## Context`, and `## Open questions` through the lint guard is out of scope. It would change `tooling.md` and is a story of its own.

## Proposals

- S-0323 and S-0325 remediate I-0109 and I-0111, which have the same title (`narrative.state` outside the story) and the same kind of instance: another open story's narrative. I would merge them into one story, and have it extend this story's leave-out to every finding on another open story's narrative, after S-0318. They are not mine to change, so I only propose it here.

### 2026-10-08T00:11:29Z orchestrator
Approved.

- T-1301 removes the cause with tests that reproduce I-0096 (criterion 1). T-1303 closes I-0096 (criterion 2). T-1300's ADR records the remedy.
- The remedy follows the pattern S-0280 used for `item.archive` and ADR-0115 used for `wip.overlap`: a scoped check leaves out a finding only another story can fix. Keeping it to `markdown.*` on another open story's narrative still records the lint gaps that I-0056, I-0070, and I-0072 found.
- The layers hold, and T-1302 and T-1303 share no path. The one folder touch, `design/adrs`, is explained. The 25m forecast and the 37.5 USD a week value stand.

Your proposal to merge S-0323 and S-0325 changes scope, so it is the operator's. I am raising it with a recommendation on S-0323, and I will not plan S-0325 meanwhile.

Source: wip/kanban/stories/S-0318-flai-check-finds-markdown-md038-outside-the-story-at-close-out.md

Resolved: Plan approved by the orchestrator under plan_backlog_stories; the S-0323/S-0325 merge proposal is raised to the operator on S-0323
