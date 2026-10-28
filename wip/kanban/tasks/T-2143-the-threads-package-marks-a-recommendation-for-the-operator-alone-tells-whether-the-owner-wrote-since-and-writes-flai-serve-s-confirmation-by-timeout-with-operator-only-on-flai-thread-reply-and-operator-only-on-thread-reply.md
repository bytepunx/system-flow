---
id: T-2143
type: task
nature: improvement
title: The threads package marks a recommendation for the operator alone, tells whether the owner wrote since, and writes flai serve's confirmation by timeout, with --operator-only on flai thread reply and operator_only on thread_reply
status: done
parent: S-0502
owner: alex
created: 2026-10-26T04:43:57Z
updated: 2026-10-28T13:05:09Z
transitions:
  - to: ready
    at: 2026-10-28T12:47:49Z
    by: agent-S-0502
  - to: in-progress
    at: 2026-10-28T12:47:50Z
    by: agent-S-0502
  - to: done
    at: 2026-10-28T13:05:09Z
    by: agent-S-0502
stream: S-0502
tags: [flai]
touches: [flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/server_test.go, flai/internal/mcpserver/folder.go, docs/users/flai-reference.md]
after: [T-2141]
usage:
  source: log
  seconds: 1039
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 14675
      cache_read: 3787886
      cache_write: 124501
      cost: 1.8784
---
# T-2143 The threads package marks a recommendation for the operator alone, tells whether the owner wrote since, and writes flai serve's confirmation by timeout, with --operator-only on flai thread reply and operator_only on thread_reply

## Work

The thread side of T-2141's ADR, points 3 and 4, in `flai/internal/threads/threads.go` with tests in `threads_test.go`:

- **The mark.** `ReplyOptions` gains `OperatorOnly`, allowed only with `Recommendation`; `ReplyWith` writes the mark in the form the ADR fixes (recommended: the one-line paragraph `Awaits: the operator alone` just above the `Source:` line), and `Entries()` reads it into `Entry.OperatorOnly`, keeping it in the text as the source line is kept. `View()` carries `operator_only` on each entry and on `pending_recommendation`, absent when false, so `flai thread show --json`, the host API, and the MCP tools show it.
- **Who wrote since.** `(*Thread).WroteSince(author string, at time.Time) bool`: whether an entry by that author, the opener included, comes after `at`. `PendingRecommendation` is unchanged, since it skips the opener's entries by ADR-0090 and a test in `threads_test.go` asserts that.
- **The confirmation by timeout.** `ConfirmByTimeout(r, id, author, now, waited time.Duration)`: as `Confirm`, with the line `Confirmed the recommendation of <time> <author> after <duration> with no answer from the operator (orchestration.permissions.confirm_after).`, the recommendation's source cited, status `answered`, refused when no recommendation is pending, when it is marked for the operator alone, and to its own author. `Confirms()` keeps matching it by its prefix. A test shows the digest's `confirmed` pattern in `flai/internal/digest/answer.go` still matches the line, and that `Answers` yields nothing for it since its author is not the owner.
- `flai/cmd/thread.go`: `--operator-only` on `flai thread reply`, refused without `--recommend`; `thread_test.go` covers both. `flai guard` reads `--recommend` from the arguments and lets the other flags through, so it needs no change.
- `flai/internal/mcpserver/server.go`: `operator_only` on `thread_reply`, a boolean described as the ADR words it, passed to `ReplyWith`; `server_test.go` covers the write and the field on `thread_get`'s view.

## Done when

- A recommendation posted with `operator_only` reads back with the mark in the file and `operator_only: true` in `flai thread show --json`, and a test shows the entry's author and `recommendation` read as before without the new parsing.
- `ConfirmByTimeout` writes the entry the ADR fixes, sets `answered`, and is refused for a marked recommendation and when none is pending.
- `WroteSince` is true for an owner entry after the recommendation, whether or not the owner opened the thread.
- `flai test` on the changed paths passes.

## Notes

Layer 2, waits for T-2141 (the mark's and the line's form). The orchestrator's prompt task, the serve look, and the metrics task wait for it. It shares no path with T-2142.
