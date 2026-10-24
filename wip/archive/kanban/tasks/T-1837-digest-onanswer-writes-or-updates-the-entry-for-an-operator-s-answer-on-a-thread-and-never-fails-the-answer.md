---
id: T-1837
type: task
nature: feature
title: digest.OnAnswer writes or updates the entry for an operator's answer on a thread, and never fails the answer
status: done
parent: S-0433
owner: alex
created: 2026-10-16T06:07:29Z
updated: 2026-10-24T20:35:43Z
transitions:
  - to: ready
    at: 2026-10-24T20:30:58Z
    by: agent-S-0433
  - to: in-progress
    at: 2026-10-24T20:30:59Z
    by: agent-S-0433
  - to: done
    at: 2026-10-24T20:35:43Z
    by: agent-S-0433
stream: S-0433
tags: [cli]
touches: [flai/internal/digest/onanswer.go, flai/internal/digest/onanswer_test.go]
usage:
  source: log
  seconds: 284
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10388
      cache_read: 1605528
      cache_write: 88704
      cost: 1.0458
---
# T-1837 digest.OnAnswer writes or updates the entry for an operator's answer on a thread, and never fails the answer

## Work

Add `digest.OnAnswer` in `onanswer.go`, the one function every path that saves a thread entry calls after saving it, with the thread just saved, the runner, and the time:

- It acts only when the thread's newest entry is the operator's answer by S-0429's rule: `digest.Answers` with the manifest's `owner` gives a candidate whose `Ref` is that entry's time. A recommendation, an agent's reply, a resolution, the operator's reply on their own thread, and a permission thread write nothing.
- When no entry cites the thread, it writes one with `digest.New`, with Seed's title and body (`seedTitle`, `seedBody`), kind, topics, and paths, so that `flai digest seed` later finds it recorded by the answer's time and writes no second.
- When an active entry already cites the thread, it edits that one with `digest.Edit`, appending the new answer under its own sentence holding its time, and keeps everything the entry already says, since the curator may have merged other facts into it. An entry whose body already holds the answer's time is left as it is. So an entry the curator merged into another, the thread's source moved there, is updated there and not written again.
- An answer with a secret's shape (`digest.Secrets`) writes nothing, and the warning says so by the rule, `digest.secret`, never by the text.
- Entries are written in the main checkout, where threads are written, whichever checkout the caller runs in.
- It returns what it wrote or updated (the entry's ID and path, and `written` or `updated`), or a warning; it never returns an error that would undo the answer, and the caller logs the warning with `slog` at warn.

First layer: the three entry points call it.

## Done when

`onanswer_test.go` covers a first answer, a second that updates, an entry the curator merged into another, a confirmation, an overruled recommendation, an agent's reply, a recommendation, the operator's own thread, a secret, a run from a worktree, and a write that fails and only warns; `flai test flai/internal/digest` passes.

## Notes

Called from the callers rather than inside `threads.ReplyWith`, so `flai/internal/threads` does not import `digest`, which reads threads.
