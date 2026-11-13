---
id: ADR-0299
title: "A close-out records no story.unaccepted, and a check scoped to a story leaves out every story.unaccepted outside it"
status: accepted
date: 2026-11-13
supersedes: []
superseded_by: []
refines: [ADR-0085]
---

# ADR-0299 A close-out records no story.unaccepted, and a check scoped to a story leaves out every story.unaccepted outside it

## Context

[ADR-0085](0085-a-close-out-s-flai-check-reports-findings-outside-the-story-as-notes-and.md) scoped a close-out's `flai check` to its story and recorded each rule's findings outside the story in one issue per rule. Later ADRs took out of that recording the findings the closing story's agent cannot clear, among them [ADR-0137](0137-a-close-out-records-no-threads-archived-and-a-check-scoped-to-a-story-leaves.md) a `threads.archived`, [ADR-0256](0256-a-close-out-records-no-narrative-missing-and-a-check-scoped-to-a-story-leaves.md) a `narrative.missing`, and [ADR-0296](0296-a-close-out-records-no-item-blocked-closed-and-a-check-scoped-to-a-story-leaves.md) an `item.blocked-closed` outside the story.

[I-0371](../issues/I-0371-flai-check-finds-story-unaccepted-outside-the-story-at-close-out.md), ``flai check finds `story.unaccepted` outside the story at close-out``, has six instances, from the close-outs of S-0778, S-0730, S-0777, S-0780, S-0779, and S-0783 on 2026-11-13. All have the same cause.

1. **The finding names another story.** Each is on the archived S-0001, S-0002, S-0003, and S-0005: "S-0001 is done but its branch story/S-0001 was never merged; merge or delete it". `story.unaccepted` warns on a done story that was never archived, and on a done story whose `story/S-nnnn` branch still exists unmerged. None of the closing stories touched those four.
2. **The branches were a test's leak.** All four point at 5c3b47d0c, a commit of S-0777's from before its rebase. S-0777 ran the Go tests with `TMPDIR` inside its worktree. A test project with no `git init` of its own let git walk up into that worktree, and story start added real `story/S-nnnn` branches and worktrees to the main repository. S-0781 stops that leak with a git ceiling in the tests' `TestMain`, and the operator removes the four branches.
3. **Only the main checkout clears it.** A done story not yet archived is accepted with `flai accept`, and a branch left after acceptance is a ref of the main repository, which the operator merges or deletes. No story branch archives another story or deletes another story's ref, so the closing story's agent can fix nothing in its worktree.
4. **It cannot arise inside the story.** A story at its close-out is in progress, not done, so its own item file never carries the finding.
5. **Recording it adds nothing.** Each close-out bumped I-0371 with the same four findings, which counts how many close-outs ran while the branches stood, not a cost any closing story could remove.

## Decision

A close-out records no `story.unaccepted` in an issue: a check scoped to a story leaves out every `story.unaccepted` outside the story, as ADR-0137 does with a `threads.archived`, ADR-0256 with a `narrative.missing`, and ADR-0296 with an `item.blocked-closed`.

1. **A check scoped to a story leaves it out.** `flai check --story S-nnnn` drops from its result every finding of rule `story.unaccepted` whose path is not inside the story, and takes back the counts it added. The finding's path is the done story's item file, so another story's is outside.
2. **One inside the story stays.** A `story.unaccepted` on the story's own item file, or on an item file its branch changes or that is uncommitted in its worktree, is inside the story and is reported as before.
3. **`--record-issues` records no `story.unaccepted` outside the story.** What is left out of the result is not recorded.
4. **Without `--story` every `story.unaccepted` is reported**, so the main checkout's `flai check` still warns on each done story not accepted and each branch left unmerged, naming `flai accept` or the branch to merge or delete.

This refines ADR-0085, whose rule that every finding outside the story is a note recorded in an issue now leaves out `story.unaccepted`, and extends the treatment of ADR-0137, ADR-0256, and ADR-0296 to one more finding no story branch clears.

## Consequences

- A close-out no longer bumps an issue because a story it did not touch is done and unaccepted, or left a branch, and I-0371 can be closed once the leak is stopped and the four branches are gone.
- The closing story's agent no longer sees, at its close-out, a finding only `flai accept` or the operator's merge or deletion of a branch in the main checkout clears. The operator still sees it in the main checkout's `flai check`.
- A leak like S-0777's is no longer seen at close-outs, only in the main checkout's check. The git ceiling in the tests is what stops the leak itself.
- An open issue for `story.unaccepted` recorded by an older flai stays as it is until it is closed.

## Alternatives considered

- **The root fix alone.** Stop the leak and remove the four branches, and change nothing in the check. That clears every instance today, and S-0781 does both. But a real done story left unaccepted, or a branch an older flai or a hand acceptance left behind, would again reach every close-out, which can do nothing about it, and each would bump an issue until the operator acted in the main checkout.
- **Keep it as a note but record none**, as [ADR-0115](0115-a-close-out-records-no-wip-overlap-and-notes-only-an-overlap-naming-its-story.md) does with a `wip.overlap` naming the story. The note would name another story that no close-out can act on, and the main checkout's check already reports it, where the operator sees it.
