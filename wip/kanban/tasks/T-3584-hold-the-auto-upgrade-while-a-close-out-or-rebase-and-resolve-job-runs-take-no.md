---
id: T-3584
type: task
nature: improvement
title: Hold the auto-upgrade while a close-out or rebase and resolve job runs, take no new job meanwhile, and upgrade once the bound passes
status: backlog
parent: S-0846
owner: alex
created: 2026-11-17T12:33:31Z
updated: 2026-11-17T12:33:31Z
transitions: []
stream: S-0846
tags: [serve, code]
touches: [flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/verifyjob.go, flai/internal/serve/verifyjob_test.go, flai/internal/serve/serve.go]
after: [T-3583]
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
