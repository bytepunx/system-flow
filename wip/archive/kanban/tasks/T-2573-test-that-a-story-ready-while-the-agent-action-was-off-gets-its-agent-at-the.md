---
id: T-2573
type: task
nature: remediation
title: Test that a story ready while the agent action was off gets its agent at the first look once the action is on and its agent names a harness
status: done
parent: S-0589
owner: alex
created: 2026-11-07T10:05:50Z
updated: 2026-11-07T10:17:25Z
transitions:
  - to: ready
    at: 2026-11-07T10:11:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T10:11:54Z
    by: system-flow
  - to: done
    at: 2026-11-07T10:17:25Z
    by: agent-S-0589
stream: S-0589
tags: [flai, serve, test]
touches: [flai/internal/serve/agents_test.go]
usage:
  source: log
  seconds: 331
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 4985
      cache_read: 3977336
      cache_write: 56483
      cost: 1.3408
  strategic:
    - kind: orchestrator
      seconds: 3
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 1
          output: 8
          cache_read: 28935
          cache_write: 760
          cost: 0.0074
---
# T-2573 Test that a story ready while the agent action was off gets its agent at the first look once the action is on and its agent names a harness

## Work

I-0300's guess, that flai serve starts an agent only when a story enters ready, does not hold. `(*launcher).look` in `flai/internal/serve/agents.go` reads the `agent` action at every look, on every change under `kanban/` and `threads/` and every minute, and starts each ready story that has had no run since it entered ready (ADR-0041). `~/.flai/serve/serve.log` shows what happened to S-0001 in flai-scratch-import: at 08:25:49Z `agent not started ... the agent host action is off`; the operator enabled `agent` at 08:26:06Z; the next minute look, at 08:26:49Z, logged `S-0001 names no harness, and no command is set on the host`. The imported project's default agent named no harness, so nothing could start S-0001. The operator set the default agent at 08:42:02Z and S-0001 started at 08:42:14Z. That missing default is S-0587's.

In `flai/internal/serve/agents_test.go`, beside `TestNothingIsStartedWhenItShouldNotBe`, add one test that replays the instance with `newAgentLab`:

1. A story enters ready while `lab.cfg.Enabled` is false; a look starts nothing and the state's `waiting` says the action is off.
2. The action is turned on, with nothing to start the story with (no harness on its agent or the project's default, and no command); a look starts nothing and `waiting` says it names no harness.
3. The story's agent gains what starts it, without the story re-entering ready (the lab's command, or the project's default agent in `system-flow.yaml` if the lab can set it); the next look starts exactly one run for it, and the journal names it once.

Also cover the plain case: ready while off, then on with something to start it, then one look starts it. Both may be subtests of one test. They pass on the code as it is: say so in the test's comment, naming I-0300 and that its cause was the missing harness.

## Done when

- The test is in `agents_test.go` and passes under `flai test flai/internal/serve/agents_test.go`.
- It would fail if the launcher skipped a story that was in ready before the action was turned on, as read from the code path it covers.

## Notes
