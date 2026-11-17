---
id: T-3505
type: task
nature: remediation
title: The guard finds the foreground sub-agents whose launch has its result in the session transcript and drops them from the session's record
status: done
parent: S-0819
owner: alex
created: 2026-11-17T03:34:46Z
updated: 2026-11-17T07:29:12Z
transitions:
  - to: ready
    at: 2026-11-17T07:17:18Z
    by: system-flow
  - to: in-progress
    at: 2026-11-17T07:17:18Z
    by: system-flow
  - to: done
    at: 2026-11-17T07:29:12Z
    by: agent-S-0819
stream: S-0819
tags: [guard, go]
touches: [flai/internal/guard/transcript.go, flai/internal/guard/transcript_test.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/internal/guard/claudecode.go, flai/internal/guard/claudecode_test.go, flai/internal/guard/call.go]
after: [T-3504]
usage:
  source: log
  seconds: 714
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 3811
      cache_read: 1326633
      cache_write: 29107
      cost: 0.5402
---
# T-3505 The guard finds the foreground sub-agents whose launch has its result in the session transcript and drops them from the session's record

## Work

Build in `flai/internal/guard` the rule T-3504's ADR states; it waits for T-3504 because the signal it reads is the one that task confirms.

- `transcript.go`, new: read a Claude Code session transcript and give the IDs of the sub-agents whose foreground Agent launch has its result there, whether the result is an error or not; a background launch's notice does not count. Read it line by line and pass over a line that does not parse.
- `subagents.go`: a function that, under the record's lock, drops from a session's record the IDs it is given, removing the file when nothing is left, as `RecordChange` does.
- `claudecode.go`: read `transcript_path` from a `PreToolUse` hook input onto the call, so the wait check can find the transcript.

Tests: a transcript with an errored foreground result, a successful one, a background notice, and an unparsable line in `transcript_test.go`; a record pruned under a held lock and pruned to empty in `subagents_test.go`; the path read from the hook input in `claudecode_test.go`.

## Done when

- `flai test` on the six paths passes.
- A record that lists a sub-agent whose foreground launch has its result in the transcript loses it, and one launched in the background keeps it.

## Notes
