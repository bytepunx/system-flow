---
id: ADR-0189
title: "The orchestrator's acceptance lands on a dirty main as the merge queue's does, finishes an acceptance of its own that stopped before its commit, sends a story whose rebase it cannot settle back to its agent, and is held by a thread of its own only until it resolves it"
status: proposed
date: 2026-10-23
supersedes: []
superseded_by: []
refines: [ADR-0093, ADR-0175, ADR-0176, ADR-0178]
---

# ADR-0189 The orchestrator's acceptance lands on a dirty main as the merge queue's does, finishes an acceptance of its own that stopped before its commit, sends a story whose rebase it cannot settle back to its agent, and is held by a thread of its own only until it resolves it

## Context

With `accept_reviews` on, the orchestrator accepts a story in review with `flai accept <S-nnnn> --by orchestrator` ([ADR-0093](0093-with-accept-reviews-on-the-orchestrator-accepts-a-story-in-review-through-flai.md)). [I-0150](../issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md) found three stories that sat in review for 259 hours together in W41. Each waited because the orchestrator's acceptance stopped on something only the operator could clear:

- **A dirty main checkout (TH-0389, TH-0391, TH-0394).** `flai accept` refuses any uncommitted file outside `wip/` unless `--yes`. The orchestrator may not pass `--yes`, because that would commit the operator's files with the acceptance. S-0338's acceptance had also stopped before its commit, and `flai accept` refuses to let the orchestrator finish an acceptance it made itself.
- **A rebase on a shared path that it cannot settle (TH-0397).** [ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md) merges a shared markdown file by its sections. A stop that survives that merge leaves the story's worktree mid-rebase. `flai guard` refuses the orchestrator any git write, so it can neither abort the rebase nor send the story back. The story's agent, the one party that can resolve the conflict, is not started.
- **A gate it raised itself (TH-0387).** Any open thread on a story blocks the orchestrator's acceptance, whoever opened it. The orchestrator opened a thread asking for a gate that was never built, and the story waited for the operator to resolve it.

Since [ADR-0178](0178-the-merge-queue-lands-each-story-in-the-main-checkout-with-git-s-own-check.md), the merge queue already lands on a dirty main. Since [ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md), it aborts a rebase it cannot settle and sends the story back. The orchestrator's acceptance does neither.

On TH-0598 the planner proposed one remedy per cause, and the orchestrator recommended the first two as proposed. For the third, it recommended a stricter form than the planner's.

## Decision

The orchestrator's acceptance lands on a dirty main as the merge queue's does, finishes an acceptance of its own that stopped before its commit, sends a story whose rebase it cannot settle back to its agent, and is held by a thread of its own only until it resolves it.

- **A dirty main.**
  - Under `--by orchestrator`, `flai accept` does not refuse the main checkout's uncommitted files outside `wip/`. It reads them before anything changes and leaves them as `--batch` does: out of the lint fix and out of the acceptance commit.
  - It refuses only a story whose change, `main...story/S-nnnn`, would overwrite one of those files or an untracked one. The dry run lists that as the blocker `dirty_main`, naming the files and the command that frees them, as the queue's eligibility does.
  - The operator's acceptance by hand keeps its refusal and `--yes`.
- **An acceptance of its own that stopped.**
  - The orchestrator may finish an acceptance whose move to done it made, when that acceptance stopped before its archive or its commit: `flai accept <S-nnnn> --by orchestrator` resumes it as the operator's resume does.
  - A stopped acceptance that anyone else made stays the operator's, and the refusal names them.
- **A rebase it cannot settle.**
  - When the rebase of an orchestrator's acceptance stops after the section merge and the generated-file steps, `flai accept` runs `git rebase --abort` in the story's worktree. The branch and main are then as they were.
  - It records the conflict, as today.
  - It moves the story back to in-progress as `flai`, with a reason naming the conflicting paths.
  - It opens one thread on the story, by `flai`, naming the paths and the story's branch.
  - It answers with an error that says the story was sent back.
  - `flai serve` starts the story's agent for that thread, as it does for a story the queue sent back, resuming its session where it has one.
  - The operator's acceptance by hand still leaves the stop in the worktree for them to finish or abort.
- **A thread of its own.**
  - A thread on the story or one of its tasks that the orchestrator opened, and that nobody else wrote on, still holds the story. The orchestrator resolves it with a reason before it accepts.
  - The dry run lists such a thread under its own blocker, `own_thread_open`, apart from `thread_open`, and says the orchestrator resolves it.
  - A thread that anyone else opened or wrote on holds the story until it is resolved. The orchestrator resolves only a thread it opened, as before.
  - The orchestrator accepts on flai's blockers and the story's criteria. It opens a thread on a story in review only for a reason not to accept it, and resolves that thread once the reason is gone.

## Consequences

- The three waits I-0150 recorded no longer need the operator:
  - A file uncommitted on main that the story does not change stops nothing.
  - A conflict reaches the story's agent.
  - The orchestrator's own concern is the orchestrator's to clear.
- A concern the orchestrator raises still holds the story, and is resolved with a reason that the thread keeps. It is not lost at archive.
- A send-back moves a story out of review without the operator. The thread and the narrative name the paths, and the story returns to review through its agent's close-out.
- `flai accept --by orchestrator` writes git in the story's worktree, through `git rebase --abort`. The orchestrator itself still runs no git write: `flai guard` is unchanged.
- The orchestrator's dry run reads `main...story/S-nnnn` and the checkout's status once more.

## Alternatives considered

- **A thread of its own does not hold the story, and the acceptance resolves it at archive.** This was the planner's form. It lost on the orchestrator's recommendation: a concern the orchestrator raised could then never stop a story, as TH-0589's hold on S-0483 rightly did.
- **The orchestrator may pass `--yes`.** Lost: that commits the operator's uncommitted files with the acceptance.
- **Let `flai guard` grant the orchestrator `git rebase --abort` and `flai move`.** Lost: it widens what a strategic agent may write, when flai can make the one safe write itself.
- **Abort and send back on every acceptance.** Lost: the operator at the keyboard can settle the stop in place, and may want to.
- **Keep the story in review and open a thread for the operator.** Lost: that is the wait I-0150 records.
