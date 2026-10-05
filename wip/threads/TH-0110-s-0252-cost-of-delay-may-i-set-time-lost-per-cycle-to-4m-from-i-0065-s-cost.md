---
id: TH-0110
title: "S-0252 cost of delay: may I set time_lost_per_cycle to 4m, from I-0065's cost?"
anchor:
  path: wip/kanban/stories/S-0252-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md
  item: S-0252
status: resolved
participants: [planner-S-0252, alex]
created: 2026-10-04T23:23:39Z
updated: 2026-10-05T00:06:01Z
---

# TH-0110 S-0252 cost of delay: may I set time_lost_per_cycle to 4m, from I-0065's cost?

On wip/kanban/stories/S-0252-flai-issue-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-issue-number.md.

## Entries

### 2026-10-04T23:23:39Z planner-S-0252
S-0252 has no cost of delay inputs and no epic, so `flai cod S-0252` cannot work out a value. The inputs are yours to give.

**Recommended:** `time_lost_per_cycle: 4m`, with no revenue and no penalty. That is I-0065's recorded cost over its three instances. Since S-0203, flai copies an issue's cost into this input on the stories it makes from the issue: S-0258, S-0260 and S-0262 got theirs that way. S-0252 was made before that, so it has none.

Alternatives:

1. A larger figure, such as 10m, if you count the rebase stops and the hand renumbering. S-0209's instance stopped a rebase on `design/issues/summary.md`.
2. No inputs. S-0252 keeps no value and sorts last by cost of delay.

Reply with the inputs to set, or "yes" for the recommendation. I will hold the cost of delay until you answer, and plan the touches, forecast and tasks meanwhile.

### 2026-10-04T23:28:29Z planner-S-0252
The plan for S-0252. Everything is written except the cost of delay value, which waits for your answer above.

**Cause and proposed fix.** `issues.NextID` globs only the current checkout's `design/issues`. It should number one past the highest I-number found in four places:

- the main checkout;
- every story worktree, which catches issues not yet committed;
- every local `story/*` branch, read with `git ls-tree`;
- the main branch.

**Tasks.** They run in three layers, one task per layer:

1. **T-0831** adds a `storygit` helper that lists the file names under a folder in those four places, with a repository fixture test. It waits for nothing.
2. **T-0832** has `flai issue new` number past every one of those names, with a test that reproduces I-0065, and updates the help and the reference. It waits for T-0831, whose helper it calls.
3. **T-0833** updates `flai-cli.md`, `continuous-improvement.md` and `docs/users/flai.md`, then closes I-0065. It waits for T-0832: the documents describe what that task built, and the issue closes only once its test passes.

**Enrichment.**

- 11 touches, from co-change, layout, and the criteria.
- Topic `cli` added.
- Forecast raised from flai's 12m to 40m, from comparable remediations. Delivery is 2026-10-05T00:30Z, allowing for the hold behind S-0248, which shares `flai-cli.md` and `design/issues/summary.md`.

**Assumptions.**

- Branches and worktrees that cannot be read are skipped. Remote-only story branches are not consulted: agents work on local branches in this host's worktrees.
- An accepted story whose branch was never deleted only makes the numbering skip numbers. That is harmless, since gaps are allowed.
- The fix needs no ADR. It refines how an ID is allocated, not the ID scheme of ADR-0017.

**Proposals.**

- S-0245, the same defect for ADR numbers (I-0063), should reuse T-0831's helper for `design/adrs`. S-0245 is in the backlog and unplanned. I propose giving it `after: S-0252` when it is planned, rather than merging the two stories.
- I would split, merge or drop no task.

### 2026-10-04T23:55:55Z alex
use 10m

### 2026-10-04T23:56:20Z planner-S-0252
Done. I set `time_lost_per_cycle: 10m` on S-0252, and `flai cod` gives a value of 25.00 USD a week, which I wrote unadjusted. It matches S-0258 and S-0262. The Planning notes say why it stands. The plan is complete: T-0831, T-0832 and T-0833 are in the backlog.

### 2026-10-05T00:06:01Z alex
Resolved.
