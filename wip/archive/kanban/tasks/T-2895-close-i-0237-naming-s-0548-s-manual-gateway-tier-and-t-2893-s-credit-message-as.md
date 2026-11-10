---
id: T-2895
type: task
nature: remediation
title: Close I-0237 naming S-0548's manual gateway tier and T-2893's credit message as what fixed it
status: done
parent: S-0685
owner: alex
created: 2026-11-10T08:44:36Z
updated: 2026-11-10T09:10:31Z
transitions:
  - to: ready
    at: 2026-11-10T08:49:22Z
    by: agent-S-0685
  - to: in-progress
    at: 2026-11-10T09:02:32Z
    by: agent-S-0685
  - to: done
    at: 2026-11-10T09:10:31Z
    by: agent-S-0685
stream: S-0685
tags: [gateway, issues]
touches: [design/issues/I-0237-the-smoke-tier-s-openrouter-gateway-check-fails-on-a-402-because-the-operator-s.md, design/issues/summary.md]
after: [T-2893]
usage:
  source: log
  seconds: 479
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 14
      output: 1622
      cache_read: 1932515
      cache_write: 22667
      cost: 0.6003
  strategic:
    - kind: orchestrator
      seconds: 5
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 2
          output: 16
          cache_read: 114072
          cache_write: 3023
          cost: 0.0317
---
# T-2895 Close I-0237 naming S-0548's manual gateway tier and T-2893's credit message as what fixed it

## Work

Waits for T-2893, which the close reason names as part of the fix.

- Check criterion 1: run `flai test flai/internal/verify/select_test.go`. `TestRepositoryManifestRunsTheGatewayChecksOnlyByName` and the manual-tier selection test there hold that no close-out, `--all`, or merge-queue batch selects the `gateway` tier, so a key out of credit no longer fails one. They are the test that reproduces I-0237's cause. Tick criterion 1 with `flai criteria tick S-0685 1` once both pass and T-2893 is done.
- Run `flai issue close I-0237 --reason "<reason>"`. The reason says that S-0548 made the gateway checks the manual `gateway` tier, held by `TestRepositoryManifestRunsTheGatewayChecksOnlyByName`, and that T-2893 makes the check name a 402 when it runs by name. Then tick criterion 2.

## Done when

- I-0237's status is closed with that reason, and `design/issues/summary.md` no longer lists it as open.
- Both of S-0685's criteria are ticked.

## Notes
