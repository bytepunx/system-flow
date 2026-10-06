---
id: TH-0206
title: "S-0286: what cost of delay inputs should it carry?"
anchor:
  path: wip/kanban/stories/S-0286-a-story-that-changes-a-path-claude-code-protects-is-flagged-at-review-and-accepted-only-by-the-operator-flai-checks-permission-prompt-against-each-new-claude-code-and-the-prompt-covers-every-protected-path-in-a-worktree-but-git.md
  item: S-0286
status: open
participants: [planner-S-0286, alex]
created: 2026-10-06T22:50:12Z
updated: 2026-10-06T22:55:27Z
---

# TH-0206 S-0286: what cost of delay inputs should it carry?

On wip/kanban/stories/S-0286-a-story-that-changes-a-path-claude-code-protects-is-flagged-at-review-and-accepted-only-by-the-operator-flai-checks-permission-prompt-against-each-new-claude-code-and-the-prompt-covers-every-protected-path-in-a-worktree-but-git.md.

## Entries

### 2026-10-06T22:50:12Z planner-S-0286
S-0286 has no cost of delay inputs and no epic, so `flai cod S-0286` has nothing to work from. The inputs are yours to set.

**My recommendation:** `time_lost_per_cycle: 1h`, with no revenue and no penalty. At `planning.hour_rate` 150 and a 168h cycle, that gives a value of 150 USD a week.

Why 1h:

- Until auto-approve is safe to leave on, each `.claude/` write in a story waits on a permission thread for you. S-0299's write is waiting on TH-0204 now.
- A Claude Code update that breaks `permission_prompt` hangs a story with no warning. I-0082 and I-0081 recorded 10m and 30m per occurrence, and both undercounted.

Reply `yes` to set that, or give your own figures (`revenue_per_week`, `penalty_per_week`, `time_lost_per_cycle`). I'm drafting the touches, forecast and tasks in the meantime.

### 2026-10-06T22:53:20Z alex
yes

### 2026-10-06T22:55:27Z planner-S-0286
Thanks. I set `time_lost_per_cycle: 1h` on S-0286, and `flai cod` gives 150 USD a week, which I wrote as its value. The plan is on TH-0217.
