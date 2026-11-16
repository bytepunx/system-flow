---
id: T-3400
type: task
nature: remediation
title: A container story's close-out syncs and verifies again when its first flai verify stops at sync, as the script does
status: done
parent: S-0798
owner: alex
created: 2026-11-16T01:42:46Z
updated: 2026-11-16T02:17:02Z
transitions:
  - to: ready
    at: 2026-11-16T02:08:24Z
    by: system-flow
  - to: in-progress
    at: 2026-11-16T02:08:24Z
    by: system-flow
  - to: done
    at: 2026-11-16T02:17:02Z
    by: agent-S-0798
stream: S-0798
tags: [flai, close-out]
touches: [flai/internal/serve/verifyjob_container.go, flai/internal/serve/verifyjob_test.go]
usage:
  source: log
  seconds: 518
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 22
      output: 133
      cache_read: 2134127
      cache_write: 39208
      cost: 0.8235
---
# T-3400 A container story's close-out syncs and verifies again when its first flai verify stops at sync, as the script does

## Work

A story whose agent runs in a container is closed out by `containerCloseOut.run` in `flai/internal/serve/verifyjob_container.go` (line 97), not by the script (ADR-0274). It has the same gap. The container's `flai verify` returns at once when it fails (line 144, `return end(code, err)`). Only the host's `--sync-only` check (line 165) goes on to `flai stream sync` and the next round.

- Keep the container verify's output, as the sync check keeps it in `said`. When it fails and its last line starts `verify: <story> stopped at sync (`, and `round < closeOutRounds`, go on as the sync check does. Increment `round` and print `close-out: the branch is behind main; flai stream sync <story> (round <n> of <rounds>)`, the line T-3399 prints too. Run the commit step on the host first, since `flai stream sync` refuses uncommitted changes. Then run `gitPointers` and `flai stream sync` on the host, and loop. At any other step, or after the last round, stop as today.
- Update the function's comment and the comment at line 24 to say so.
- Extend `TestAContainerStorysCloseOutCommitsAndSyncsOnTheHost` in `flai/internal/serve/verifyjob_test.go` (line 656), or add a test beside it. Have the lab's first container verify print `verify: <story> stopped at sync (main has 1 commit the branch lacks)` and exit 1. Assert that the host ran `flai stream sync` and a second container verify (`<id>-r2`), and that the job's last line is the pass in round 2 of 3. Before the change it reproduces I-0383 for a container story.

This task waits for none and shares no file with T-3399.

## Done when

- A container story's close-out whose first verify stops at sync syncs on the host and verifies again, and the new test passes.
- `flai test flai/internal/serve/verifyjob_container.go flai/internal/serve/verifyjob_test.go` passes.

## Notes
