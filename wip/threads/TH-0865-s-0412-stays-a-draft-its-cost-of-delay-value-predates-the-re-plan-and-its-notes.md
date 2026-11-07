---
id: TH-0865
title: "S-0412 stays a draft: its cost of delay value predates the re-plan, and its Notes promise figures that are not there"
anchor:
  path: wip/kanban/stories/S-0412-a-story-is-worked-to-review-end-to-end-on-the-opencode-adapter-by-a-scripted.md
  item: S-0412
status: resolved
participants: [orchestrator]
created: 2026-11-07T09:36:38Z
updated: 2026-11-07T09:37:28Z
---

# TH-0865 S-0412 stays a draft: its cost of delay value predates the re-plan, and its Notes promise figures that are not there

On wip/kanban/stories/S-0412-a-story-is-worked-to-review-end-to-end-on-the-opencode-adapter-by-a-scripted.md.

## Entries

### 2026-11-07T09:36:38Z orchestrator
**Recommended:** run the planner on S-0412 again, or on E-0023, so it rewrites the cost of delay with `flai cod` and fills in the two Planning notes. I finalize the draft on my next run once those agree.

**Why I have not finalized it:** `flai promote --drafts` lists S-0412 as complete (50m). The criteria, the touches, and the forecast are consistent:
- All five touches are covered by T-1716 to T-1719, and each task's touches stay inside them.
- The 50m forecast was set at 2026-11-07T09:35:53Z, with a basis that fits the fixture-only scope.

The value is not consistent with the re-plan:
- The value is 46.15 USD a week, set by planner-E-0023 at 2026-10-20T21:13:46Z, with no inputs. That date is before the re-plan, so it priced the old `experiment` on a funded OpenRouter key, not today's test on fixtures.
- The Notes say "**Cost of delay.** Rewritten from `flai cod` on this run; see below", and "**Forecast.** Rewritten from flai's figure on this run; see below". Nothing follows either line.

**Alternative:** if 46.15 still stands for the fixture-only story, say so here. The planner or you then drop the two "see below" lines, and I finalize it.

I have not changed S-0412.

### 2026-11-07T09:37:28Z orchestrator
Resolved: planner-E-0023 rewrote S-0412's cost of delay from flai cod at 2026-11-07T09:36:57Z (74.63 USD a week, E-0023's penalty shared by forecast, 50m of 5h35m) and wrote the Forecast and Cost of delay paragraphs under Planning; the value now prices the fixture-only story, so I finalize it.
