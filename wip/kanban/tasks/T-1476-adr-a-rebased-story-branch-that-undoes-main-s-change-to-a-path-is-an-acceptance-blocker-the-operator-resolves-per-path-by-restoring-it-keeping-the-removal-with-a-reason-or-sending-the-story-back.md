---
id: T-1476
type: task
nature: remediation
title: "ADR: a rebased story branch that undoes main's change to a path is an acceptance blocker the operator resolves per path, by restoring it, keeping the removal with a reason, or sending the story back"
status: done
parent: S-0371
owner: alex
created: 2026-10-12T17:15:17Z
updated: 2026-10-14T21:35:56Z
transitions:
  - to: ready
    at: 2026-10-14T21:34:47Z
    by: agent-S-0371
  - to: in-progress
    at: 2026-10-14T21:34:48Z
    by: agent-S-0371
  - to: done
    at: 2026-10-14T21:35:56Z
    by: agent-S-0371
stream: S-0371
tags: [cli, docs]
touches: [design/adrs/README.md]
usage:
  source: log
  seconds: 68
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 16
      output: 136
      cache_read: 1345501
      cache_write: 11284
      cost: 0.6255
---
# T-1476 ADR: a rebased story branch that undoes main's change to a path is an acceptance blocker the operator resolves per path, by restoring it, keeping the removal with a reason, or sending the story back

## Work

Record the remedy for I-0124 as an ADR with `flai adr new`, refining ADR-0069, ADR-0093, and ADR-0098, and accept it with `flai adr accept`. The operator confirmed its shape on TH-0420 (2026-10-14T18:14:21Z): resolve it from the dashboard, with flai doing the git work, not with git commands or hand-edited files. It waits for nothing: the code tasks wait for it.

The decision to record:

- **The check.** `storygit.Rebase` reads the branch's old merge base with main before `git rebase`. After a clean rebase, for each path that main changed since that base and that the rebased branch still changes against main, it compares the lines main added since the old base with the lines the branch's net change removes. A path where the branch removes any of them is `undoes_main`: the path, the count of main's lines it removes, and the old base. Git drops a commit that becomes empty on a rebase and replays the revert after it cleanly, so no conflict shows; the check therefore runs after the rebase, against the base read before it.
- **The record.** A later rebase starts from a base that already holds main's change, so it cannot find the path again. Each finding is therefore recorded per story when the sync makes it, under `.flai-cache/`, and the acceptance preview reads the record. A recorded path whose branch no longer removes main's lines, against main now, drops out by itself. Name the file and its shape in the ADR.
- **At sync.** `flai stream sync` and `flai task done` warn on each such path, inside the story's claim or outside it, and name the restore: `flai stream restore S-nnnn --path <path>`. Neither fails, as ADR-0046's checks do not.
- **The blocker.** The acceptance preview (`flai accept --dry-run`, the host read `accept.preview`) lists each recorded path that is not kept as a blocker, and gives it as data, `undoes_main`, with the lines main gained that the branch removes, so the dashboard can show each path as a card. `flai accept` refuses on it before anything is merged. Its own rebase runs the check too, so a path that only the acceptance's rebase finds is refused, merging nothing.
- **Restore.** `flai stream restore S-nnnn --path <path>` commits main's version of the path onto `story/S-nnnn` in its worktree, refusing while the worktree has uncommitted changes, and clears the path. The dashboard runs it through a new host write, `stream.restore`.
- **Keep.** `flai accept S-nnnn --keep-removal <path> --reason "<why>"` passes that path. The reason is recorded under `### Removals kept` in the story's Notes, in the acceptance commit. `accept.run` takes `keep_removals`, each a path and a reason. flai refuses `--keep-removal` with `--by orchestrator`: an undone path waits for the operator.
- **Send back.** Moving the story back to in-progress with a reason, which exists, and the story's agent restores or redoes the path.
- **The dashboard.** The review page's cards and the board's badge are a story of their own, after this one; this ADR names the data they read.

## Done when

- An ADR under `design/adrs/` records the decision above with its context, I-0124's instance, TH-0420's choice, its consequences, and the alternatives considered (the `git checkout` advice the first plan gave, and the orchestrator choosing Keep), and is accepted.
- `design/adrs/README.md` lists it.
- `flai check --strict` and the markdown lint pass on both files.

## Notes

The ADR's file is named by `flai adr new`; its number depends on which story writes an ADR first, so the touch is the index alone and `flai task done` widens it with the file.
