---
id: T-2574
type: task
nature: remediation
title: Start a story already in ready on the first look after the agent action is turned on for its project
status: cancelled
parent: S-0589
owner: alex
created: 2026-11-07T10:06:13Z
updated: 2026-11-07T10:10:57Z
transitions:
  - to: cancelled
    at: 2026-11-07T10:10:57Z
    by: agent-S-0589
stream: S-0589
tags: [flai, serve]
touches: [flai/internal/serve/agents.go, flai/internal/serve/serve.go, flai/internal/serve/agents_test.go, flai/internal/serve/serve_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go]
after: [T-2573]
usage:
  source: sum
  seconds: 0
  models: []
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
# T-2574 Start a story already in ready on the first look after the agent action is turned on for its project

## Work

Fix the cause T-2573 recorded under `### Cause` in the story's Notes; it waits for T-2573 because the fix depends on which candidate that task's tests confirmed. ADR-0041 and ADR-0043 already decide the behaviour: a story in ready gets its agent whenever the in-progress limit has room, whenever it entered ready. This task makes the code meet them when the `agent` action comes on after the story is ready. Fix the confirmed cause, choosing among these:

- Turning the `agent` action on, from a shell (`flai serve enable agent`, `flai/cmd/serve_actions.go`) or from the dashboard, makes the project's launcher look again at once (`starter.again` in `flai/internal/serve/serve.go`), rather than on the next change under `kanban/` or `threads/`.
- The project loop's minute look runs every minute however often the orchestrator's and the roles' wakes fire: a ticker made once, not a `time.After` made again on every pass.
- A project the board's Import adds while serve runs gets a launcher with its configuration, and its looks run.

Correct the `ActionAgent` comment and its description in `flai/internal/hostapi/writes.go`, which say the action starts an agent only when a story becomes ready; change `writes_test.go` or `serve_actions_test.go` only where they pin that text. Make T-2573's tests pass without weakening them.

## Done when

- T-2573's failing test passes, and the tests beside it in `agents_test.go` and `serve_test.go` still pass.
- `flai test` passes on every file this task changed.
- The `ActionAgent` description in `writes.go` says a story already in ready gets its agent once the action is on.

## Notes
- 2026-11-07T10:10:57Z: moved to cancelled: No fix needed in the launcher: it already starts a story that was ready before the agent action came on (ADR-0041), and serve.log shows it looked at 08:26:49Z and found S-0001 named no harness. The missing default harness is S-0587's.
