---
id: T-0327
type: task
nature: feature
title: "flai checks run/status/cancel/tail: sequential, in the story's worktree, one at a time, cancellable by killing the whole process group"
status: done
parent: S-0082
owner: alex
created: 2026-09-22T22:36:33Z
updated: 2026-09-22T22:53:49Z
transitions:
  - to: ready
    at: 2026-09-22T22:44:39Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T22:44:39Z
    by: system-flow
  - to: review
    at: 2026-09-22T22:53:49Z
    by: system-flow
  - to: done
    at: 2026-09-22T22:53:49Z
    by: system-flow
stream: S-0082
tags: []
---
# T-0327 flai checks run/status/cancel/tail: sequential, in the story's worktree, one at a time, cancellable by killing the whole process group

## Work
New CLI surface (`flai/cmd/checks.go` or similar) and host-side state (`internal/serve/checks.go`, a state file per story under `serveDir/checks/<id>.json`, the same shape as `agents.json` but keyed by story instead of by project root, since the criterion is one run per *story* at a time, not per project — two different stories' checks may run together).

- `flai checks run <story-id>`: resolves the configured checks (T-0326's resolver); refuses if the story has no worktree (`workitem.Repo.WorktreePath`, checked with `os.Stat`, worded like `syncStoryBranch`'s existing "no worktree at %s; open one with flai stream open %s"); refuses if a run for this story is already active (state file's `running` plus `Alive(pid)`, mirroring `agents.go`'s guard) rather than queuing or replacing it. Writes `running:true` to the state file before starting, runs each named check in order **in the worktree**, stopping at the first non-zero exit; each check's stdout and stderr go to one log file per run (`serveDir/checks/<id>-<ts>.log`); the whole sequence is bounded by the configured time limit (`context.WithTimeout`), and each child is started detached into its own process group (`Detach`, already `Setsid` on Unix) so cancelling can reach anything it spawned, not just the direct child. On completion writes the final state: `outcome` one of `passed`, `failed`, `timed-out`, `cancelled`, which check (if any), its exit code, `started`/`ended`, duration implied by the two. This command **blocks until the run ends** — real minutes, not seconds — the same shape `flai dashboard upgrade` already has.
- `flai checks status <story-id>`: reads and prints the state file; a story with no run yet is `{story, running: false}`, not an error.
- `flai checks cancel <story-id>`: reads the state file's `pid`; if running, sends SIGTERM to the **process group** (new `TerminateGroup`/`KillGroup` in `proc_unix.go`/`proc_windows.go`, extending the existing single-PID `Terminate`), waits a short grace period polling `Alive`, then SIGKILLs the group if it is still there; sets `outcome: "cancelled"`. Refuses plainly if nothing is running.
- `flai checks tail <story-id> --from <byte-offset> [--wait <seconds>, default 20]`: reads the run's log file from the offset, and if the run is still active, polls for new content until `--wait` elapses or the run ends, printing each new complete line as a structured log event as it is found (so that, run through `ExecRunner` with `progress: true`, the existing stderr-JSON-becomes-`$/progress` machinery streams it with no new transport). Exits with `{offset, running, outcome?}` so the caller knows whether to call again.

## Done when
Tests cover: two named checks, the second not reached after the first fails; the worktree-missing refusal; the already-running refusal; a real timeout (a short one in the test) produces `outcome: "timed-out"` and the process is actually gone, not orphaned; `cancel` on a real running process actually stops it and anything it spawned (spawn a child of the check command itself in the test fixture and confirm it is gone too, not just the direct child); `tail` returns exactly the new lines from an offset and the right `running`/`outcome` once the run ends. `go test -race -short ./...` and `golangci-lint run ./...` clean.

## Notes
Depends on T-0326. `flai checks run` is meant to be spawned by the host action (T-0328) with `detachTimeout` set to the configured time limit plus a buffer, the same lesson S-0081 learned the hard way: the connection that asked for a run must not be able to kill it by dying, since nothing about running checks touches that connection at all, but a generic per-request context would still cancel it if not detached.
