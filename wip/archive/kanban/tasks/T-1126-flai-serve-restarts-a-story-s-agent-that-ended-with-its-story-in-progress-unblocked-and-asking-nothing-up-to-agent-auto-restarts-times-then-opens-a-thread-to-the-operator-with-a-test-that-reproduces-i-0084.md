---
id: T-1126
type: task
nature: remediation
title: flai serve restarts a story's agent that ended with its story in progress, unblocked, and asking nothing, up to agent.auto_restarts times, then opens a thread to the operator, with a test that reproduces I-0084
status: done
parent: S-0294
owner: alex
created: 2026-10-06T23:07:58Z
updated: 2026-10-07T01:59:54Z
transitions:
  - to: ready
    at: 2026-10-07T01:49:51Z
    by: agent-S-0294
  - to: in-progress
    at: 2026-10-07T01:49:51Z
    by: agent-S-0294
  - to: done
    at: 2026-10-07T01:59:54Z
    by: agent-S-0294
stream: S-0294
tags: [cli, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, flai/internal/serve/restart.go, flai/cmd/serve_actions.go, flai/internal/harness/harness.go]
after: [T-1124, T-1125]
usage:
  source: log
  seconds: 603
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 112
      output: 39848
      cache_read: 7324158
      cache_write: 219978
      cost: 3.6601
---
# T-1126 flai serve restarts a story's agent that ended with its story in progress, unblocked, and asking nothing, up to agent.auto_restarts times, then opens a thread to the operator, with a test that reproduces I-0084

## Work

This is criterion 1. Today `judge` in `flai/internal/serve/agents.go` records a run that ended with its story in progress as `OutcomeFailed`, and nothing starts it again. `resume` restarts only an `OutcomeAsked` run whose thread was answered. Make the launcher restart a failed run, as T-1124's ADR decides:

1. **The field.** Add the limit to `AgentConfig`. `serve_actions.go` fills it from T-1125's setting.
2. **The count.** Record the number of automatic restarts on the story's run in `serve/agents.json`, in `AgentRun`. Mark a run started automatically, so that its prompt and the dashboard can say so.
3. **The restart.** When a run ends `OutcomeFailed` with its story in progress, it is started again through the same path as `start(..., after)`, in a new session with the same agent name, and the new agent is told how the last run ended. It happens only when all of these hold:
   - the story is not blocked;
   - the run was not stopped by the operator (`Stopped` unset);
   - the `agent` action is enabled;
   - the count is under the limit.
4. **The reset.** The count goes back to 0 when the operator restarts the agent (`Restart` in `restart.go`) and when the story leaves in-progress.
5. **The thread.** At the limit, open one thread on the story with `threads.New`, authored by flai serve. It says how many times the agent ended and the last run's `why`, and names `flai serve agent restart` and Retry. Open it once per limit reached, not on every look.

Write the tests in `agents_test.go`, on the agent lab the existing launcher tests use, with a fake agent command that exits leaving its story in progress. The first test reproduces I-0084: the run ends, the story is still in progress, and the test asserts that a new run is started. Then test each of these:

- it stops after the limit, and opens exactly one thread;
- 0 means no restart;
- a blocked story is not restarted;
- an operator stop is not restarted;
- an asked run is still restarted only on its answer;
- a manual restart resets the count.

## Done when

- The test that reproduces I-0084 fails on main and passes here.
- The other tests above pass, all with `scripts/flai-test.sh`.
- A story in progress whose agent exits is restarted up to the limit, and then has one thread to the operator.

## Notes

Drafted by the planner. The operator confirmed TH-0208's recommendation on 2026-10-06. S-0272 also changes `flai/internal/serve/restart.go` and tests the restart on an answer (T-1104). Whichever story goes second rebases on the other's change.
