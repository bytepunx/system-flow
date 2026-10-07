---
id: T-1073
type: task
nature: improvement
title: The host method test.run runs flai test in a story's worktree behind the checks host action and answers its result
status: done
parent: S-0273
owner: alex
created: 2026-10-06T22:52:18Z
updated: 2026-10-07T00:28:39Z
transitions:
  - to: ready
    at: 2026-10-07T00:22:04Z
    by: agent-S-0273
  - to: in-progress
    at: 2026-10-07T00:22:04Z
    by: agent-S-0273
  - to: done
    at: 2026-10-07T00:28:39Z
    by: agent-S-0273
stream: S-0273
tags: [hostapi, testing]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, design/system/dashboard-host-channel.md, flaiover/src/lib/server/agent.ts]
after: [T-1069]
usage:
  source: log
  seconds: 395
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 108
      output: 47676
      cache_read: 6394031
      cache_write: 153894
      cost: 3.1218
---
# T-1073 The host method test.run runs flai test in a story's worktree behind the checks host action and answers its result

## Work

Criterion 2, the host channel half. This task waits for T-1069: a host method builds a `flai` command line, as `checks.run` does, so `flai test --json` has to exist first.

- Add `test.run` to `writes.go` beside `checks.run`.
  - Its arguments: `id`, the story whose worktree it runs in, with the main checkout when none is given; `paths`; `all`; and `max`.
  - It builds `flai test --json` with them and runs it in that worktree.
  - Give it `progress` and a detach timeout as `checks.run` has. Its exits keep a test failure, which is an answer, apart from an error.
- Gate it behind the existing `checks` host action (`ActionChecks`). It runs the project's commands on the host as `checks.run` does, so it needs no new action. This is an assumption on the plan thread.
- Register it where `MethodsFor` in `hostapi.go` needs it, and keep `contract_test.go` passing.
- Describe the method in `design/system/dashboard-host-channel.md`.

## Done when

- [ ] `test.run` is refused while the `checks` action is off, and with it on it answers `flai test --json`'s result for a fixture story's worktree (tests in `writes_test.go`).
- [ ] A test failure comes back as a result, not a channel error.
- [ ] `go test -race ./internal/hostapi/...`, golangci-lint, and the markdown lint pass.

## Notes

Drafted by the planner for S-0273.
