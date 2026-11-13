---
id: D-0590
title: "S-0712 plan: take main's head once as a commit in storygit.Rebase and check undoes_main against it, four tasks in three layers; proposes merging S-0709: Recommended: merge the two stories into…"
kind: answer
sources: [TH-1017, S-0712]
topics: [git, flai, cli, go, code, all]
paths: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go, flai/internal/storygit/sync_test.go, design/system/workflow.md, docs/users/flai.md, design/issues/I-0338-flai-task-done-s-sync-reports-undoes-main-for-paths-the-story-never-changed.md, design/issues/summary.md]
date: 2026-11-11
status: absorbed
absorbed_by: "design/system/workflow.md#branches-and-collisions-adr-0019"
asked_again: [TH-1101]
cited: [S-0709, S-0712, TH-1091, TH-1101]
---

# D-0590 S-0712 plan: take main's head once as a commit in storygit.Rebase and check undoes_main against it, four tasks in three layers; proposes merging S-0709: Recommended: merge the two stories into…

The operator, alex, answered TH-1017 at 2026-11-11T09:03:11Z.

What was asked:

```text
The plan for S-0712 is written: four tasks in three layers, touches, a forecast, and a cost of delay value.
```

The answer:

```text
**Recommended:** merge the two stories into S-0712, and cancel S-0709.

- Add to S-0712 a third criterion: "I-0335 is closed with `flai issue close I-0335 --reason` saying what fixed it".
- Widen T-2970 to close I-0335 and I-0338 together.
- Have T-2963 also pass the pinned commit to `SettleChangelogs`, which S-0709's plan (TH-1014) found reads main by name after the rebase too.
- Set S-0712's `time_lost_per_cycle` to 6m, I-0335's 3m added to I-0338's 3m. That makes the value about 15 USD a week.

**Why S-0712 and not S-0709:** S-0712's fix resolves main once, before the merge-base read, and uses that one commit for the merge base, the rebase, and the check. S-0709's fix pins it only from the rebase on. S-0712 also already plans the separate real-git reproduction and the `workflow.md` paragraph.

**Why not both:** both stories change the same lines of `storygit.Rebase` for the same race. Two stories would rebuild one fix and conflict on `sync.go`.

This changes scope and an input on S-0712, so it is yours. The same recommendation is on TH-1014.

Source: design/conventions/strategic-agents.md § As the orchestrator
```

Source: TH-1017.

S-0712, accepted on 2026-11-12, built the merge. Its narrative records one widening the recommendation did not name, for T-2963:

```text
T-2963 also passes the pinned commit to the issue fold (fold's o.Files.Fold(res.Base), which runs git ls-tree on main by name after the rebase): the same race would fold the branch's issue into an issue file main gained after the rebase, which the rebased branch does not have.
```

So the sync takes main's head once and uses that one commit for the merge base, the rebase, the undoes_main check, the changelog settle, and the issue fold. Source: S-0712's `## Decisions` (`wip/archive/agents/S-0712.md`).
