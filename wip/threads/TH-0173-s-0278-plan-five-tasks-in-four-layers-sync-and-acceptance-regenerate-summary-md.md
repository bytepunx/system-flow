---
id: TH-0173
title: "S-0278 plan: five tasks in four layers, sync and acceptance regenerate summary.md"
anchor:
  path: wip/kanban/stories/S-0278-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md
  item: S-0278
status: answered
participants: [planner-S-0278, alex]
created: 2026-10-06T11:33:12Z
updated: 2026-10-06T11:35:16Z
---

# TH-0173 S-0278 plan: five tasks in four layers, sync and acceptance regenerate summary.md

On wip/kanban/stories/S-0278-two-stories-that-each-record-or-close-an-issue-always-conflict-in-design-issues-summary-md-whose-updated-line-both-rewrite.md.

## Entries

### 2026-10-06T11:33:12Z planner-S-0278
Plan for S-0278: five tasks in four layers.

1. T-1011 (layer 1, waits for none): choose the fix from I-0074's instances and record it in an ADR. The story asks for a proposal before anything is built.
2. T-1012 (layer 2, after T-1011): `flai stream sync` regenerates `design/issues/summary.md` and continues the rebase when it is the only conflict. `flai accept` gets the same fix through `mergeStoryBranch`. The test reproduces I-0074 with two branches that each record or close an issue.
3. T-1013 (layer 3, after T-1012): the trial merge leaves `summary.md` out of a pair's conflicts, so it opens no thread for it. It reuses T-1012's list of generated paths and the same test file.
4. T-1014 and T-1015 (layer 4, after T-1012 and T-1013, run together since they share no path). T-1014 updates the design, the user guide, and the help, and regenerates the reference. T-1015 closes I-0074.

Assumptions:

- I recommend regenerating on sync and acceptance in T-1011. A merge driver needs `git config` in every clone, which flai does not install. Dropping the `updated:` line does not stop the rows' hunks meeting. Not committing the file is wider than this story. If T-1011's ADR chooses otherwise, T-1012 and T-1013 are rewritten to match.
- Out of scope: two branches bumping the same issue file, as I-0073 did in the third instance, still conflict in its front matter. Say if you want a story for it.
- Forecast 40m, up from flai's 26m, for the decision and the two-branch git test (S-0253 took 35m, S-0197 31m). Cost of delay 22.50 USD a week, from flai's 9m-a-cycle input, unchanged.

Proposal: I-0089 has the same cause as I-0074 (`summary.md` conflicting with issues recorded on main), and this fix removes it too. I would add a criterion closing I-0089, plus a sixth task beside T-1015 that closes it. That would add I-0089's 3m a cycle to the inputs. The story is finalized, so I have not changed its words. Shall I?

### 2026-10-06T11:35:16Z alex
Make a story for front matter collision on issues
