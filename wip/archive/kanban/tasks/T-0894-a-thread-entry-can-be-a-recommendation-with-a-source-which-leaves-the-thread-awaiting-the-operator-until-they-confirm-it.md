---
id: T-0894
type: task
nature: feature
title: A thread entry can be a recommendation with a source, which leaves the thread awaiting the operator until they confirm it
status: done
parent: S-0220
owner: alex
created: 2026-10-05T04:46:50Z
updated: 2026-10-06T06:06:48Z
transitions:
  - to: ready
    at: 2026-10-06T05:59:55Z
    by: agent-S-0220
  - to: in-progress
    at: 2026-10-06T05:59:55Z
    by: agent-S-0220
  - to: done
    at: 2026-10-06T06:06:48Z
    by: agent-S-0220
stream: S-0220
tags: [flai]
touches: [flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, design/adrs, design/system/flai-cli.md]
usage:
  source: log
  seconds: 413
  estimated: true
  models:
    - model: claude-haiku-4-5-20251001
      input: 276
      output: 8815
      cache_read: 2186103
      cache_write: 92819
      cost: 0.379
    - model: claude-opus-5-5
      input: 80
      output: 30535
      cache_read: 3781763
      cache_write: 112804
      cost: 2.0178
---
# T-0894 A thread entry can be a recommendation with a source, which leaves the thread awaiting the operator until they confirm it

## Work

Give a thread entry two optional marks in `flai/internal/threads/threads.go`: that it is a recommendation, and the source it cites, a repository path with an optional heading, such as an ADR, a design document, or a convention. Today `Entry` (lines 57-62) has only `At`, `Author`, and `Text`, and `Entries()` (line 175) reads them from the `### <time> <author>` heading. Choose a representation that the wip markdown lint accepts and that an older flai still reads as a plain entry. A suffix on the heading, such as `### <time> <author> (recommendation)`, with the source as the entry's last line, `Source: <path> § <heading>`, is one way.

- `Reply` (line 373) takes the marks, through an options value or a sibling function. A recommendation never sets the status to `answered` and never sets it to `open`; it leaves the status as it was. An answer that cites a source sets `answered` as any reply by another author does.
- Add `Confirm(r, id, author, now)`. It is refused unless the newest entry that is not the opener's is a pending recommendation. It appends an entry by `author` saying the recommendation is confirmed, naming the entry's time and author, and sets `answered`. The recommendation becomes the answer.
- `View` (line 556) and the JSON the MCP tools, `flai thread show --json`, and the dashboard read carry `recommendation`, `source`, and, on the thread, `pending_recommendation` while one waits for the operator.
- Record the representation and the confirm rule as an ADR (`flai adr new`), linked from the thread format in `design/system/flai-cli.md`. Describe the entry marks and `Confirm` there.

This task waits for nothing: every other task of the story builds on these marks.

## Done when

- a test replies with a recommendation and finds the status unchanged, the marks read back, and the lint clean
- a test confirms a pending recommendation, finds the status `answered` and the confirming entry, and is refused confirming a thread with none pending
- a thread file written by an older flai reads as before, with no recommendation
- the ADR is proposed, and `flai-cli.md` describes the marks and `Confirm`
- `go test ./internal/threads/` passes

## Notes
