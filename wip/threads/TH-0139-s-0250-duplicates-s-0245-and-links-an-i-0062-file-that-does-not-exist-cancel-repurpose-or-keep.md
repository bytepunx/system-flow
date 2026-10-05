---
id: TH-0139
title: "S-0250 duplicates S-0245, and links an I-0062 file that does not exist: cancel, repurpose, or keep?"
anchor:
  path: wip/kanban/stories/S-0250-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md
  item: S-0250
status: resolved
participants: [planner-S-0250, alex]
created: 2026-10-05T05:44:07Z
updated: 2026-10-05T05:48:22Z
---

# TH-0139 S-0250 duplicates S-0245, and links an I-0062 file that does not exist: cancel, repurpose, or keep?

On wip/kanban/stories/S-0250-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md.

## Entries

### 2026-10-05T05:44:07Z planner-S-0250
Before I plan S-0250 I need two answers from you.

**1. S-0250 is a duplicate.** S-0250 and S-0245 have the same title, goal, and criteria. Both cover the defect that `flai adr new` numbers ADRs from its own worktree only. On main that defect is now I-0063: commit dd3ef02a renumbered it from I-0062 when two issues took the same ID. S-0245 names I-0063. S-0250 still names I-0062 and links `design/issues/I-0062-flai-adr-new-numbers-…md`, a file that does not exist. The I-0062 on main is a different issue: "flai guard lets a task sub-agent run flai adr new but refuses flai adr topics" (count 2, cost 4m).

My recommendation, in this order:

- **(a) Recommended:** cancel S-0250 as a duplicate of S-0245, and plan S-0245 for I-0063. You can cancel it yourself, or reply "cancel" and I will move it to cancelled with that reason. Nothing else is planned on S-0250.
- **(b)** Repurpose S-0250 for the real I-0062, the guard defect. It would get a new title, a goal and criteria naming I-0062, and a correct link. I would rewrite those words only with your yes.
- **(c)** Keep S-0250 for the ADR-number defect, fix its link and IDs to I-0063, and cancel S-0245 instead.

**2. Cost of delay input.** S-0250 has no cost of delay inputs and no epic, so `flai cod` cannot work out a value. My recommendation:

- With (c): `time_lost_per_cycle: 10m`, I-0063's recorded cost. It matches S-0252, the sibling fix for issue numbers.
- With (b): `time_lost_per_cycle: 4m`, I-0062's recorded cost.

**The plan I would write under (c), or for S-0245 under (a):** `adr.NextNumber` (`flai/internal/adr/adr.go`) reads only the local `design/adrs`. S-0252 fixed the same defect for issues with `storygit.FolderNames`, which already lists a folder's files on main, in every story worktree, and on every story branch. The fix reuses it.

Touches:

- `flai/internal/adr`
- `flai/cmd/adr.go` and `flai/cmd/adr_test.go` (co-change 3 of 3)
- `docs/users/flai.md`, `docs/users/flai-reference.md`
- `design/system/flai-cli.md`, `design/system/documentation-standard.md`
- the I-0063 file and `design/issues/summary.md`

Tasks, one per layer:

1. `adr.NextNumber` numbers past main, every worktree, and every story branch, with a test that reproduces the collision.
2. The design and the user guide say so, and I-0063 is closed.

Forecast: about 25m, against S-0252's 34m, because the helper already exists. I am holding until you answer.

### 2026-10-05T05:45:25Z alex
a

### 2026-10-05T05:45:36Z planner-S-0250
Thank you: (a). I could not cancel S-0250 myself, because flai guard refuses a planner any move past backlog, cancellation included (ADR-0060). Please cancel it from the board, with this reason: "Duplicate of S-0245, which remediates I-0063 (renumbered from I-0062 in dd3ef02a)."

I planned nothing on S-0250 and wrote no touches, forecast, cost of delay, or tasks for it. S-0245 has its own planner, planner-S-0245, which has already written its touches, forecast, cost of delay, and tasks T-0915 and T-0916. I left S-0245 alone.

One more thing worth a story of its own: the real I-0062 (the guard lets a task sub-agent run `flai adr new` but refuses `flai adr topics`) has no remediation story. Tell me if you want me to draft one.

### 2026-10-05T05:48:22Z alex
Resolved.
