---
id: T-1297
type: task
nature: remediation
title: Judge a run whose question was answered before the look as asked, so flai serve starts it again on the answer
status: done
parent: S-0317
owner: alex
created: 2026-10-08T00:01:46Z
updated: 2026-10-08T00:16:17Z
transitions:
  - to: ready
    at: 2026-10-08T00:10:58Z
    by: agent-S-0317
  - to: in-progress
    at: 2026-10-08T00:10:58Z
    by: agent-S-0317
  - to: done
    at: 2026-10-08T00:16:17Z
    by: agent-S-0317
stream: S-0317
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 319
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 46
      output: 12204
      cache_read: 2054628
      cache_write: 96923
      cost: 1.3066
---
# T-1297 Judge a run whose question was answered before the look as asked, so flai serve starts it again on the answer

## Work

In `flai/internal/serve/agents.go`, `judgeRun` records `asked` only when `asking` finds a thread on the story whose last entry is still the agent's. When the operator answers between the agent's end and flai serve's look, the thread no longer awaits the agent, and the run is recorded `failed` (I-0095). S-0335 closed the same gap for conversations with `unanswered`, which counts a message that came after the run started and has had no reply.

Do the same for threads. Add a helper beside `asking`, such as `answeredSince(repo, story, agent, started)`, that returns a thread on the story or one of its tasks (`threads.StoryOf`) on which the agent wrote an entry at or after the run's `Started`, and whose question has an answer after the agent's last entry: an entry by someone else that is not a recommendation awaiting confirmation (ADR-0090), a confirmation, or the thread resolved. Have `judgeRun` take that thread as the run's `Thread` with outcome `asked` and a `Why` that says the question was answered before the end was judged. `resume` then finds `threadAnswered` true at once and starts the agent again in its session, without counting an automatic restart (ADR-0108).

A thread the agent asked on in an earlier run, whose answer it was started again for, has its agent's last entry before this run's `Started`, so it is not counted again. This is the first direction I-0095 names; the second, recording the wait at ask time from `thread_open` and `thread_reply`, is not taken, because it needs new run state written from the MCP server for a gap this closes from what the threads already record.

Add a test in `flai/internal/serve/agents_test.go`, modelled on `TestAReplyBeforeTheEndIsJudgedStartsTheAgentAgain`, that reproduces I-0095: the agent opens a thread on its story after its run started, someone else answers, then `judgeRun` runs with no exit code (as `settleOrphans` calls it). It asserts `OutcomeAsked` and the thread's ID, that `answered` returns the thread, that a look starts the agent again in the same session with `Answered` set, and that `AutoRestarts` is not raised. Cover too: a pending recommendation is no answer, and a run started after the answer is judged `failed` as before.

This task waits for none.

## Done when

- `judgeRun` records `asked`, naming the thread, for a run whose agent asked on its story's thread during the run and was answered before the end was judged.
- The new test fails on the code before the change and passes after it; the existing tests in `flai/internal/serve/agents_test.go` pass.
- `flai test flai/internal/serve` passes.

## Notes
