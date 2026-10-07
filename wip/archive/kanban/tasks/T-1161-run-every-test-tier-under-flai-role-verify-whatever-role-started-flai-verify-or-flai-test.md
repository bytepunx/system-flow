---
id: T-1161
type: task
nature: feature
title: "Run every test tier under FLAI_ROLE=verify, whatever role started flai verify or flai test"
status: done
parent: S-0311
owner: alex
created: 2026-10-07T14:29:21Z
updated: 2026-10-07T14:32:31Z
transitions:
  - to: ready
    at: 2026-10-07T14:30:47Z
    by: agent-S-0311
  - to: in-progress
    at: 2026-10-07T14:30:47Z
    by: agent-S-0311
  - to: done
    at: 2026-10-07T14:32:31Z
    by: agent-S-0311
stream: S-0311
tags: [cli]
touches: [flai/internal/verify/run.go, flai/internal/verify/run_test.go, flai/internal/verify/proc_test.go]
usage:
  source: log
  seconds: 104
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 21
      output: 5993
      cache_read: 893600
      cache_write: 37815
      cost: 0.5424
---
# T-1161 Run every test tier under FLAI_ROLE=verify, whatever role started flai verify or flai test

## Work

The tiers `flai verify` and `flai test` run, and the MCP tools `verify` and `test` with them, inherit flai's environment through `OS.Run` in `flai/internal/verify/proc.go`. Under the orchestrator that environment holds `FLAI_ROLE=orchestrate`. A Go test that runs flai in process then hits flai's own orchestrator checks (`flai/cmd/move.go`, `accept.go`, `edit.go`, `order.go`, `plan.go`, and `mcpserver/server.go` read `FLAI_ROLE`), and go-test fails (TH-0260).

- In `runTier` in `flai/internal/verify/run.go`, add `FLAI_ROLE=verify` (`conventions.RoleVerify`) after `opts.Env` in the environment each tier gets. It is the role of a run that only checks, and none of flai's role checks match it. It comes after `opts.Env` and after `os.Environ()`, so it wins over an inherited role and over a caller's `Env`, since exec keeps the last value of a key.
- Leave `FLAI_STORY`, `FLAI_ITEM`, `FLAI_FOCUS`, and `CLOSE_OUT_STORY` as they are.
- In `run_test.go`, check with a fake `Proc` that every tier gets `FLAI_ROLE=verify`, with and without `opts.Env`.
- In `proc_test.go`, run a real process with `t.Setenv("FLAI_ROLE", "orchestrate")` and check that it sees `verify`.

This task waits for nothing.

## Done when

- [ ] Every tier `verify.Verify` and `verify.Test` run gets `FLAI_ROLE=verify`, whatever `FLAI_ROLE` flai itself has.
- [ ] `run_test.go` and `proc_test.go` cover it, and `flai test flai/internal/verify` passes.

## Notes
