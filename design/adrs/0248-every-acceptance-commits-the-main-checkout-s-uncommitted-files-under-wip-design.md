---
id: ADR-0248
title: "Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, in a commit of their own before the merge when the story's change would overwrite one, leaves the rest as they are, and sends a story whose rebase it cannot settle back to its agent, whoever accepts"
status: accepted
date: 2026-11-10
supersedes: []
superseded_by: []
refines: [ADR-0175, ADR-0178, ADR-0189]
---

# ADR-0248 Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, in a commit of their own before the merge when the story's change would overwrite one, leaves the rest as they are, and sends a story whose rebase it cannot settle back to its agent, whoever accepts

## Context

S-0566 records what the operator meets at nearly every acceptance from the board. Two things stop the acceptance by hand and leave the operator to do git's work:

- **A rebase that stops.** `flai accept` rebases the story's branch onto the main branch before it merges. When the rebase stops on conflicts that the merge of the issue files, the ADR index, the issue summary, and the shared markdown files do not settle ([ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md), [ADR-0126](0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md), [ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md), [ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md)), the acceptance by hand leaves the stopped rebase in the story's worktree and tells the operator to resolve each path, `git add` it, and `git rebase --continue`, or to abort. [ADR-0189](0189-the-orchestrator-s-acceptance-lands-on-a-dirty-main-as-the-merge-queue-s-does.md) gave the orchestrator's acceptance a better answer: abort the rebase, send the story back to in-progress, open a thread naming the paths, and let `flai serve` start the story's agent on it. It kept the stop in place for the acceptance by hand, on the ground that the operator at the keyboard may want to settle it. In practice the operator is asked to resolve conflicts in files they did not write, in a story they did not work, and the story waits in review until they do.
- **Uncommitted files in the main checkout.** flai itself writes in the main checkout between landings: transitions and narratives under `wip/`, the digest, the analysis reports, the issue files and their summary, and the ADR index under `design/`. The acceptance commit takes `wip/` as a matter of course, but every other uncommitted file is "outside wip": the acceptance by hand refuses it unless `--yes` includes it, and the dashboard asks the operator, with a checkbox, whether to include the files in the acceptance commit. The orchestrator's acceptance and the merge queue's landing leave those files as they are ([ADR-0178](0178-the-merge-queue-lands-each-story-in-the-main-checkout-with-git-s-own-check.md), ADR-0189), and refuse or hold back a story whose change would overwrite one, naming `git stash push` as the way to free it. A dirty `design/issues/summary.md` or digest entry, which flai wrote, thus asks the operator a question at every acceptance, and holds a story back from the queue.

## Decision

Every acceptance commits the main checkout's uncommitted files under wip, design, and the shared paths, in a commit of their own before the merge when the story's change would overwrite one, leaves the rest as they are, and sends a story whose rebase it cannot settle back to its agent, whoever accepts.

- **The shared set.** The main checkout's uncommitted and untracked files under the manifest's `wip` folder, under its `design` folder, and under or at each path in `claims.shared` are the *shared set*: what flai and the agents write there between landings is flai's own work, which the acceptance commit exists to land. Every acceptance that lands a story here commits them, by hand, from the dashboard, by the orchestrator, and as the merge queue's landing alike:
  - with the acceptance commit, as `wip/` was committed before, when the story's change, `main...story/S-nnnn`, would overwrite none of them;
  - otherwise in a commit of their own on the main branch before the story's branch is rebased and merged, `chore: [S-nnnn] commit the main checkout's wip and design before accepting`, so that the rebase settles each such path by its own rules (the issue merge, the summary and index regeneration, the section merge of a shared markdown file) in place of git refusing to overwrite it. The lint fix of the tiers that carry one runs over them before that commit, as it runs before the acceptance commit.
- **The rest.** An uncommitted file outside the shared set is the operator's, and every acceptance leaves it as it is, out of the lint fix and out of its commits, as the orchestrator's and the queue's did. An acceptance is refused, before anything is merged, only when the story's change would overwrite one of them, as the blocker `dirty_main`, naming the files and the command that sets them aside, for the acceptance by hand and the dashboard as for the orchestrator's; the queue holds such a story back as before. `flai accept --yes` includes those files in the acceptance commit, for the shell; the dashboard offers no choice, since nothing in it needs one.
- **The preview.** `flai accept --dry-run` and the acceptance preview list the shared set the acceptance will commit as `shared_uncommitted`, and the files left as they are as `uncommitted`, and the merge queue's eligibility, preflight, and landing count a file in the shared set as clean.
- **A rebase that stops.** When the rebase of any acceptance that lands here, by hand, from the dashboard, or by the orchestrator, stops on conflicts after the merge of the issue files, the generated files, and the shared markdown files, `flai accept` does what ADR-0189 had it do for the orchestrator alone: aborts the rebase, so the branch and the main branch are as they were; records the conflict; moves the story back to in-progress as `flai`, with a reason naming the paths; opens one thread on the story, by `flai`, naming the paths and the branch, or writes on the one still open from an earlier send back; and answers that the story was sent back, with `--json` as the refusal `{"refused": {"id", "reason", "sent_back", "branch", "base", "conflicts", "thread"}}`. `flai serve` starts the story's agent on that thread, resuming its session where it has one, and the agent syncs, resolves, closes out, and moves the story to review again. The merge queue's landing keeps its own send back ([ADR-0176](0176-the-merge-queue-lands-each-story-of-a-passing-batch-through-flai-accept-batch.md)).
- **A branch carrying conflict markers** is refused before the merge as before, and the blocker says to send the story back to in-progress so that its agent resolves them, not to resolve them by hand.
- No acceptance tells the operator to rebase a story's worktree, to resolve a conflict, or to choose which files to include.

## Consequences

- The operator accepts from the board and is never left mid-rebase: a conflict reaches the story's agent, which is the one party that knows both sides of the change.
- An acceptance may make two commits on the main branch where it made one: the shared set first, then the acceptance. The release evaluation reads both.
- The dashboard's acceptance confirmation and the review page lose the checkbox for uncommitted files and show the shared set as what the acceptance commits.
- ADR-0189's alternative "abort and send back on every acceptance", which lost then, is adopted: the operator asked for it after a fortnight of resolving other agents' conflicts.
- A story sent back from the operator's acceptance moves out of review without the operator's say; the thread and the narrative name the paths, and the story returns to review through its agent's close-out.

## Alternatives considered

- **Keep the stop in place for the acceptance by hand** (ADR-0189). Lost: that is the intervention S-0566 removes.
- **Resolve the conflicts in `flai accept` with a model.** Lost: the story's agent already holds the story's context and its worktree; starting it on a thread costs nothing new.
- **Commit every uncommitted file of the main checkout with the acceptance.** Lost: a half-done code edit of the operator's would land under a story's acceptance.
- **Commit the shared set always in a commit of its own.** Lost: the common case, no overlap, needs no second commit, and a history of one commit per acceptance stays readable.
- **Take `docs` into the shared set.** Lost: nothing of flai's writes under `docs` in the main checkout; a project that wants a docs path committed names it in `claims.shared`.
