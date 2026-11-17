---
id: T-3429
type: task
nature: feature
title: "flai serve runs a resolve job: it rebases the story's branch, runs its close-out, and sends the story back only on conflicts or a failed step"
status: done
parent: S-0804
owner: alex
created: 2026-11-16T23:19:25Z
updated: 2026-11-17T00:25:13Z
transitions:
  - to: ready
    at: 2026-11-16T23:20:38Z
    by: agent-S-0804
  - to: in-progress
    at: 2026-11-16T23:57:25Z
    by: agent-S-0804
  - to: done
    at: 2026-11-17T00:25:13Z
    by: agent-S-0804
stream: S-0804
tags: []
touches: [flai/internal/sendback/sendback.go, flai/internal/sendback/sendback_test.go, flai/cmd/accept.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/resolvejob.go, flai/internal/serve/resolvejob_test.go]
after: [T-3428]
usage:
  source: log
  seconds: 1668
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 128
      output: 53920
      cache_read: 9727856
      cache_write: 365533
      cost: 5.0639
---
# T-3429 flai serve runs a resolve job: it rebases the story's branch, runs its close-out, and sends the story back only on conflicts or a failed step

## Work

Waits for T-3428, whose `closeout` contract (`KindResolve`, `Job.Outcome`, `Thread`, `Conflicts`, `Reason`, the `resolve:` progress lines) this builds on.

1. Move the send back's shared parts out of `flai/cmd` into a new package `flai/internal/sendback`, so `flai serve` can send a story back without importing `cmd`: `sendBackTarget` (the state the workflow's send-back transition takes a story to) and `sendBackThread` (open the `Acceptance sent back:` thread by `flai`, or write on the one still open). `flai/cmd/accept.go` calls them; its behaviour is unchanged. `storygit.SentBackAuthor` and `SentBackPrefix` stay where they are, so `agents_restart.go` starts the agent on the thread as today.
2. In `flai/internal/serve`, a request whose `Kind` is `resolve` (put the job's own code in `resolvejob.go`):
   - is refused, running nothing, when the story is not in review;
   - while any job of the story runs, is answered with that job's id and starts nothing, as a close-out request is;
   - is never answered from the verify record: it always rebases first;
   - otherwise starts a job of kind `resolve`, recorded in `S-nnnn.job.json` with its output in `S-nnnn.log` as a close-out job's are, whose body: writes the `resolve:` step lines; runs this flai's `stream sync S-nnnn` in the worktree (the rebase, as `flai stream sync` does); when that stops on conflicts, reads the conflicting paths, aborts the rebase (`storygit.AbortRebase`), and ends; else runs the story's close-out exactly as a close-out job would (the worktree's script on the host, or the container close-out of ADR-0274), with a commit message `chore: [S-nnnn] close out again after rebase and resolve`.
   - Every close-out and resolve job's environment sets `FLAI_VERIFY_PROGRESS=1`.
3. At the job's end:
   - a clean rebase and a passing close-out: `Outcome` stayed in review, nothing sent back, no agent started;
   - conflicts: the story is sent back to in-progress as `flai` with a reason naming the paths, and the send-back thread names them, the branch, the base, and that flai aborted the rebase;
   - a close-out that stopped or failed: sent back the same way, the thread naming the step it stopped at, that step's findings from the story's last verify report (at most ten, the rest counted), and the job's log path;
   - before sending back, read the story again: one no longer in review is left as it is, with `Reason` saying so;
   - write `Outcome`, `Thread`, `Conflicts`, `Reason` in the job, journal it, and announce the verified event as a close-out job does.
   - A job flai serve's restart interrupts is swept interrupted, as today, and sends nothing back.
4. Tests with a fake runner, as `verifyjob_test.go` does: each of the three ends, the second request, the refusal of a story not in review, the interrupted sweep sending nothing, and the moved package.

## Done when

- `flai test` passes on the changed paths.
- Each end above is covered by a test that fails without the change.

## Notes
