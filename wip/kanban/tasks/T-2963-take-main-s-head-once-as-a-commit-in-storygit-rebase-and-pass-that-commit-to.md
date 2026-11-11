---
id: T-2963
type: task
nature: remediation
title: Take main's head once as a commit in storygit.Rebase and pass that commit to the undone check and the changelog settle
status: done
parent: S-0712
owner: alex
created: 2026-11-10T10:27:51Z
updated: 2026-11-11T22:06:46Z
transitions:
  - to: ready
    at: 2026-11-11T21:57:06Z
    by: agent-S-0712
  - to: in-progress
    at: 2026-11-11T21:57:06Z
    by: agent-S-0712
  - to: done
    at: 2026-11-11T22:06:46Z
    by: agent-S-0712
stream: S-0712
tags: [flai, git]
touches: [flai/internal/storygit/sync.go, flai/internal/storygit/undone.go, flai/internal/storygit/undone_test.go]
usage:
  source: log
  seconds: 580
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 66
      output: 443
      cache_read: 2892989
      cache_write: 170663
      cost: 1.1372
---
# T-2963 Take main's head once as a commit in storygit.Rebase and pass that commit to the undone check and the changelog settle

## Work

The cause of I-0338 and I-0335 is a race in `storygit.Rebase` (`flai/internal/storygit/sync.go`). It reads the merge base with the branch name `main`, rebases onto `main` by name, and then `FindUndone` (`flai/internal/storygit/undone.go`) runs `git rev-parse main` again and diffs `oldBase..main` and `main..HEAD`. When main gains a commit after the rebase, the branch lacks it, so `diff main HEAD` lists that commit's lines as removed. `diff oldBase main` lists them as added, so the path is reported and recorded with the old base. `settleChangelogs` passes `res.Base`, the name, to `SettleChangelogs`, which reads main's `CHANGELOG.md` after the rebase the same way (found by S-0709's plan, TH-1014).

- In `Rebase`, after `MainBranch`, resolve the branch to one commit with `git rev-parse <base>`. Use that commit for the `merge-base` read, the `git rebase`, `checkUndone`, and `settleChangelogs`. Keep `res.Base` as the branch's name, so the JSON answer keeps its shape.
- Give `FindUndone` the commit the rebase used. It compares against that commit and never reads main's head again. Update its doc comment and `checkUndone`'s to say so.
- Have `settleChangelogs` pass that commit, not `res.Base`, to `SettleChangelogs`. `SettleChangelogs`' signature stays: the merge queue's `build.settleChangelogs` (`flai/internal/mergequeue/batch.go`) already passes the batch head it rebased onto.
- Add a scripted case to `undone_test.go`. `FindUndone` is given the rebase's commit while the scripted `rev-parse main` would answer a newer one. The case asserts that the check diffs against the given commit and reports nothing for lines only the newer commit added.
- `flai accept` reaches `storygit.Rebase` through `flai/cmd/branch.go`, so it takes the fix with no change of its own.

It waits for no task: it is the fix the others build on.

## Done when

- `storygit.Rebase` reads main's head once and passes that commit to the merge-base read, the rebase, `FindUndone`, and `SettleChangelogs`.
- `FindUndone` runs no `rev-parse` of main's name.
- The new case in `undone_test.go` and the existing undone and sync tests pass under `flai test flai/internal/storygit`.

## Notes
