---
id: T-1188
type: task
nature: feature
title: Acceptance and flai archive close a story's conversations, and flai check validates message files
status: backlog
parent: S-0330
owner: alex
created: 2026-10-07T20:14:23Z
updated: 2026-10-07T20:14:23Z
transitions: []
stream: S-0330
tags: [flai]
touches: [flai/cmd/accept.go, flai/cmd/accept_threads_test.go, flai/cmd/archive.go, flai/cmd/archive_test.go, flai/internal/check/check.go, flai/internal/check/check_test.go]
after: [T-1186]
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
