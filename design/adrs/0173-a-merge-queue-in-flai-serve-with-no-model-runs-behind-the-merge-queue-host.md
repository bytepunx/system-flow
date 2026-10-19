---
id: ADR-0173
title: "A merge queue in flai serve, with no model, runs behind the merge-queue host action and under accept_reviews, batches at most merge_queue.batch_max enqueued eligible stories under local acceptance only, and never batches a story whose branch undoes main's change"
status: proposed
date: 2026-10-19
supersedes: []
superseded_by: []
refines: [ADR-0093, ADR-0155]
topics: [orchestration, flai, cli, dashboard]
---

# ADR-0173 A merge queue in flai serve, with no model, runs behind the merge-queue host action and under accept_reviews, batches at most merge_queue.batch_max enqueued eligible stories under local acceptance only, and never batches a story whose branch undoes main's change

## Context

E-0022 brings a merge queue to flai. It verifies the stories in review as one batch, bisects a red batch to the story that fails, and lands the green ones.

Acceptance is one story at a time, and it stops. [I-0150](../issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md) measures acceptances that stop on a dirty main checkout or on a rebase inside a shared path. S-0446's acceptance stopped because its branch was 8 commits behind main (TH-0511).

The epic left three questions to the designer: what runs the queue, how large a batch may be, and whether a story whose branch undoes main's change ([ADR-0155](0155-a-rebased-story-branch-that-undoes-main-s-change-to-a-path-is-recorded-at-sync.md)) may be batched. The planner recommended answers on TH-0465, questions 2 to 4. After E-0021 landed pull-request acceptance, the planner made the queue local-only (TH-0465, TH-0466). TH-0516 asks the operator to confirm these answers. This ADR records them.

## Decision

A merge queue in `flai serve`, with no model, runs behind the `merge-queue` host action and under `accept_reviews`. It batches at most `merge_queue.batch_max` enqueued, eligible stories, under local acceptance only. It never batches a story whose branch undoes main's change.

- **What runs it.** A loop of `flai serve`'s own, with no model. A new host action, `merge-queue`, turns it on; it is off by default. It runs only while `orchestration.permissions.accept_reviews` is on. Checking criteria, verifying, and bisecting are arithmetic and git, so a model would add cost and no judgement.
- **Enqueueing.** A story enters the queue only when someone enqueues it. The orchestrator enqueues with [ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)'s evidence. The operator enqueues with `flai queue add` or from the review page. Eligibility alone enqueues nothing.
- **Landing.** The queue lands each story as the one who enqueued it, through `flai accept`'s own flow. `flai accept` by hand is unchanged, and takes the story out of the queue.
- **Eligibility.** A story in review is eligible when all of these hold:
  - its last `flai verify` passed at the head of its branch;
  - every acceptance criterion is ticked;
  - no thread on the story or one of its tasks is open;
  - no path is recorded as undoing main's change (ADR-0155).

  Eligible stories are ordered by when they last entered review. `flai queue` lists each other story in review with every reason it is not eligible.
- **A branch that undoes main's change.** It is never batched. `flai queue` shows it as not eligible and names the paths, until the operator restores them, keeps the removal, or sends the story back. The orchestrator may not keep a removal, so a landing as the orchestrator could not land it either.
- **Batch size.** At most 4 stories by default, set as `merge_queue.batch_max` in `system-flow.yaml`. A red batch of 4 bisects in at most two more runs. 4 is under this project's review limit of 5.
- **Local acceptance only.** Under `acceptance.mode: pull-request` ([ADR-0165](0165-a-project-chooses-local-or-pull-request-acceptance-in-its-manifest-flai-reaches.md)), the forge's required checks gate each story ([ADR-0169](0169-under-pull-request-acceptance-every-acceptance-the-operator-s-and-the.md)). A batch rebases each story, so the kept forge state would be read at a commit other than the branch's head, and every landing would stop on `forge_stale`. In that mode the queue builds no batch. `flai queue`, `flai board`, the MCP `board`, and the dashboard say that the queue does not run, and list no story as eligible.
- **State.** One file in the main checkout, `.flai-cache/merge-queue/state.json`. It is host-local, like `flai verify`'s last result and ADR-0155's record. Only the package `flai/internal/mergequeue` reads and writes it. It holds:
  - the stories enqueued, with who enqueued each and when;
  - the batch under verification, and its step;
  - each story's last outcome.

  With no file, the queue has no batch and no outcomes. Nothing goes into front matter.
- **Which story builds each part.**

  | Story | Part |
  |-------|------|
  | S-0398 | Eligibility; the state's reader; `flai queue`; the queue on `flai board`, the MCP `board`, `board.get`, and the dashboard, read-only |
  | S-0399 | Enqueueing; the `merge-queue` host action; `batch_max`; building a batch in a worktree of its own, rebased on main in order, and verifying it once |
  | S-0400 | Landing a green batch through today's acceptance; bisecting a red one; sending the failing story back with the failing step and a thread to its agent |
  | S-0401 | Settling a rebase that stops inside the shared paths within the batch; telling the story's agent of a conflict outside them |
  | S-0402 | Landing while the main checkout has uncommitted files no queued story claims, and telling the operator what the checkout is behind |
  | S-0403 | The orchestrator enqueueing each story it verified instead of accepting it, with evidence |
  | S-0404 | The queue wait per story, the batch size, and the bisections per week, in `metrics.md`, `flai stats`, and the dashboard |
  | S-0405 | The close-out and `flai verify` leaving the test tiers the manifest marks `queue` to the batch |

## Consequences

- The operator sees which stories would batch, in which order, and why each other story would not, before anything is verified.
- The operator turns the queue on or off with one host action. With it off, acceptance is as before.
- A story waits in the queue only for the batch's run, not for another story's conflict or for the operator's checkout.
- The orchestrator stops accepting stories one by one and enqueues them instead (S-0403). Its turns never wait on the test tiers.
- A story that undoes main's change still waits for the operator, as ADR-0155 requires. The queue shows why.
- Under pull-request acceptance the queue gives nothing. The forge's CI verifies each story, and each acceptance waits for it.
- The queue's state is on one host. A clone on another host has no queue, as it has no verify result.
- A later ADR, from S-0405, refines [ADR-0110](0110-a-story-s-agent-runs-the-close-out-which-runs-flai-verify-itself-before-review.md)'s close-out for the tiers left to the batch.

## Alternatives considered

- **The orchestrator runs each batch through an MCP tool.** Rejected: its turns would wait on the integration and smoke tiers.
- **No cap on batch size, every enqueued story in one batch.** Rejected: a red batch takes more runs to bisect, and one failure holds back more stories.
- **Batch a story that undoes main's change, and let its landing stop.** Rejected: it wastes the batch's run, since only the operator can keep the removal.
- **Enqueue every eligible story by itself.** Rejected: a story enters the queue only on someone's word, the orchestrator's evidence or the operator's.
- **Run the queue under pull-request acceptance too.** Rejected: every rebased story would stop on `forge_stale`.
