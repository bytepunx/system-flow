---
id: T-1837
type: task
nature: feature
title: digest.OnAnswer writes or updates the entry for an operator's answer on a thread, and never fails the answer
status: backlog
parent: S-0433
owner: alex
created: 2026-10-16T06:07:29Z
updated: 2026-10-16T06:07:29Z
transitions: []
stream: S-0433
tags: [cli]
touches: [flai/internal/digest/onanswer.go, flai/internal/digest/onanswer_test.go]
---
# T-1837 digest.OnAnswer writes or updates the entry for an operator's answer on a thread, and never fails the answer

## Work

Add `digest.OnAnswer(root, thread, entry)` in `onanswer.go`, the one function every path that saves a thread entry calls after saving it:

- When the entry is the operator's answer by S-0429's rule (`digest.Answers`) on a thread the operator did not open, write the entry with `digest.New`, or, when an entry already cites the thread, update it with `digest.Edit`; an entry whose source the curator moved into another is left alone.
- A recommendation, an agent's reply, and the operator's reply on their own thread write nothing.
- An answer with a secret's shape writes nothing, and the result says so.
- It returns what it wrote or updated, or a warning; it never returns an error that would undo the answer.

First layer: the three entry points call it.

## Done when

`onanswer_test.go` covers a first answer, a second that updates, a confirmation, an overruled recommendation, an agent's reply, the operator's own thread, a secret, and a write that fails and only warns; `flai test flai/internal/digest` passes.

## Notes

Called from the callers rather than inside `threads.ReplyWith`, so `flai/internal/threads` does not import `digest`, which reads threads.
