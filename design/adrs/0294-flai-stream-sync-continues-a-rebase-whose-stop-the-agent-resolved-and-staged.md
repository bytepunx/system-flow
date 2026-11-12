---
id: ADR-0294
title: "flai stream sync continues a rebase whose stop the agent resolved and staged, settling the stops after it as its own rebase does, and every stop's advice says git add and sync again"
status: accepted
date: 2026-11-12
supersedes: []
superseded_by: []
refines: [ADR-0069, ADR-0098, ADR-0126, ADR-0155, ADR-0286]
---

# ADR-0294 flai stream sync continues a rebase whose stop the agent resolved and staged, settling the stops after it as its own rebase does, and every stop's advice says git add and sync again

## Context

`flai stream sync` settles some stops of its rebase itself: it writes `design/issues/summary.md` and the ADR index again ([ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md), [ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md)), merges an issue file by its instances ([ADR-0126](0126-a-rebase-of-a-story-branch-merges-an-issue-file-both-sides-changed-by-its.md)), and merges a design, docs, or shared markdown file by sections ([ADR-0175](0175-a-rebase-that-stops-on-a-markdown-file-inside-the-manifest-s-shared-paths.md), [ADR-0276](0276-a-rebase-that-stops-on-a-markdown-file-under-the-manifest-s-design-or-docs.md)). It does so only inside the rebase it started, through `ContinueOverResolved`. A stop it cannot settle is left for the agent, and every message said to resolve it, `git add` it, and run `git rebase --continue`, then sync again. A sync run on a rebase in progress refused it (`StopRebaseInProgress`).

[I-0351](../issues/I-0351-a-git-rebase-continue-run-by-hand-after-a-sync-s-stop-meets-later-stops-on.md) records the cost. S-0689's sync stopped on `design/system/flai-cli.md`, which the section merge could not settle. The agent resolved it and ran `git rebase --continue` by hand, as told. That replayed the next commit, which stopped on `design/issues/I-0347` and `summary.md`, stops the sync would have settled itself. The agent settled them by hand. And the steps after a clean rebase, the generated files written again at the head (ADR-0286), the changelog lines put back (ADR-0257), the undone check (ADR-0155), the fold, and the fix, never ran for that rebase.

## Decision

**`flai stream sync`, run on a story's worktree where its rebase waits on a stop the agent resolved and staged, continues that rebase itself, settles each later stop as its own rebase does, and finishes with the steps after a clean rebase. Every stop's advice says to `git add` the resolved paths and run `flai stream sync` again.** This refines ADR-0069: the agent still never starts, continues, or merges a rebase by hand; `git rebase --abort` stays the way back.

1. **Resolved.** A stop is resolved when no path is unmerged (`git diff --diff-filter=U`), no staged file holds a conflict marker (`git diff --cached --check`), nothing tracked is changed and unstaged, and the rebase is of the story's own branch (the rebase state's `head-name` is `refs/heads/story/S-nnnn`). Anything else is refused as before, `StopRebaseInProgress`, touching nothing, naming the unmerged, marked, or unstaged paths.
2. **The continue.** The sync runs `git rebase --continue` with no editor, then the loop its own rebase runs, `ContinueOverResolved`, with the same files, so that later stops on generated files, issue files, and markdown merged by sections are settled. A later stop it cannot settle stops the sync with `StopConflicts`, as a fresh rebase does.
3. **The old base and head.** The steps after a clean rebase that read the branch as it was before the rebase take it from the rebase's own state in the worktree's git directory, read before the continue: `onto` is the commit the rebase was onto, the old head is `orig-head`, and the old base is `git merge-base onto orig-head`. Those steps run against that `onto`, not main's head now, so that a commit main gained since the rebase started is none of the branch's (I-0338). A rebase state that lacks either file is refused, naming `git rebase --abort` and a fresh sync.
4. **After it.** Once the rebase is clean, the sync runs what follows a clean rebase: the generated files again at the head (ADR-0286), the changelogs (ADR-0257), the undone check (ADR-0155), the fold (ADR-0126), and the fix, then the trial merge and the checks, as any sync. The result says `resumed: true`. When main moved past `onto` meanwhile, the branch does not contain it yet; the next sync, which every task close and the close-out run, rebases onto it.
5. **Where.** In `storygit.Sync`, which `flai stream sync` and `flai task done` run. `flai task done`, called again after its sync stopped, passes a resolved stop to its sync step instead of refusing it at its commit step, which commits nothing while the rebase waits. `storygit.Rebase` run alone, as `flai accept` and the merge queue run it, still refuses a rebase in progress: an acceptance takes the branch as committed.
6. **The advice.** Every message, help text, prompt, and convention that told an agent to run `git rebase --continue` after a sync's stop says instead to resolve each path, `git add` it, and run `flai stream sync` again, or call `flai task done` again, or ask for the close-out again, where that is what stopped.

## Consequences

- A stop the agent resolves is settled with the same merges as one flai settles, and the steps after a clean rebase always run, so a hand-continued rebase no longer leaves issue files, `summary.md`, the ADR index, or a changelog for the agent, nor skips the undone check.
- The agent learns no new command: the advice already ended "then run flai stream sync again".
- A rebase continued by hand with `git rebase --continue` still works, as before, but loses what the sync would have settled; nothing forbids it beyond the conventions.
- A sync that resumes may leave the branch behind main's newest commits; the next sync takes them.

## Alternatives considered

- **A `--continue` flag on `flai stream sync`.** One more thing to learn, and the advice would differ from the plain sync the agent runs everywhere else. A plain sync of a resolved stop has one meaning.
- **Settle the stops in a git hook or merge driver.** Each clone must configure it, which flai does not install (ADR-0098), and it would not run the steps after a clean rebase.
- **Abort and rebase again from the start.** It throws away the agent's resolution, and the same stop comes back.
- **Record the old base and head in a file of flai's when the sync starts.** git's rebase state already holds them, for any rebase of the branch, including one an acceptance left.
