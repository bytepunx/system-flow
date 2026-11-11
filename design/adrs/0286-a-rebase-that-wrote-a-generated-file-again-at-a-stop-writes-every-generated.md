---
id: ADR-0286
title: "A rebase that wrote a generated file again at a stop writes every generated file again at the branch's head once it finishes, and commits what differs from the head's"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0098, ADR-0228]
---

# ADR-0286 A rebase that wrote a generated file again at a stop writes every generated file again at the branch's head once it finishes, and commits what differs from the head's

## Context

[ADR-0098](0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md) writes `design/issues/summary.md` again when a rebase stops on it alone, from the issue files as the worktree is at that stop. [ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md) adds the ADR index to the generated files, written again the same way, and has `flai check` report `issues.summary` and `adr.index` when either differs from what its folder renders.

A file written at a stop is right for that commit only. A later commit of the branch that adds an issue or an ADR file, and leaves the generated file alone, replays cleanly over it, so the branch's head ends with a generated file its folder does not render. [I-0266](../issues/I-0266-flai-check-finds-issues-summary-outside-the-story-at-close-out.md) is that case. S-0520's sync stopped on `summary.md` at `66a374f21`, which closed I-0218, and wrote it again there. The next commit, `44beb162d`, added I-0264's file and replayed cleanly. The branch was accepted with a summary that lacked I-0264's row. Every close-out after the acceptance found `issues.summary` outside its story until S-0576's close-out wrote the summary again.

The merge queue's batch rebase settles its stops with the same `ContinueOverResolved`, so a batch can land the same stale file.

## Decision

**A rebase that wrote a generated file again at a stop writes every generated file again at the branch's head once it finishes, and commits what differs from the head's.** This refines ADR-0098 and ADR-0228; the rest of each stands.

1. **When.** Only after a rebase that stopped and that `ContinueOverResolved` continued to its end, writing at least one generated file again on the way. A clean rebase writes nothing at a stop and leaves each generated file as the branch's own commits made it; the fold after a clean rebase already writes them again when it folds (ADR-0126).
2. **What.** Each generated file `issues.Generated` names, the issues summary and the ADR index, written from the files at the branch's head.
3. **What counts as a difference.** A file whose new content equals the head's, the summary's `updated:` line left out as the `issues.summary` check leaves it out, is put back as the head has it. The one comparison is shared by the check and the rebase, so the two never disagree on what is stale.
4. **The commit.** What differs is committed on the branch as `docs: [S-nnnn] write the generated files again at the branch's head`. The sync's result lists those paths, as it lists the files written again at a stop.
5. **Where.** In `flai stream sync`'s rebase, which `flai task done`, the close-out, and `flai accept` run, and in the merge queue's batch rebase, on the story's rebased copy in the batch's worktree, never on the story's own branch.
6. **Failure.** A file written again at the head but not committed is an error naming the paths to add and the subject to commit them with, as the fold and the changelog commit say (ADR-0126, ADR-0257).

## Consequences

- A branch whose sync stopped on a generated file reaches main with that file as its folder renders it, so no close-out after its acceptance finds `issues.summary` or `adr.index` outside its story for this cause.
- A sync that stopped on a generated file may add one commit to the branch. One that wrote nothing different adds none.
- A generated file the branch left stale by other means, such as a fold that failed after it wrote some files, is not covered: this runs only after a rebase that stopped.

## Alternatives considered

- **Write the generated files again at the head after every rebase.** It would also mend a stale file the branch brought by hand, but it commits on branches whose rebase met no stop, and the check already reports such a file inside the story, where its agent fixes it.
- **Write the generated files again at each stop, and again at every later commit the rebase replays.** git gives no hook between the commits of a clean replay without running the rebase one commit at a time, which is slower and more code for the same head.
- **Leave it to `flai check --fix` at the close-out.** The close-out of the story whose sync made the stale file passes, since its own check runs on its branch before the acceptance; the finding falls on every other story after it.
