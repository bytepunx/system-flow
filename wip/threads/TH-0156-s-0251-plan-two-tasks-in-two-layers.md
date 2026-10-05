---
id: TH-0156
title: "S-0251 plan: two tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0251-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md
  item: S-0251
status: resolved
participants: [planner-S-0251, alex]
created: 2026-10-05T05:51:47Z
updated: 2026-10-05T05:52:32Z
---

# TH-0156 S-0251 plan: two tasks in two layers

On wip/kanban/stories/S-0251-flai-stream-sync-s-trial-merge-blames-a-story-for-conflicts-between-main-and-another-story-s-stale-branch.md.

## Entries

### 2026-10-05T05:51:47Z planner-S-0251
Here is the plan for S-0251.

**Cause.** `trialMerge` (`flai/cmd/stream_sync.go`) runs `git merge-tree --write-tree ours theirs`, and git picks the merge base. The syncing story has just rebased onto main. The other story's branch is stale, so the base is that branch's old fork point. Everything main gained since then therefore counts as the syncing story's change. A path where main and the stale branch disagree is then blamed on the syncing story. That is TH-0088 for S-0203, and `template/CHANGELOG.md` in TH-0104 and TH-0105 for S-0225.

**Proposed fix.** Keep a conflicted path only when both stories changed it since they left main: `git diff --name-only <main>...story/<id>` intersected with the same diff for the other story. A conflict between main and the other branch alone is left to that story's own sync. The alternative is to merge the other branch onto main first, then merge with ours over main as the base. It is more exact, but it costs a second merge-tree per story. The story's agent proposes the fix in `## Decisions` before building it, as the goal asks.

**Tasks, in two layers:**

- Layer 1: **T-0984**, flai stream sync reports a conflict only on paths both stories changed, not on what main brought. This covers the code, two reproduction tests (S-0203's instance and S-0225's changelog instance), the help text in `stream.go`, and the reference regeneration.
- Layer 2: **T-0985**, the design and the users' guide say how the trial merge finds a conflict, and I-0064 is closed. This covers `flai-cli.md` and `docs/users/flai.md`, then `flai issue close`, which writes the issue and `summary.md`. It waits for T-0984 because it describes what that task built, and the issue closes only once that task's tests pass.

**Assumptions:**

- ADR-0046 and `agent-coordination.md` stay true, so no ADR is needed. If the agent finds the change alters the safety net's contract, it writes one.
- Conflict threads that are already open, such as TH-0088 and TH-0104, resolve at the next sync through `reportConflicts` as they do today. No migration is needed.
- The forecast is 30m, up from flai's 15m by comparison with S-0197 and S-0252, which took about 31m each.
- The cost of delay waits on TH-0150.

### 2026-10-05T05:52:32Z alex
Resolved.
