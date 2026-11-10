---
id: T-2852
type: task
nature: remediation
title: Close I-0314 with what fixed it
status: done
parent: S-0653
owner: alex
created: 2026-11-10T07:02:00Z
updated: 2026-11-10T07:37:58Z
transitions:
  - to: ready
    at: 2026-11-10T07:37:50Z
    by: agent-S-0653
  - to: in-progress
    at: 2026-11-10T07:37:52Z
    by: agent-S-0653
  - to: done
    at: 2026-11-10T07:37:58Z
    by: agent-S-0653
stream: S-0653
tags: [issues]
touches: [design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md, design/issues/summary.md]
after: [T-2849, T-2851]
usage:
  source: log
  seconds: 4
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 6
      output: 48
      cache_read: 589878
      cache_write: 2255
      cost: 0.2178
---
# T-2852 Close I-0314 with what fixed it

## Work

Second layer: it waits for T-2849 and T-2851, because the reason it writes names what both changed, and the issue closes only once both causes are gone.

- Run `go test ./internal/serve -run 'TestAFormulaRoleStepsEndIsToldToTheInstancesAgent' -count 20` in `flai/`, with the host under the load other stories give it if it has any, and record the run and the host's load in the task's Notes.
- Close the issue with `flai issue close I-0314 --reason "<reason>"`, the reason naming both causes and both fixes: `logRunEnd` measures the ended run's own log, and the test waits for each end in the inbox.

## Done when

- The 20 runs pass; the run and the load are in the task's Notes.
- I-0314 is closed with that reason, and `design/issues/summary.md` no longer lists it.

## Notes
