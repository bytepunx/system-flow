---
id: T-3105
type: task
nature: remediation
title: flai serve agent harness clears the budget state of the program, arguments, or login it changes
status: done
parent: S-0720
owner: alex
created: 2026-11-11T08:55:12Z
updated: 2026-11-11T22:55:34Z
transitions:
  - to: ready
    at: 2026-11-11T22:36:30Z
    by: agent-S-0720
  - to: in-progress
    at: 2026-11-11T22:36:30Z
    by: agent-S-0720
  - to: done
    at: 2026-11-11T22:55:34Z
    by: agent-S-0720
stream: S-0720
tags: [flai, cli, serve, budget]
touches: [flai/cmd/serve_harness.go, flai/cmd/serve_agent.go, flai/cmd/serve_actions_test.go]
after: [T-3098]
usage:
  source: log
  seconds: 1144
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 75
      output: 21453
      cache_read: 5876546
      cache_write: 99360
      cost: 2.2073
---
# T-3105 flai serve agent harness clears the budget state of the program, arguments, or login it changes

## Work

- In `flai/cmd/serve_harness.go`'s `RunE`, compare the harness entry before and after the change, for the host's harness and for a project's. When its program, its arguments, or its login changed, call the clear T-3098 adds for the old identity's state, with why naming the harness, the project or the host, and what changed. A change of anything else, or none, clears nothing.
- The dashboard's harness setting (`settings.harness` in `flai/internal/hostapi/settings.go`) runs this command, so it needs no change of its own; check that it still does.
- Say the clear in the command's output where `showAgentCommand` (`flai/cmd/serve_agent.go`) prints the result, and in its `--json`.
- Cover in `flai/cmd/serve_actions_test.go`: a program change lifting a held state, so the next start goes as the probe; a login change doing the same; a change that touches neither leaving it held.
- Waits for T-3098: it calls the clear T-3098 adds. It runs beside T-3104, with no path in common.

## Done when

- A harness change of program or login lifts the hold for the old identity, journaled under `serve.budget` with why, in a test.
- A harness change of neither leaves the hold, in a test.
- `flai test` on the changed paths passes.

## Notes
