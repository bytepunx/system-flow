---
id: T-2129
type: task
nature: improvement
title: flai serve blocks a story whose agent ended asking, with the thread as the reason, and unblocks it when the answer starts the agent again
status: done
parent: S-0498
owner: alex
created: 2026-10-26T04:14:02Z
updated: 2026-10-26T04:33:29Z
transitions:
  - to: ready
    at: 2026-10-26T04:19:52Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:19:53Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T04:33:29Z
    by: agent-S-0498
stream: S-0498
tags: [flai]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go]
after: [T-2127]
usage:
  source: log
  seconds: 816
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 97
      output: 704
      cache_read: 4649843
      cache_write: 137525
      cost: 1.6198
---
# T-2129 flai serve blocks a story whose agent ended asking, with the thread as the reason, and unblocks it when the answer starts the agent again

## Work

In `flai/internal/serve/agents.go`:

- Where `judgeRun` settles a run as `OutcomeAsked` with its `Thread` (the agent ended its session on an open question to the designer, ADR-0154), append a blocked interval to the story with `workitem.BlockItem` and the reason `awaiting the operator on TH-nnnn`, unless `IsBlocked()` already holds; write the item as the launcher writes its other item changes, and log the block to the serve log with the story and the thread.
- Where `resume` starts the agent again because `answered()` found the thread answered, confirmed, or resolved, close the interval with `workitem.UnblockItem` before the start, so the story is not blocked while its agent works. A story blocked by someone else's reason (one `flai block` wrote) is left as it is: only an interval whose reason names the thread is closed.
- Keep ADR-0108's rule: `restartEnded` still skips a blocked story, and `startSentBack` too; the asked run is restarted by `resume` alone, as today (`resume` already restarts a blocked story for `run.Thread`).

Tests in `agents_test.go`: a run judged asked blocks the story with the reason naming the thread; a story already blocked is not blocked twice; the answer on the thread closes the interval and starts the agent; a story blocked by another reason keeps that interval; and a run that ended any other way blocks nothing. This is the test that reproduces I-0151's "story idles in progress": before the change, the asked story carries no blocked interval.

This task waits for T-2127, the ADR that decides the block. It changes no other path, so it runs with T-2128.

## Done when

- A story whose agent ended asking shows blocked on `flai board` with `awaiting the operator on TH-nnnn`, and is unblocked at the restart on the answer.
- The five tests above pass under `flai test flai/internal/serve/agents.go`, and the existing restart tests are unchanged.

## Notes

Layer 2, with T-2128. The reason's wording is the ADR's: change both if the operator asks for another.
