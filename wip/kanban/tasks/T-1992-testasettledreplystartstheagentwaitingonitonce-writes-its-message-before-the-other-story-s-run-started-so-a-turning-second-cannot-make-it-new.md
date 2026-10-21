---
id: T-1992
type: task
nature: remediation
title: TestASettledReplyStartsTheAgentWaitingOnItOnce writes its message before the other story's run started, so a turning second cannot make it new
status: done
parent: S-0473
owner: alex
created: 2026-10-21T12:08:03Z
updated: 2026-10-21T12:10:53Z
transitions:
  - to: ready
    at: 2026-10-21T12:10:39Z
    by: agent-S-0473
  - to: in-progress
    at: 2026-10-21T12:10:40Z
    by: agent-S-0473
  - to: done
    at: 2026-10-21T12:10:53Z
    by: agent-S-0473
stream: S-0473
tags: []
touches: [flai/internal/serve/agents_test.go, design/issues/I-0158-testasettledreplystartstheagentwaitingonitonce-in-flai-internal-serve-fails-now-and-then-in-the-close-out-s-integration-tier-on-a-branch-that-does-not-touch-serve.md, design/issues/summary.md]
usage:
  source: log
  seconds: 11
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 2
      output: 16
      cache_read: 139725
      cache_write: 1683
      cost: 0.0479
---

# T-1992 TestASettledReplyStartsTheAgentWaitingOnItOnce writes its message before the other story's run started, so a turning second cannot make it new

## Work

The test gives the other story's run and the message their times from separate `time.Now()` calls: `putEnded("Other")` stamps the other run's start and end, then `putEnded("Messages and waits")` stamps the run whose start the message is sent at. `messages.ToSince` takes what is strictly after the other run's start, so when the clock's second turns between the two calls, the message is new to the other story's agent. `look` then starts that agent for it, and the test fails "started before the reply". flai serve is right to start it: a message after a run ended is new to it (ADR-0154). The test's timing is what is wrong.

Reproduce it every time by putting the other run a second earlier, as when the second turns. Then fix it by putting the waiting story's run, and so the message, a minute before the other run started, with `endedAgo`. The message is then never new to the other story's agent, however the clock's seconds fall. Close I-0158 saying so. Waits for no task: it is the story's only one.

## Done when

- The test puts the other run a second earlier, and fails without the fix.
- With the fix it passes, through `flai test` on the file.
- I-0158 is closed with the reason.

## Notes
