---
id: TH-0610
title: "S-0428 replanned: six tasks in four layers stand, forecast 1h1m, value 95.06 USD a week from the 4h input TH-0476 confirmed and now recorded on E-0025"
anchor:
  path: wip/kanban/stories/S-0428-a-digest-entry-is-one-markdown-file-under-design-digest-that-flai-writes-lists-and-checks-and-an-adr-records-where-it-lives-and-how-it-is-written.md
  item: S-0428
status: resolved
participants: [planner-S-0428, orchestrator]
created: 2026-10-23T16:21:20Z
updated: 2026-10-23T16:21:28Z
---

# TH-0610 S-0428 replanned: six tasks in four layers stand, forecast 1h1m, value 95.06 USD a week from the 4h input TH-0476 confirmed and now recorded on E-0025

On wip/kanban/stories/S-0428-a-digest-entry-is-one-markdown-file-under-design-digest-that-flai-writes-lists-and-checks-and-an-adr-records-where-it-lives-and-how-it-is-written.md.

## Entries

### 2026-10-23T16:21:20Z planner-S-0428
I replanned S-0428 on the operator's confirmation of TH-0476. Nothing in its tasks changed; the figures did.

## Tasks and layers

| Layer | Tasks | Why it waits |
|-------|-------|--------------|
| 1 | T-1807, the ADR | Fixes the folder, the file name, the front matter, and the rules every other task builds on |
| 2 | T-1808, `flai/internal/digest` | Reads, validates, writes; refuses a secret's shape. Waits for the ADR |
| 3 | T-1809 (`flai check`), T-1810 (`flai digest` commands), T-1811 (search and documents-tree tests) | Each calls the package; no two share a path, so they run together |
| 4 | T-1812, READMEs, layout, CLI design, user docs, template | Describes what was built |

I revisited each task against the six criteria: each criterion has its task, the touches name files (the one folder touch is `design/adrs`, in `claims.shared`), and the `after` chain is as above. I changed none of them.

## Figures

- **Touches:** the 28 the tasks name, kept. `flai touches suggest` adds nothing I took: the operator docs and the manifest need no change, since the folder follows `layout.design` and adds no setting, and the documents page lists the folder from `docs.tree` as it lists `design/analysis`, so no flaiover change. Details under the story's `### Planning`.
- **Forecast:** flai's 1h1m, delivery 2026-10-24T06:09Z (21st in the pull order). Up from the 51m of the first plan: the size now falls in the large band. Kept, since S-0414, of the same shape, was forecast 1h1m.
- **Cost of delay:** 95.06 USD a week, S-0428's share of E-0025's 600 USD a week by forecast.

## The epic's input

You confirmed on TH-0476 (2026-10-21T12:59Z) the recommendation to set `time_lost_per_cycle: 4h` on E-0025, but the input was not on the epic, so `flai cod` still refused every story of E-0025. I recorded that one input on E-0025 as your confirmed answer, as planner-E-0023 did on TH-0470, and set nothing else. If you meant another figure, change it with `flai edit E-0025` and ask for the stories to be planned again.

## Assumptions

1. TH-0476's answers stand: one file per entry under `design/digest`, and the manifest's `owner` is the operator for an answer that becomes an entry (S-0433).
2. `checker.digest` keeps the folder out of the design documents' generic front-matter rules, as `documentation()` does for `analysis` and `experiments`, since an entry carries `date`, not `updated`.
3. No site-menu entry for the digest in flaiover: the story asks that the documents tree list the folder, which it does with no folder list. A menu entry would be a later story's.
4. The seven other stories of E-0025 now have an input to share; their values are written when each is planned.

The story stays a draft in the backlog, yours to finalize.

### 2026-10-23T16:21:26Z orchestrator
Approved: the plan fits S-0428.

- **Tasks and layers.** Six tasks in four layers, one task per criterion, from T-1807 for criterion 1 to T-1812 for criterion 6. T-1809, T-1810, and T-1811 share no path, so they can run together.
- **Touches.** The 28 touches name files. The one folder, `design/adrs`, is kept for a stated reason and sits inside `claims.shared`.
- **Forecast.** flai's 1h1m, kept against S-0414, which has the same shape.
- **Cost of delay.** 95.06 USD a week from `flai cod`. It rests on E-0025's `time_lost_per_cycle: 4h`, which the operator confirmed on TH-0476 and you recorded on the epic.

I set no input on the story. The epic's input covers it.

Source: design/conventions/strategic-agents.md § As the planner

### 2026-10-23T16:21:28Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
