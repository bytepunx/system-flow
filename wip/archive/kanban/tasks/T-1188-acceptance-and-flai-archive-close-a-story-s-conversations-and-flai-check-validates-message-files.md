---
id: T-1188
type: task
nature: feature
title: Acceptance and flai archive close a story's conversations, and flai check validates message files
status: done
parent: S-0330
owner: alex
created: 2026-10-07T20:14:23Z
updated: 2026-10-07T20:48:28Z
transitions:
  - to: ready
    at: 2026-10-07T20:33:08Z
    by: agent-S-0330
  - to: in-progress
    at: 2026-10-07T20:33:08Z
    by: agent-S-0330
  - to: done
    at: 2026-10-07T20:48:28Z
    by: agent-S-0330
stream: S-0330
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_threads_test.go, flai/cmd/archive.go, flai/cmd/archive_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go, flai/internal/check/scope.go, flai/internal/check/scope_test.go]
after: [T-1186]
usage:
  source: log
  seconds: 920
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 103
      output: 45266
      cache_read: 5987128
      cache_write: 196273
      cost: 3.2841
---
# T-1188 Acceptance and flai archive close a story's conversations, and flai check validates message files

## Work

Close conversations with their stories, and check the files. It waits for T-1186, whose `Close` and parser it calls. It shares no path with T-1187 and runs beside it.

- `flai accept` and `flai archive` close each open conversation the archived story is part of, with an entry saying why, beside where they resolve threads (ADR-0109).
- `flai check` validates each message file: its front matter, that both stories exist, and its markdown; a malformed one is an error.
- No message shows in `flai thread list` or among the threads awaiting anyone: a test proves it.

## Done when

- Tests cover closing at acceptance and at archive, a malformed message file, and a thread list with messages present.
- `flai test` passes on the changed paths.
- `flai check --strict` passes.

## Notes
