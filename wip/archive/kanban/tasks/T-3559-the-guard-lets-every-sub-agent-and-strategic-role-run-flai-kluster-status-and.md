---
id: T-3559
type: task
nature: remediation
title: The guard lets every sub-agent and strategic role run flai kluster status and deployments
status: done
parent: S-0829
owner: alex
created: 2026-11-17T12:28:43Z
updated: 2026-11-17T13:11:25Z
transitions:
  - to: ready
    at: 2026-11-17T13:04:58Z
    by: agent-S-0829
  - to: in-progress
    at: 2026-11-17T13:04:58Z
    by: agent-S-0829
  - to: done
    at: 2026-11-17T13:11:25Z
    by: agent-S-0829
stream: S-0829
tags: [flai, guard]
touches: [flai/internal/guard/guard_tables.go, flai/internal/guard/kluster_test.go]
usage:
  source: log
  seconds: 387
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 757
      cache_read: 458061
      cache_write: 16851
      cost: 0.2416
---
# T-3559 The guard lets every sub-agent and strategic role run flai kluster status and deployments

## Work

Reproduce I-0355 first: in `flai/internal/guard/kluster_test.go`, a test that runs `flai kluster status`, `scripts/flai.sh kluster status --json`, and `flai kluster deployments --json` through `Decide` as a story's sub-agent, a planner's sub-agent, the operator's sub-agent, and the planner, orchestrator, analyzer, and curator sessions, and asserts each is allowed. Today the sub-agent cases fall through `asSubAgent` to the generic refusal, since `cliReads` has no `kluster` entry; `TestFlaiKlusterWritesAreTheOperators` checks only that a read is not refused as a write, so it passed over the gap.

Then add `"kluster": {"status", "deployments"}` to `cliReads` in `flai/internal/guard/guard_tables.go`, with a comment naming `design/system/kluster.md` § Reads and S-0829. `subcommandsOf` in `roles.go` reads `cliReads`, so the roles a file holds take the same subcommands. Check in the same test that `flai kluster up`, `deploy`, `exec`, `log`, and `env` keep the refusals they have now. This task waits for no other: it is the fix the others describe.

## Done when

- The new test fails on main's `guard_tables.go` and passes with the fix.
- `flai test flai/internal/guard/guard_tables.go flai/internal/guard/kluster_test.go` passes.

## Notes
