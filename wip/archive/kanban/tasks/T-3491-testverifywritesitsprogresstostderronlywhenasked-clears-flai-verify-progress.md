---
id: T-3491
type: task
nature: remediation
title: TestVerifyWritesItsProgressToStderrOnlyWhenAsked clears FLAI_VERIFY_PROGRESS before its run without it, so it passes in a close-out job (I-0395)
status: done
parent: S-0797
owner: alex
created: 2026-11-17T02:14:51Z
updated: 2026-11-17T02:15:12Z
transitions:
  - to: ready
    at: 2026-11-17T02:15:02Z
    by: agent-S-0797
  - to: in-progress
    at: 2026-11-17T02:15:02Z
    by: agent-S-0797
  - to: done
    at: 2026-11-17T02:15:12Z
    by: agent-S-0797
stream: S-0797
tags: []
touches: [flai/cmd/verify_test.go]
usage:
  source: log
  seconds: 10
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 4
      output: 1498
      cache_read: 342369
      cache_write: 12599
      cost: 0.1832
---
# T-3491 TestVerifyWritesItsProgressToStderrOnlyWhenAsked clears FLAI_VERIFY_PROGRESS before its run without it, so it passes in a close-out job (I-0395)

## Work

S-0797's close-out stopped at integration on `flai/cmd/verify_test.go`: flai serve's close-out job sets `FLAI_VERIFY_PROGRESS=1`, and the test's run without the variable inherits it (I-0395). Clear it with `t.Setenv(closeout.ProgressEnv, "")` before that run. It waits for no task: it is a fix to main's test that S-0797's close-out needs.

## Done when

- The test passes with `FLAI_VERIFY_PROGRESS=1` in its environment and without it.
- S-0797's close-out passes the integration tier.

## Notes
