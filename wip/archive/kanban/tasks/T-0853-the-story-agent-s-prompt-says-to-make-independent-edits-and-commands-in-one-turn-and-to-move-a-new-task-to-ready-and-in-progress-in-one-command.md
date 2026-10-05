---
id: T-0853
type: task
nature: improvement
title: The story agent's prompt says to make independent edits and commands in one turn and to move a new task to ready and in-progress in one command
status: done
parent: S-0268
owner: alex
created: 2026-10-05T00:29:19Z
updated: 2026-10-05T02:36:57Z
transitions:
  - to: ready
    at: 2026-10-05T02:35:35Z
    by: agent-S-0268
  - to: in-progress
    at: 2026-10-05T02:35:35Z
    by: agent-S-0268
  - to: done
    at: 2026-10-05T02:36:57Z
    by: agent-S-0268
stream: S-0268
tags: [harness, flai]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
usage:
  source: log
  seconds: 82
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 20
      output: 4741
      cache_read: 616546
      cache_write: 27373
      cost: 0.4109
    - model: claude-sonnet-5
      input: 48
      output: 14406
      cache_read: 1356858
      cache_write: 79230
      cost: 0.6136
---
# T-0853 The story agent's prompt says to make independent edits and commands in one turn and to move a new task to ready and in-progress in one command

## Work

Add two sentences to the story agent's prompt in `flai/internal/harness/harness.go`. Put them in `delegation()`, beside the per-task cycle, or in `rules()`, whichever reads better:

- Make independent edits and commands in one turn, as several tool calls in one message. Examples are consecutive edits to one file, reads of files already known, and commands that do not wait on each other.
- Move a task you have just written to ready and in-progress in one command, such as `flai move T-nnnn ready && flai move T-nnnn in-progress`. `flai move` refuses a task from backlog straight to in-progress, and S-0248's agent spent a turn on that refusal.

Assert both sentences in `flai/internal/harness/harness_test.go`, in the per-task cycle test or a new one. Keep the prompts that a resumed run and a commit run get as they are, unless they already carry the block you change.

This task waits for none. It shares no path with T-0854's conventions, so the two can run together.

S-0266's T-0851 also edits `delegation()`. S-0268 is held until S-0266 is accepted, so sync over it first.

## Done when

- The prompt holds both sentences.
- A harness test asserts them. It fails without the sentences and passes with them, shown without `git stash`, as T-0854 says.
- The harness tests and the flai lint pass.

## Notes

Drafted by the planner for S-0268.
