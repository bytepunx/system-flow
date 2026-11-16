---
id: T-3399
type: task
nature: remediation
title: scripts/close-out.sh and the template's sync and verify again when the first flai verify stops at sync, within the three rounds
status: done
parent: S-0798
owner: alex
created: 2026-11-16T01:42:36Z
updated: 2026-11-16T02:16:47Z
transitions:
  - to: ready
    at: 2026-11-16T02:08:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T02:08:24Z
    by: system-flow
  - to: done
    at: 2026-11-16T02:16:47Z
    by: agent-S-0798
stream: S-0798
tags: [flai, close-out, template]
touches: [scripts/close-out.sh, template/root/scripts/close-out.sh, flai/internal/closeout/script_test.go]
usage:
  source: log
  seconds: 503
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 31
      output: 232
      cache_read: 3015681
      cache_write: 255791
      cost: 1.2397
---
# T-3399 scripts/close-out.sh and the template's sync and verify again when the first flai verify stops at sync, within the three rounds

## Work

In `scripts/close-out.sh` and `template/root/scripts/close-out.sh`:

- Today the loop's first step, `run "$story" --record-issues` (line 125), exits at once when `flai verify` fails (`[ "$rc" -eq 0 ] || exit "$rc"`, line 131). Only the last check, `--sync-only` (line 151), goes on to `flai stream sync` when its last line is `verify: <story> stopped at sync (`. That is I-0383: a story asked for its close-out while its branch is behind main, as the orchestrator asks for one after `verify_behind_main`, ends in under a second at sync and leaves a failed report as its last.
- When the round's `flai verify` stops at sync (the same last-line match the sync check uses), go on as the sync check does, within the same three rounds. Increment `round` and print `close-out: the branch is behind main; flai stream sync <story> (round <n> of <rounds>)`, the line T-3400 prints too. Then run `flai stream sync "$story"` and loop. At any other step, or after the last round, stop as today.
- `flai stream sync` refuses uncommitted changes. Before the sync, commit what is outstanding with the git commit options given, as the commit step does, and stop with its message when there are changes and no options.
- A sync that stops on a conflict stops the close-out at `flai stream sync` with the hint it has today.
- Update the header comment of each script to say a close-out asked for a branch behind main syncs first. Keep the template's free of this repository's IDs, as its comments are now.

Add `flai/internal/closeout/script_test.go`. No test runs either script today. Copy each script into a temporary root with a stub `scripts/env.sh` and a stub `scripts/flai.sh` that records its arguments. Have the stub's first `verify --record-issues` print `verify: S-0001 stopped at sync (main has 3 commits the branch lacks)` and exit 1, its `stream sync` exit 0, and its later calls pass. Run the script with `sh` on a git repository whose branch is `story/S-0001`. Assert that it ran `stream sync` and a second `verify`, and that its last line is the pass in round 2 of 3. Add a case where `stream sync` fails: the close-out stops at `flai stream sync`. Before the change the first case reproduces I-0383.

This task waits for none and shares no file with T-3400.

## Done when

- Both scripts sync and verify again when the first `flai verify` stops at sync, and the new test passes against each.
- `flai test scripts/close-out.sh template/root/scripts/close-out.sh flai/internal/closeout/script_test.go` passes, ShellCheck included.

## Notes
