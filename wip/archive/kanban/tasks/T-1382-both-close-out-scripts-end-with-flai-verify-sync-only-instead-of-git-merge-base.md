---
id: T-1382
type: task
nature: improvement
title: Both close-out scripts end with flai verify --sync-only instead of git merge-base
status: done
parent: S-0347
owner: alex
created: 2026-10-08T08:44:29Z
updated: 2026-10-08T10:43:59Z
transitions:
  - to: ready
    at: 2026-10-08T10:43:05Z
    by: agent-S-0347
  - to: in-progress
    at: 2026-10-08T10:43:05Z
    by: agent-S-0347
  - to: done
    at: 2026-10-08T10:43:58Z
    by: agent-S-0347
stream: S-0347
tags: [template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, scripts/README.md, template/root/scripts/README.md, flai/cmd/verify_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
after: [T-1381, T-1379]
usage:
  source: log
  seconds: 53
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 34
      output: 11890
      cache_read: 1943197
      cache_write: 62567
      cost: 1.036
---
# T-1382 Both close-out scripts end with flai verify --sync-only instead of git merge-base

## Work

Both close-out scripts end with `git merge-base --is-ancestor "$base" HEAD`, which fails on any commit main gained during the run (I-0119).

- In `scripts/close-out.sh`, replace "the sync check" with `"$ROOT/scripts/flai.sh" verify "$story" --sync-only`; in `template/root/scripts/close-out.sh`, with `flai verify "$story" --sync-only`, dropping its resolution of the main checkout's branch, which flai does. Keep the step's name and the last line's form, so "stopped at the sync check" still names it, with a hint that says to run `flai stream sync` when it fails.
- Update the header comments of both scripts and the `close-out.sh` rows of `scripts/README.md` and `template/root/scripts/README.md`: the last check passes over commits on main that change only `wip/` paths the branch does not change.
- Run the close-out of a scratch story by hand, or reproduce it in a temporary repository: verify, commit on main a change under `wip/`, then the last check passes and prints the note; a change outside `wip/` stops it.

Waits for T-1381, whose flag it calls, and for T-1379, whose `template/CHANGELOG.md` entry names this change.

## Done when

- Both scripts' last check is `flai verify --sync-only`, and neither calls `git merge-base`.
- The reproduction above passes on a `wip`-only commit and stops on another; say in the task's log how it was run.
- `flai test scripts/close-out.sh template/root/scripts/close-out.sh` passes.

## Notes

Drafted by the planner for S-0347. S-0341 changes the close-out's last line to say when tiers were reused; sync first and keep both.
