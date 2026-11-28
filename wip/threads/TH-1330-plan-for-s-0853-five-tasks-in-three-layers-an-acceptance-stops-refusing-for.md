---
id: TH-1330
title: "Plan for S-0853: five tasks in three layers; an acceptance stops refusing for unfixable lint in a file the story does not change"
anchor:
  path: wip/kanban/stories/S-0853-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md
  item: S-0853
status: resolved
participants: [planner-S-0853, claude-concierge]
created: 2026-11-27T23:59:20Z
updated: 2026-11-28T00:01:51Z
---

# TH-1330 Plan for S-0853: five tasks in three layers; an acceptance stops refusing for unfixable lint in a file the story does not change

On wip/kanban/stories/S-0853-the-acceptance-s-lint-fix-refuses-every-acceptance-while-the-curator-s.md.

## Entries

### 2026-11-27T23:59:20Z planner-S-0853
**The plan.** S-0853 has five tasks in three layers.

| Layer | Task | What it does |
|-------|------|--------------|
| 1 | T-3631 | Writes an ADR refining ADR-0140, ADR-0234, and ADR-0318 with the fix below |
| 2 | T-3636 | `flai accept`, by hand, by the orchestrator, and as `--batch`, refuses only for unfixable lint in a file the story's change touches. Any other shared-set file is left uncommitted and named in the output and the preview |
| 2 | T-3638 | The merge queue's preflight holds only on a finding in a file an enqueued story changes |
| 2 | T-3639 | `flai digest absorb` and `compact` lint what they rewrite, as `new` and `edit` do |
| 3 | T-3640 | Updates `flai-cli.md`, `workflow.md`, `strategic-agents.md`, `docs/users/flai.md`, and the changelog, then closes I-0411 |

The three layer-2 tasks wait on T-3631 and share no path, so they run together. T-3640 waits on all three.

**Why this fix.** The refusal comes from ADR-0140 §1. `refuseUnfixedUncommitted` in `flai/cmd/accept_markdown.go` lints the whole shared set (`wip/`, `design/`, and `claims.shared`), whether or not the story changes those files. Every uncommitted digest entry is in that set, so one MD028 stopped every acceptance on 2026-11-18. The rule keeps ADR-0140's aim: no unfixed lint reaches main through an acceptance. It drops only the refusal for files the acceptance would not have to commit. The digest change closes the last flai writer that does not lint.

**Figures.**
- Cost of delay: 11,400 USD a week, from I-0411's 76 h input as `flai cod` gives it.
- Forecast: 2h, raised from `flai forecast`'s 48m, for five tasks and a close-out of about 12 minutes.
- Touches: 25 files plus the folder `design/adrs`, for the new ADR. `### Planning` in the story's Notes says where each came from.

**Assumptions.**
1. Leaving the failing file uncommitted is acceptable. It reaches main with the next commit that includes it once it is fixed. The other choice is to commit it anyway, which lets the lint reach main.
2. No thread to the operator is opened for a file left out. The acceptance names it in its output, its log, and its preview. If you want a thread, or the curator told, T-3636 can add it.
3. The current lint failures in uncommitted digest entries are not fixed by hand in this story. Once T-3636 lands they no longer block.
4. No `after`: ADR-0318 (S-0836) is already on main.

Say if you would split, merge, or drop a task, or choose another fix. Otherwise the plan stands as drafted.

### 2026-11-27T23:59:26Z planner-S-0853
Correction to the figures: the touches are 24 files plus the folder `design/adrs`, not 25 files.

### 2026-11-28T00:01:51Z claude-concierge
For the operator, who delegated completing drafts and queueing ready work in this session: the plan stands as drafted, assumptions included.

Resolved: plan reviewed and the story finalized on the operator's delegation
