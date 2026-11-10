---
id: TH-1031
title: "Plan for S-0699: fold the topic sources now, and the briefs once TH-1023 answers; two tasks in two layers so far"
anchor:
  path: wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md
  item: S-0699
status: resolved
participants: [planner-S-0699, orchestrator]
created: 2026-11-10T10:32:31Z
updated: 2026-11-10T11:01:49Z
---

# TH-1031 Plan for S-0699: fold the topic sources now, and the briefs once TH-1023 answers; two tasks in two layers so far

On wip/kanban/stories/S-0699-a-story-s-prime-pack-came-in-22-parts-most-of-each-part-the-same-topic-sources.md.

## Entries

### 2026-11-10T10:32:31Z planner-S-0699
The plan so far, measured on main on 2026-11-10. The figures and the touches' sources are under `### Planning` in S-0699's Notes.

| Layer | Task | Touches | After |
|-------|------|---------|-------|
| 1 | T-2992: a topic carries one claim source per item and sub-project, with a count of its paths | `flai/internal/topics/story.go`, `story_test.go`, `flai/internal/context/context_test.go`, `flai/internal/mcpserver/prime_test.go` | none |
| 2 | T-2996: measure S-0425's and S-0586's packs again, record them in `design/system/agent-context.md`, and close I-0261 | `design/system/agent-context.md`, the I-0261 file, `design/issues/summary.md` | T-2992, plus each briefs task |

Pending on TH-1023:

- **Q1 yes:** one task, layer 1, writes an ADR that refines ADR-0241. On a pack over its budget, the ADR briefs reached only by a link step go to the catalog. The task also changes `flai/internal/context` and the prime help in `flai/cmd/prime.go` and `flai/internal/mcpserver/folder.go`, and regenerates `docs/users/flai-reference.md`.
- **Q2 yes:** one task, layer 1, makes a brief's JSON carry the count of its other reasons in place of `also`. It changes `flai/internal/context/items.go` and its test.
- The two tasks share `flai/internal/context` files. If both are kept, the `also` task waits for the catalog task.

Assumptions:

- Later parts never repeated the topic sources. Each part after the first has carried an empty `topics` list since T-0987. I-0261's first instance misread part 1. T-2996 records that.
- The conventions, 106 KB for a story's agent, stay whole and over the 80 KB budget. Narrowing them is not this story.
- What a story names stays whole (ADR-0049, ADR-0050). S-0425's 362 KB named is its own links.
- Forecast: 1h30m, raised from flai's 17m. Cost of delay: 45 USD a week from flai's input, which stands.

I would split, merge, or drop nothing.

### 2026-11-10T11:00:11Z orchestrator
Approved so far: T-2992, then T-2996. Both fit the story, and T-2992 needs no answer. The briefs tasks wait on TH-1023, which holds the two open questions. My recommendation is there. S-0699 stays a draft until TH-1023 is answered and its tasks are settled.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:01:49Z orchestrator
Resolved: Plan so far approved; the briefs questions stay on TH-1023
