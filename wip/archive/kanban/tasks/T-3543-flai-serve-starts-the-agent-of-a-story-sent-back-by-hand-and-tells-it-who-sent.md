---
id: T-3543
type: task
nature: remediation
title: flai serve starts the agent of a story sent back by hand and tells it who sent it back and why
status: done
parent: S-0823
owner: alex
created: 2026-11-17T08:30:58Z
updated: 2026-11-17T10:48:39Z
transitions:
  - to: ready
    at: 2026-11-17T10:27:37Z
    by: agent-S-0823
  - to: in-progress
    at: 2026-11-17T10:27:40Z
    by: agent-S-0823
  - to: done
    at: 2026-11-17T10:48:39Z
    by: agent-S-0823
stream: S-0823
tags: [flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_restart.go, flai/internal/serve/agents_start.go, flai/internal/serve/agents_test.go]
after: [T-3541]
usage:
  source: log
  seconds: 1248
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 64
      output: 15193
      cache_read: 2589054
      cache_write: 128534
      cost: 1.7028
---
# T-3543 flai serve starts the agent of a story sent back by hand and tells it who sent it back and why

## Work

`startAcceptSentBack` in `flai/internal/serve/agents_restart.go` already starts the agent of a story in progress with an open `Acceptance sent back:` thread by `flai` whose last run started before the thread, and `resolveAcceptSentBack` resolves the thread once the story is back in review. Its words, and `harness.SentBack`'s `Acceptance` prompt in `flai/internal/harness/harness.go`, speak of a rebase the acceptance could not settle. For the thread T-3541 opens, tell the agent instead who sent the story back by hand and the reason, as the thread holds them, and to read the thread, do what it asks, and move the story to review again: tell the two threads apart by the title T-3541 writes, and give `harness.SentBack` what it needs to word the prompt (a field of its own, or `What` alone, as fits). Word the journal entry and the run's why in `flai/internal/serve/agents_start.go` the same way, naming who sent it back rather than `flai`.

Test in `flai/internal/serve/agents_test.go` the case I-0397 records: a story whose agent ended, then moved from review to in-progress by the operator with a reason, gets its agent started at the next look, in its last run's session, told of the thread; the thread is resolved once the story is back in review. Test the prompt's words in `flai/internal/harness/harness_test.go`.

This task waits for T-3541, whose thread title and text it reads.

## Done when

- A story sent back by hand with its agent ended has its agent started by flai serve on the thread, told who sent it back and why, and the thread resolves once the story is in review again.
- An acceptance's and a rebase and resolve job's send backs are started with their words as before.
- The tests above pass.

## Notes
