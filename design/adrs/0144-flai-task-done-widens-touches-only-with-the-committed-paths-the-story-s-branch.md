---
id: ADR-0144
title: "flai task done widens touches only with the committed paths the story's branch still changes against main"
status: accepted
date: 2026-10-12
supersedes: []
superseded_by: []
refines: [ADR-0107, ADR-0128]
topics: [cli, touches]
---

# ADR-0144 flai task done widens touches only with the committed paths the story's branch still changes against main

## Context

[ADR-0107](0107-flai-task-done-closes-a-task-in-one-call-commit-sync-move-log-widen-touches.md) made `flai task done` close a task in one call, and its sixth step widens the task's and the story's touches with the paths the commit changed. [ADR-0128](0128-flai-task-done-commits-the-paths-the-closing-task-covers-and-those-no-other.md) kept that step and left out the paths the commit left for another task.

A commit can change a path back to the main branch's version. [I-0121](../issues/I-0121-flai-task-done-widens-touches-with-a-path-the-commit-reverted-to-main-re-growing-a-claim-the-story-narrowed.md) records it on S-0232: the story agreed with S-0340 to leave `scripts/flai-snapshot.sh` alone, removed it from T-1223's and the story's touches, reverted the file, and closed T-1223 again. The touches step added the path back to both, because the revert commit changed it, though the branch no longer differed from main there. The claim overlapped S-0340's again until the touches were narrowed by hand a second time.

A story's claim says what its branch will change on main when it is accepted. A path the branch does not change against main is not part of that change.

## Decision

The touches step of `flai task done` widens the task's and the story's touches only with the committed paths that the story's branch still changes against its merge base with the main branch.

1. **What it reads.** After the sync step, it reads the paths the branch changes with `git diff --name-only --no-renames <main>...HEAD` in the story's worktree. With `--no-renames`, both sides of a rename count as changed.
2. **What it widens.** It widens with each path of the commit that the branch still changes, as before: never the wip folder, never a path left for another task.
3. **What it reports.** The committed paths it leaves out are listed in the result as `unchanged`, and the command prints them, so the agent sees why a path it committed was not added.
4. **When it cannot read the diff.** The step stops, naming the error and the command to run again, as any other step stops.
5. **What it never does.** It never removes a touch. Narrowing the touches stays the agent's, with `flai touches`.

## Consequences

- A revert to main's version no longer re-grows a claim the story narrowed, and no overlap with another story comes back from it.
- A commit that changes a path and a later commit that reverts it leave the path in the touches the first one widened: only the agent's narrowing removes it.
- The step runs one more `git diff`, which reads names only.
- The tell step (S-0333) and the `--commit` writes of `flai adr new`, `flai issue new`, and the like are unchanged. Those writes add only content they wrote, and a revert is still worth telling a story whose claim covers the path.

## Alternatives considered

- **Widen with the branch's whole diff against main.** That would add paths earlier tasks committed and the agent narrowed on purpose, or that another task covers. The commit's paths, filtered, change only what this close committed.
- **Remove a touch the branch no longer changes.** A touch also declares what a task will change before it does. Removing it would undo the agent's plan, so narrowing stays a deliberate act.
- **Leave it to the agent to narrow again.** I-0121 shows that this is what costs the time.
