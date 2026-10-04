---
id: T-0805
type: task
nature: improvement
title: A planner activity's apportioned usage is charged to the item it planned and the items above it
status: done
parent: S-0225
owner: alex
created: 2026-10-04T04:07:28Z
updated: 2026-10-04T04:21:10Z
transitions:
  - to: ready
    at: 2026-10-04T04:07:53Z
    by: agent-S-0225
  - to: in-progress
    at: 2026-10-04T04:12:43Z
    by: agent-S-0225
  - to: done
    at: 2026-10-04T04:21:10Z
    by: agent-S-0225
stream: S-0225
tags: []
touches: [flai/internal/serve, design/system/strategic-agents.md, design/system/flai-cli.md]
after: [T-0803]
usage:
  source: log
  seconds: 507
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 88
      output: 626
      cache_read: 5180606
      cache_write: 126340
      cost: 2.1646
---
# T-0805 A planner activity's apportioned usage is charged to the item it planned and the items above it

## Work

When a planner activity is logged, by `LogRunEnd` as its run ends or by `LogActivity` for the MCP tool `activity_log` during a run, charge the usage apportioned to the activity's span (the same `spent(s)` its entry's cost comes from) to the item the run was started for, under the planner's `strategic` entry, and to the items above it. The item is the one `serve/agents.json` records under `plans` for the run whose log is the newest. An activity outside any run flai serve logged charges nothing. `Measure`, which writes a story's usage from its agents' logs, keeps the `strategic` the story carries.

`design/system/strategic-agents.md` and `design/system/flai-cli.md` say where the planner's cost lands.

Waits for the usage task, whose types and charging it uses. Shares no path with the stats task, so the two run together, this one in a worktree of its own since `flai/cmd` builds both.

## Done when

- [ ] A planner run that ends charges its planned story and that story's epic with the activity's apportioned tokens, cost, and seconds, `estimated: true`, and the charge equals the activity entry's cost
- [ ] Tests pin the apportioning: two activities in one run each charge only their span's share, a resumed session's cumulative result is not charged twice, and an activity outside a logged run charges nothing
- [ ] `Measure` keeps a story's `strategic`, with a test
- [ ] `go test ./internal/serve/...` passes
- [ ] The design documents say where the planner's cost lands

## Notes
