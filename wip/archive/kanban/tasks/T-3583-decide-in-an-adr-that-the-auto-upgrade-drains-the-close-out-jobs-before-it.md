---
id: T-3583
type: task
nature: improvement
title: Decide in an ADR that the auto-upgrade drains the close-out jobs before it restarts and that serve asks again once a job a restart interrupted
status: done
parent: S-0846
owner: alex
created: 2026-11-17T12:33:22Z
updated: 2026-11-27T22:37:02Z
transitions:
  - to: ready
    at: 2026-11-27T22:35:48Z
    by: agent-S-0846
  - to: in-progress
    at: 2026-11-27T22:35:48Z
    by: agent-S-0846
  - to: done
    at: 2026-11-27T22:37:02Z
    by: agent-S-0846
stream: S-0846
tags: [adr, serve]
touches: [design/system/flai-cli.md]
usage:
  source: log
  seconds: 74
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 1956
      cache_read: 1820631
      cache_write: 19245
      cost: 0.5221
---
# T-3583 Decide in an ADR that the auto-upgrade drains the close-out jobs before it restarts and that serve asks again once a job a restart interrupted

## Work

Propose the solution I-0388 asks for, from its instances, before any code. Write a new ADR with `flai adr new` that refines ADR-0158 §5 and ADR-0229, and update `design/system/flai-cli.md` § Commands › The close-out job to match. The proposal the plan assumes:

1. **Drain before the restart.** While a project's verify pass has a job running, of either kind (close-out or rebase and resolve, which share `verifyPass.running`), `flai serve` asks for no upgrade, as it waits for a dashboard write today (`upgrade.go` `look`, `writing`). While an upgrade is due it takes no new request: the request file stays in `.flai-cache/verify/jobs/` for the next serve to take, so no new job starts that the restart would cut. The asker's ten-second wait must still be answered, so the pass answers the request with a pending state rather than letting it be withdrawn; say which in the ADR.
2. **A bound.** The drain waits at most a fixed time from the first due ask, 45 minutes (the longest job recorded is the 33-minute rebase and resolve of S-0824; the close-out p50 is 474 s), then upgrades and interrupts what still runs. No new manifest setting.
3. **Asked again once.** When serve starts and `sweep` ends a job left running `interrupted`, it also writes a request for the same story, kind, message, `all`, and `fresh`, marked as asked again, so the job starts again without the agent; a job that was itself asked again is not asked again a second time. A rebase and resolve is asked again only while its story is still in review.

Record the alternatives: waiting without a bound, re-exec while jobs run detached and adopting them, and resuming a job mid-tier. Leave the kluster, merge-queue, and image jobs out of scope and say so.

This task waits for nothing: it is the first layer.

## Done when

- A new ADR, status proposed or accepted as the operator decides, states the drain, its bound, and the ask-again rule, and refines ADR-0158 and ADR-0229.
- `design/system/flai-cli.md` § The close-out job says what a restart does to a running job under the ADR.
- `flai check --strict` passes.

## Notes

The ADR's file name is not known until `flai adr new` numbers it; `flai task done` widens this task's touches with it.
