---
id: ADR-0284
title: "The touches drift leaves out the folders flai writes, wip and the design folder's digest, analysis, and issues, on both sides, and flai task done narrows the closing task's touches, and the story's, to the files the story's branch changes"
status: accepted
date: 2026-11-11
supersedes: []
superseded_by: []
refines: [ADR-0113, ADR-0144]
---

# ADR-0284 The touches drift leaves out the folders flai writes, wip and the design folder's digest, analysis, and issues, on both sides, and flai task done narrows the closing task's touches, and the story's, to the files the story's branch changes

## Context

[I-0286](../issues/I-0286-declared-touches-match-what-a-story-changes-in-16-percent-of-completed-stories.md) found that a story's declared touches match what it changed in 16 percent of completed stories, as `flai stats` counts `claims.weeks[].exact_share` ([ADR-0113](0113-flai-stats-reports-held-stories-per-day-held-time-by-reason-per-week-and-the.md)). S-0674 read `claims.drift[]` for the 115 stories from S-0560 on:

| Measure | Stories |
|---------|---------|
| exact as flai counts it | 25 |
| with a file `outside` the touches | 80 |
| every `outside` file under `design/digest/`, `design/issues/`, or `design/analysis/` | 80 |
| exact once those folders and `wip/` are left out of both sides | 97 |
| still inexact: 1 to 4 predicted files the branch never changed | 18 |

Two causes, neither a claim the agents missed:

- **Attribution.** `storygit.Committed` gives a story every file of each commit whose subject names it in brackets, leaving out only the `wip` folder (S-0205). Three of those commits are flai's bookkeeping. `chore: [S-nnnn] accept and archive` sweeps the main checkout's uncommitted digest entries into the acceptance commit ([ADR-0248](0248-every-acceptance-commits-the-main-checkout-s-uncommitted-files-under-wip-design.md)). `chore: [S-nnnn] create story` writes the issue the story came from. `docs: [S-nnnn] record I-nnnn` is a close-out recording a finding outside the story. flai writes these folders as it writes `wip/`. On the other side, a touch under `wip/` can never be matched, since its commits are left out, and is counted `unchanged`.
- **Predictions never dropped.** A task's touches are a guess made before its code is read ([delegation.md](../conventions/delegation.md)). [ADR-0144](0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md) has `flai task done` widen them with what the commit changed, and never remove one. A predicted file the work never changed stays in the task's and the story's touches, and in the claim, for good.

A claim that names files the story will not change holds ready stories it need not hold, sends trial merges and conversations about paths nobody changes, and lowers `exact_share`, which then reports neither.

## Decision

The touches drift leaves out the folders flai writes, on both sides, and `flai task done` narrows the closing task's touches, and the story's, to the files the story's branch changes.

1. **flai's folders.** They are the manifest's `wip` folder and three folders of its `design` folder: `digest`, `analysis`, and `issues`. One function names them from the manifest, and the drift reads them from it.
2. **The commits side.** `storygit.Committed` leaves out every file under flai's folders, as it left out `wip` alone before. What it keeps between processes is keyed by the folders it left out, so an answer read under the old rule is read again.
3. **The touches side.** The drift leaves out of a story's claim each touch that is one of flai's folders or lies under one, before it counts `outside` and `unchanged`. A story whose touches lie only in flai's folders, or whose committed files lie only there, has no entry in `claims.drift[]`, as a story with no committed file outside `wip` had none before.
4. **Narrowing at task done.** After the touches step widens ([ADR-0144](0144-flai-task-done-widens-touches-only-with-the-committed-paths-the-story-s-branch.md)), `flai task done` removes from the closing task's touches each touch that names a file the story's branch does not change against its merge base with the main branch, read from the same `git diff --name-only --no-renames <main>...HEAD`. A touch that is a sub-project's name or tag, or a folder in the worktree, is never removed: it is not a file, and a story's folder touch is the planner's to keep ([ADR-0188](0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md)).
5. **The story's touches.** Each touch removed from the task is removed from the story's touches too, when the story names it and no other open task of the story names it. A story's touch that the closing task did not name is left alone: it may still be a later task's prediction.
6. **Whether anything was committed.** Narrowing runs whenever the touches step runs, with a commit or without: a task closed with nothing to commit still drops the files it predicted and never changed.
7. **What it reports.** The result lists the touches removed, as `touches.narrowed`, with `task` and `story`, and the command prints them.

## Consequences

- `exact_share` measures the claims agents write, rather than flai's own bookkeeping commits. Measured over the stories already completed, it rises from the attribution fix alone; the narrowing shows in the stories completed after it.
- A story's claim shrinks as its tasks close. A ready story held on a file nobody changed starts once the task that predicted it closes.
- `exact_share` then measures the claim as `flai task done` leaves it, not as the plan first wrote it. How far the plan's guesses were off is no longer in the figure; the narrative's log keeps what each close removed.
- `flai accept`'s overlap notices and installs read `storygit.Committed` too, and no longer name a story's digest, analysis, or issue files. Those lie under the shared paths or the design folder and were left out of notices already ([ADR-0228](0228-flai-generates-the-adr-index-from-the-adr-files-and-checks-the-issues-summary.md)).
- ADR-0144's "it never removes a touch" no longer holds for a file the branch does not change once the task that named it closes.

## Alternatives considered

- **Leave out flai's commits by their subject**, `accept and archive`, `create story`, and `record I-nnnn`. Lost: subjects change between releases, and a story's own `docs: [S-nnnn] record I-nnnn` writes an issue the story does own. The folders are what flai writes, whatever the subject.
- **Leave out flai's folders in the drift only, keeping `storygit.Committed` as it is.** Lost: one definition of a story's change, read by the drift and the acceptance alike, is easier to keep true than two.
- **Narrow at the story's close-out, or at acceptance.** Lost: the claim would stay wide while the story is in progress, which is when it holds ready stories and starts conversations.
- **Narrow the story's touches to the branch's whole diff.** Lost: it removes predictions of tasks not yet written, and touches the agent declared on purpose.
- **Leave narrowing to the agent with `flai touches`.** Lost: the 18 stories show it is not done, as I-0121 showed for widening.
