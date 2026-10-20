---
id: ADR-0178
title: "The merge queue lands each story in the main checkout with git's own check while the checkout has uncommitted files, so it never leaves the checkout behind main, and a story whose change would overwrite one of them is held back, enqueued, until they are clean"
status: accepted
date: 2026-10-20
supersedes: []
superseded_by: []
refines: [ADR-0176]
---

# ADR-0178 The merge queue lands each story in the main checkout with git's own check while the checkout has uncommitted files, so it never leaves the checkout behind main, and a story whose change would overwrite one of them is held back, enqueued, until they are clean

## Context

[ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md) lands each story of a passing batch with `flai accept <story> --by <who enqueued it> --batch <batch>` in the main checkout. That is the acceptance of today. It refused any uncommitted file outside `wip/` unless `--yes`, and with `--yes` it ran the lint fix over those files and committed them with the acceptance.

The main checkout is where the operator and the analyzer work, and their files are often uncommitted there. On TH-0462, S-0391 waited a day because the analyzer's report, its issues, and a design system's files were uncommitted on main, though none of them were S-0391's ([I-0150](../issues/I-0150-a-story-in-review-waits-days-for-the-operator-when-the-orchestrator-s-acceptance-stops-on-a-dirty-main-a-rebase-on-a-shared-path-or-a-gate-it-raised-itself.md)).

The acceptance merges the story's branch with `git merge --ff-only` in the main checkout. That merge already carries git's own check. It refuses only when an uncommitted or untracked file in the checkout would be overwritten by the change, and otherwise it updates every file the change touches while leaving the rest as they are. So a landing either moves main and the checkout's files together, or it does not move main at all.

The question is how a batch lands while the checkout is dirty, and what the checkout is left as. The operator chose on TH-0521.

## Decision

The merge queue lands each story in the main checkout with git's own check, so it never leaves the checkout behind main, and it holds back, still enqueued, a story whose change would overwrite one of the checkout's uncommitted files, until they are clean.

- **The landing stays in the main checkout, through `flai accept --batch`.**
  - A landing never leaves the checkout's files behind main, since the merge moves the branch and the checkout's files together.
  - Under `--batch`, `flai accept` does not refuse the main checkout's uncommitted files outside `wip/`. It reads them before anything changes and leaves them as they are: out of the lint fix over the checkout's uncommitted files, and out of the acceptance commit, which holds only the acceptance.
  - `flai accept` by hand keeps its refusal, and `--yes` to include the files.
- **A story whose change would overwrite one of the checkout's uncommitted or untracked files is held back.**
  - Its change is what its branch changes from the commit it shares with main (`main...story/S-nnnn`).
  - Before the landing:
    - the merge queue's eligibility gives such a story in review the reason `dirty_main`, naming the files and the command that frees them: commit them, or `git stash push -- <files>`;
    - the queue does not batch it while the reason holds;
    - `flai queue`, `flai board`, and the dashboard show the reason, as they show every reason, and it clears once the files are clean.
  - During the landing, when a file of a batched story became dirty after the batch was built:
    - the landing checks each story just before it lands it, and holds that story back;
    - it records the story and the files in the landing as held back, and as the story's outcome, `held_back`;
    - it keeps the story enqueued and lands the rest in queue order;
    - the next batch takes the held story once the files are clean.
- The stories landed after a held-back one were verified in a batch that held it. Each is rebased onto main as it lands, as any acceptance is. A rebase that stops fails that landing, which drops the batch as ADR-0176 says, and the stories not landed are batched again.

## Consequences

- An uncommitted file that no queued story changes stops no landing, and the operator need not commit or stash first.
- The checkout is never left behind main, so there is nothing to catch up. The story's criterion 4, which assumed a checkout left behind, was reworded on TH-0521 to report the stories a dirty file holds back.
- The operator sees each held story, its files, and the command that frees them, until they are clean. A held story is not built into batch after batch while the files stay dirty.
- Eligibility now reads the main checkout's uncommitted files once per computation and, when there are any, what each story in review changes. That adds one `git status` and one `git diff` per story in review to each read of the queue.
- A story landed after a held-back one stands on a commit the batch did not verify exactly, without the held story beneath it. It was verified on its own branch and in the batch, and its own acceptance rebases it onto main.

## Alternatives considered

- **Land through the batch's worktree and move main's branch ref there, leaving the checkout behind.** Lost: HEAD moves under the checkout's index. `git status` then shows the landing reversed as staged changes, and the operator's next `git commit -a` would undo it.
- **Move main's branch with git's check from the batch's worktree.** This is the merge the acceptance already makes in the main checkout. Making it from elsewhere gains nothing and bypasses the acceptance of today.
- **Update the checkout's files around the dirty ones by hand.** Lost: it reimplements git's two-way merge check, less safely.
- **Stash the dirty files, land, and restore them.** Lost: the stash stack is shared by every worktree and session (`git.md`), and a restore that conflicts leaves the operator's work in a stash they did not make.
- **Keep refusing, and ask the operator to commit or stash first.** Lost: that is the day's wait on TH-0462 that this decision removes.
