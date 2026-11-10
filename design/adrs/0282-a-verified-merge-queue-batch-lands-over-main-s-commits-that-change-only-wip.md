---
id: ADR-0282
title: "A verified merge queue batch lands over main's commits that change only wip paths none of its stories change, and every dropped batch is logged as a dropped event that flai stats counts per week"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0173, ADR-0177, ADR-0227]
---

# ADR-0282 A verified merge queue batch lands over main's commits that change only wip paths none of its stories change, and every dropped batch is logged as a dropped event that flai stats counts per week

## Context

The merge queue drops a verified batch when the main branch's head is not the commit the batch expects: the commit it was built on, or, once its landing began, where the last landing left main ([ADR-0173](0173-a-merge-queue-in-flai-serve-with-no-model-runs-behind-the-merge-queue-host.md)). `mergequeue.DropIfMainMoved` drops it before and between landings, and `mergequeue.Landable` refuses `flai accept --batch` for the same reason. The rebuilt batch verifies the same stories again.

[I-0287](../issues/I-0287-the-merge-queue-builds-one-story-batches-and-runs-three-verifications-per-story.md) measured 2.95 verifications per landed story. The queue's log for 2026-10-21 to 2026-11-02 holds 94 passing `verified` lines, and 37 of those batches never landed: each was dropped after its whole verification, about ten minutes of the `integration` and `smoke` tiers. In the three latest drops sampled, main had moved by one commit, flai's `chore: replan forecasts after reordered`, which changed only files under `wip/`.

Such a commit cannot change what the batch verified:

1. **`wip/` is written on the main branch** ([ADR-0019](0019-story-branches-and-touches.md)); a story's branch does not change it, and the batch's tiers scope themselves to the batch's stories ([ADR-0183](0183-the-merge-queue-s-batch-and-each-bisection-step-run-their-tiers-with-merge.md)).
2. **`flai accept --batch` rebases the story's branch onto main before it merges** (`mergeStoryBranch`), so a landing onto a main that gained only `wip` commits no story of the batch changes is clean, and what lands outside `wip/` is what the batch verified.
3. **`flai verify`'s `sync` step already passes over the same kind of commit**, for the same reasons ([ADR-0135](0135-the-sync-step-of-flai-verify-and-the-close-out-s-last-check-pass-over-commits.md)).

The log records no line for a drop. A drop shows in `flai stats` only as `verifications` above `landed` ([ADR-0227](0227-the-merge-queue-logs-each-bisection-step-s-verification-as-a-step-verified.md)), and cannot be told from a bisection. [ADR-0177](0177-flai-stats-reports-per-week-the-stories-the-merge-queue-landed-or-sent-back.md) left a dropped batch uncounted. `design/system/metrics.md`, the contract between `flai stats` and the dashboard, changes only with an ADR.

## Decision

A verified merge queue batch lands over main's commits that change only `wip` paths none of its stories change, and every dropped batch is logged as a `dropped` event that `flai stats` counts per week.

1. **When main moved and the batch is kept.** Let the expected commit be the one the batch expects main at, as today. The batch is not dropped, and `Landable` gives no reason on main's account, when all of these hold:
   - the expected commit is an ancestor of main's head (`git merge-base --is-ancestor`);
   - every path `git diff --no-renames --name-only <expected> <main>` lists is under the manifest's `wip` folder;
   - none of those paths is among the paths the batch's stories change, the paths `git diff --no-renames --name-only <base> <last rebased head>` lists, from the commit the batch was built on to the rebased head of its last story.
2. **What still drops it.** A main that does not descend from the expected commit, a path outside the `wip` folder, or a `wip` path a story of the batch changes, drops the batch as today, its reason naming which and the paths. So does a git read that fails, its reason naming the error: the drop is the safe side.
3. **The expected commit does not move.** Each later look diffs from the same expected commit, the batch's base or where its last landing left main, so a commit outside `wip/` that lands after a `wip`-only one still drops the batch. A landing sets the expected commit where it leaves main, as today.
4. **The batch lands without verifying again.** It lands from the commits it verified, through `flai accept --batch` as before ([ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)).
5. **The `dropped` event.** `mergequeue.Update`, the state's only writer, appends a `dropped` line when a write sets a batch's `dropped`, so that every drop is logged whoever drops it: main moving, a landing that failed, a bisection that could not run, a story no longer enqueued, or `flai queue remove`. The line has `at` (the drop's time), `batch`, `stories` (the batch's), `by` (who dropped it), and `reason` (the drop's sentence).
6. **The figure.** Under `queue.weeks[]`, `dropped` counts the `dropped` lines of the week. A dropped batch still counts in `batches`, as ADR-0177 has it, since it was built. `flai stats` prints `dropped <n>` on the batches line when it is not 0.

This refines ADR-0173, whose queue dropped a verified batch whenever main moved other than by its landing; ADR-0177, whose log left a drop out and whose stats did not count it; and ADR-0227, whose verifications per story landed a drop raised with no figure of its own.

## Consequences

- flai's own replans, edits, and moves on main no longer cost a batch its verification. Most of I-0287's dropped batches, and the verifications they cost, go.
- A batch can land with main ahead of its base by `wip` commits. The landing merges them as any acceptance does; it cannot conflict, since no story of the batch changes their paths.
- `flai stats` tells drops from bisections: `verifications` less `dropped` less the bisection steps is the batches that ran to a landing or a send-back.
- Weeks before the release that writes `dropped` count none, which is no evidence either way.
- One-story batches are unchanged: `merge_queue.coalesce` ([ADR-0226](0226-the-merge-queue-waits-merge-queue-coalesce-3-minutes-by-default-from-the.md)) already joins stories enqueued close together.

## Alternatives considered

- **Verify again only the tiers main's new paths select, as ADR-0156 does for a close-out.** A `wip`-only commit selects no tier here, so it reduces to this rule, at the cost of a second mechanism.
- **Move the expected commit to main's head at each look that passes over.** Equivalent while main gains only `wip` commits, but it hides which commits the batch was verified against, and a later diff would no longer span everything since the verification.
- **Hold flai's `wip` commits while a batch is verified or lands.** It would delay every writer of `wip/` for the length of a verification, the reason ADR-0135 gives.
- **Count drops from `state.json`.** Each write replaces it, so a week gone cannot be counted, the reason ADR-0177 gives.
- **Have each drop site append its own line.** A drop site that forgets undercounts silently; deriving it in `Update`, as ADR-0177 does for every other event, covers every writer.
