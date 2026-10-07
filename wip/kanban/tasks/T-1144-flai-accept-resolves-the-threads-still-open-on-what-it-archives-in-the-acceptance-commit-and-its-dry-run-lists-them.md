---
id: T-1144
type: task
nature: improvement
title: flai accept resolves the threads still open on what it archives, in the acceptance commit, and its dry run lists them
status: done
parent: S-0277
owner: alex
created: 2026-10-07T01:20:20Z
updated: 2026-10-07T02:55:55Z
transitions:
  - to: ready
    at: 2026-10-07T02:10:10Z
    by: agent-S-0277
  - to: in-progress
    at: 2026-10-07T02:10:10Z
    by: agent-S-0277
  - to: done
    at: 2026-10-07T02:55:55Z
    by: agent-S-0277
stream: S-0277
tags: [flai, cli]
touches: [flai/cmd/accept.go, flai/internal/preview/accept.go, flai/cmd/accept_threads_test.go]
after: [T-1142]
usage:
  source: log
  seconds: 625
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 74
      output: 26179
      cache_read: 2662453
      cache_write: 161842
      cost: 2.0543
---
# T-1144 flai accept resolves the threads still open on what it archives, in the acceptance commit, and its dry run lists them

## Work

This task makes the cause of I-0073 stop at its source. An acceptance by the operator, from `flai accept`, `flai move <id> done`, or the board through the host channel's `flai accept`, archives the story with threads on it still `open` or `answered`.

- In `flai/cmd/accept.go`, right after `repo.Archive(ap)` and before the commit, call `threads.ResolveOnItems` with the IDs of every item the plan archived: the story, its tasks, and its epic when the epic followed it. The author is the acceptance's `--by`. The reason is `<S-nnnn> was accepted`. The thread files are then in the acceptance commit, which `git add -A` already stages. Add an `acceptStep` that names the resolved threads when there are any.
- In `flai/internal/preview/accept.go`, add `resolved_threads` to the result, with the IDs. In the dry run, list the threads `threads.OnItems` finds as ones the acceptance would resolve. This is a note, not a blocker. The orchestrator's refusal while a thread is open (ADR-0093, `preview/orchestrator.go`) stays as it is.
- Print the resolved threads in the text output.
- Test this in a new `flai/cmd/accept_threads_test.go`, following the existing `accept_*_test.go` set-up. It reproduces I-0073: a story in review with an answered thread on it and an open thread on one of its tasks. After `flai accept`, both threads are resolved with the reason, they are in the acceptance commit, and `flai check` reports no `threads.archived`. A dry run lists them and changes nothing.

This task waits for T-1142, whose functions it calls. It touches no file of T-1146, so the two run together.

## Done when

- Accepting a story leaves no thread on it, its tasks, or an epic accepted with it `open` or `answered`, and the acceptance commit holds the resolved threads.
- `flai accept --dry-run` and `--json` name the threads it resolves or would resolve.
- The new test passes, and `flai test flai/cmd flai/internal/preview` passes.

## Notes
