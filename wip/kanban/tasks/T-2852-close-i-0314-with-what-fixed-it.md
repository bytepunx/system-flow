---
id: T-2852
type: task
nature: remediation
title: Close I-0314 with what fixed it
status: backlog
parent: S-0653
owner: alex
created: 2026-11-10T07:02:00Z
updated: 2026-11-10T07:02:00Z
transitions: []
stream: S-0653
tags: [issues]
touches: [design/issues/I-0314-testaformularolestepsendistoldtotheinstancesagent-fails-on-its-timing-under.md, design/issues/summary.md]
after: [T-2849, T-2851]
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
