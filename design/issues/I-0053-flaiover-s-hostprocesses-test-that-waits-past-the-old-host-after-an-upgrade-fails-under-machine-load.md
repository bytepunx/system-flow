---
id: I-0053
title: flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load
class: defect
status: closed
count: 1
cost: 4m
first_reported: 2026-10-01T11:35:40Z
last_reported: 2026-10-01T11:35:40Z
updated: 2026-10-06T22:53:08Z
---

# I-0053 flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load

## Description
flaiover's HostProcesses test that waits past the old host after an upgrade fails under machine load

## Instances

### 2026-10-01T11:35:40Z
S-0189's pre-review verifier: HostProcesses.svelte.test.ts ('waits past the old host after an upgrade that answered before restarting') failed in close-out while the Go tiers loaded the machine, then passed 3 of 3 alone and in the second full run. The branch does not touch the file or the component.

## Remediation

Story S-0264 remediates this issue, created from it at 2026-10-05T00:03:14Z.
Closed 2026-10-06T22:53:08Z: S-0264: HostProcesses' reconnect tests wait for their outcome (vi.waitFor, 4 s) with a 3 s give-up instead of asserting after a fixed 40 ms against a 50 ms give-up; a new test that answers the old host's poll after 80 ms fails at the old timings and passes now, and the file passed 6 of 6 runs while the Go tiers loaded the machine (load 10-13).
