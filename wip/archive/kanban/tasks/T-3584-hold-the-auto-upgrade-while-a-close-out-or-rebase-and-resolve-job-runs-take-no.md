---
id: T-3584
type: task
nature: improvement
title: Hold the auto-upgrade while a close-out or rebase and resolve job runs, take no new job meanwhile, and upgrade once the bound passes
status: done
parent: S-0846
owner: alex
created: 2026-11-17T12:33:31Z
updated: 2026-11-27T23:19:55Z
transitions:
  - to: ready
    at: 2026-11-27T22:37:15Z
    by: agent-S-0846
  - to: in-progress
    at: 2026-11-27T22:37:15Z
    by: agent-S-0846
  - to: done
    at: 2026-11-27T23:19:55Z
    by: agent-S-0846
stream: S-0846
tags: [serve, code]
touches: [flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/serve.go, flai/internal/closeout/closeout.go, flai/internal/closeout/closeout_test.go, flai/internal/mcpserver/verify.go, flai/internal/mcpserver/verify_test.go, flai/cmd/verify.go, flai/internal/serve/statejob.go, design/adrs/0320-the-auto-upgrade-drains-the-verify-pass-s-jobs-before-it-restarts-holds-new.md, design/system/flai-cli.md, flai/cmd/accept_resolve.go, flai/cmd/accept_resolve_test.go, flai/cmd/verify_test.go]
after: [T-3583]
usage:
  source: log
  seconds: 2560
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 11278
      cache_read: 10497739
      cache_write: 110969
      cost: 3.0105
---
# T-3584 Hold the auto-upgrade while a close-out or rebase and resolve job runs, take no new job meanwhile, and upgrade once the bound passes

## Work

Build the drain the ADR of T-3583 decides; it waits for T-3583 because the bound, the pending answer to a request, and what counts as a running job are that ADR's.

- Give the verify pass (`verifyjob.go`) a count of its running jobs, both kinds, read under its lock, and a drain flag: while it is set, `look` leaves a new request in place, answered as the ADR says, and starts no job.
- In `upgrade.go` `look`, ask for the upgrade only when `writing()` is 0 and no project's verify pass has a job running, or once the bound has passed since the first due ask; set the drain while it waits, and log `flai upgrade waits for a close-out job` once per wait, as it logs the dashboard write.
- Wire the count and the drain from the verify passes to the upgrade in `serve.go`, where the pass is built.
- Tests that reproduce I-0388: an upgrade due while a close-out job runs is not asked until the job ends; a request made while the drain holds starts no job and is taken by the next pass; the bound passed asks for the upgrade with a job still running.

## Done when

- `upgrade_test.go` and `verifyjob_test.go` hold the three tests above, and they pass with `flai test` on the changed paths.
- `TestTheUpgradeWaitsForADashboardWrite` and the existing verify job tests still pass.

## Notes
