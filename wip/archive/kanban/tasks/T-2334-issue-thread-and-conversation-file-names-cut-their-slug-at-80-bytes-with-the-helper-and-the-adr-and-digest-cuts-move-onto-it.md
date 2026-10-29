---
id: T-2334
type: task
nature: improvement
title: Issue, thread, and conversation file names cut their slug at 80 bytes with the helper, and the ADR and digest cuts move onto it
status: done
parent: S-0542
owner: alex
created: 2026-10-29T15:23:54Z
updated: 2026-10-29T15:31:37Z
transitions:
  - to: ready
    at: 2026-10-29T15:28:48Z
    by: agent-S-0542
  - to: in-progress
    at: 2026-10-29T15:28:49Z
    by: agent-S-0542
  - to: done
    at: 2026-10-29T15:31:37Z
    by: agent-S-0542
stream: S-0542
tags: [flai]
touches: [flai/internal/issues/issues.go, flai/internal/issues/issues_test.go, flai/internal/threads/threads.go, flai/internal/threads/threads_test.go, flai/internal/messages/messages.go, flai/internal/messages/messages_test.go, flai/internal/adr/adr.go, flai/internal/digest/write.go]
after: [T-2333]
usage:
  source: log
  seconds: 168
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 6406
      cache_read: 954283
      cache_write: 52481
      cost: 0.6619
---
# T-2334 Issue, thread, and conversation file names cut their slug at 80 bytes with the helper, and the ADR and digest cuts move onto it

## Work

Name the other files flai makes from a title with the helper T-2333 adds, cut at 80 bytes:

- the issue in `flai/internal/issues/issues.go`, where `flai issue new` builds `<I-nnnn>-<slug>.md`; I-0231's own name is 239 bytes;
- the thread in `flai/internal/threads/threads.go` (`orSlug`), where TH-0465's name is 238 bytes;
- the conversation in `flai/internal/messages/messages.go` (`orSlug`), named as the thread is.

Replace the cuts `flai/internal/adr/adr.go` and `flai/internal/digest/write.go` make by hand with the helper, keeping their names as they are: the digest's `TestSlugIsCutOnAWord` and the ADR tests must still pass unchanged.

Test each of the three new cuts in `issues_test.go`, `threads_test.go`, and `messages_test.go`: a title whose slug is over 80 bytes gives a name of the ID, a dash, at most 80 bytes, and `.md`.

Waits for T-2333, whose helper it calls.

## Done when

- `flai issue new`, `flai thread new`, and `flai message send` with a long title each write a file whose name is at most the ID, a dash, 80 bytes, and `.md`, and the tests that show it pass.
- `adr.go` and `write.go` cut through the helper, and their existing tests pass unchanged.
- `flai test` on the eight files passes.

## Notes
