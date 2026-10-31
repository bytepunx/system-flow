---
id: T-1756
type: task
nature: improvement
title: Every start of a built-in role goes through startOf with its role file, the project's or the baseline, so the text below planner.md's marker reaches the planner's prompt
status: done
parent: S-0419
owner: alex
created: 2026-10-16T05:46:13Z
updated: 2026-10-30T13:03:26Z
transitions:
  - to: ready
    at: 2026-10-30T13:03:20Z
    by: agent-S-0419
  - to: in-progress
    at: 2026-10-30T13:03:21Z
    by: agent-S-0419
  - to: done
    at: 2026-10-30T13:03:26Z
    by: agent-S-0419
stream: S-0419
tags: [cli]
touches: [flai/internal/serve/agents.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go]
after: [T-1753]
usage:
  source: log
  seconds: 5
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 874
      cache_read: 313910
      cache_write: 8167
      cost: 0.1373
---
# T-1756 Every start of a built-in role goes through startOf with its role file, the project's or the baseline, so the text below planner.md's marker reaches the planner's prompt

## Work

- `startOf` in `flai/internal/serve/agents.go` is the start every run goes through: a story's agent's, each built-in role's, and a defined role's.
- In it, `builtinFile` gives the request of the planner, orchestrator, analyzer, or curator its role file (`roles.Builtin`). That is the project's file under the manifest's design folder, else the baseline flai carries.
- `serve.Plan` and the replanner, `serve.Analyze` and the analysis schedule, and the orchestrator's start, stop, and hold all reach the harness through `startOf`. So each one passes its file with no change of its own, and keeps its record under `plans`, `analyzer`, or `orchestrator`.
- Test in `role_test.go`:
  - With no role files, each of the four is given the baseline.
  - The project's `planner.md`, with an addition below the marker, reaches the planner's claude-code definition prompt.
  - A story's agent and a defined role are given none.

After T-1753, which reads the file in the harness.

## Done when

The serve package's tests pass unchanged, with the new test.

## Notes

- Rewritten by agent-S-0419 on TH-0760. The planner wanted to rewire `serve.Plan` into `l.role`. `startOf` is already the start path a defined role takes (`role.go`), so the built-in roles meet it there without moving their records or their when.
