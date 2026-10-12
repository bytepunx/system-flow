---
id: T-1476
type: task
nature: remediation
title: "ADR: a rebase of a story branch that removes lines main gained since its old branch point is reported by flai stream sync and refused by flai accept"
status: backlog
parent: S-0371
owner: alex
created: 2026-10-12T17:15:17Z
updated: 2026-10-12T17:15:17Z
transitions: []
stream: S-0371
tags: [cli, docs]
touches: [design/adrs/README.md]
---
# T-1476 ADR: a rebase of a story branch that removes lines main gained since its old branch point is reported by flai stream sync and refused by flai accept

## Work

Record the remedy for I-0124 as an ADR with `flai adr new`, refining ADR-0069 and ADR-0098, and accept it with `flai adr accept` once the plan thread on S-0371 settles it. It waits for nothing: the code tasks wait for it.

The decision to record:

- `storygit.Rebase` reads the branch's old merge base with main before `git rebase`. After a clean rebase, for each path that both main changed since that base and the rebased branch still changes against main, it compares the lines main added since the old base with the lines the branch's net change removes. A path where the branch removes any of them is reported in the result as `undoes_main`: the path, the count of main's lines it removes, and the old base.
- `flai stream sync` and `flai task done` warn on each such path, inside the story's claim or outside it, and say how to restore it: `git checkout <main> -- <path>` in the worktree and a commit. Neither check fails the sync, as ADR-0046's checks do not.
- `flai accept` refuses, merging nothing, a story whose rebase reports such a path, naming each path and the restore, as it refuses conflict markers (S-0253). The story goes back to its agent, or the operator restores the path.
- Git drops a commit that becomes empty on a rebase and replays the revert after it cleanly. That is why no conflict shows and the check runs after the rebase, on the old base read before it.

## Done when

- An ADR under `design/adrs/` records the decision above with its context, I-0124's instance, and its consequences, and is accepted.
- `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

The ADR's file is named by `flai adr new`; its number depends on which story writes an ADR first, so the touch is the index alone and `flai task done` widens it with the file.
