---
id: S-0227
type: story
nature: improvement
title: The analyzer's cost is recorded on the issues it filed and the stories made from them
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-06T10:31:55Z
transitions:
  - to: ready
    at: 2026-10-05T06:13:41Z
    by: alex
tags: [flai]
topics: [analysis, planning]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/issues, flai/internal/metrics, design/system/continuous-improvement.md, design/adrs, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/cmd/stats.go, flai/cmd/check_stats_test.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [S-0223, S-0226]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 54.15
  by: planner-S-0227
  at: 2026-10-05T05:45:09Z
forecast:
  duration: 50m
  delivery: 2026-10-06T15:59:00Z
  basis: "Its own forecast of 50m; 7th in the pull order with an in-progress limit of 3, behind S-0285, S-0221, S-0222, S-0224, S-0226 and S-0223."
  by: flai
  at: 2026-10-06T10:31:55Z
---
# S-0227 The analyzer's cost is recorded on the issues it filed and the stories made from them

## Goal

S-0225 charges a planner activity's usage to the item it planned (ADR-0083), and the entry's `kind` leaves room for the analyzer. The analyzer's work ends in a report and in issues, but issues carry no usage today.

Its cost should land on the issues it filed or bumped. It should then travel with an issue to the draft story that the issue step makes from it, so that the cost of a remediation includes the cost of finding it.

## Acceptance criteria
- [ ] An issue carries `usage` with a `strategic` entry per kind, in the shape ADR-0083 gives items. An ADR extends ADR-0083 to issues, and `front-matter-fields.txt` lists the key
- [ ] An analyzer activity's usage is charged when the activity ends, apportioned as ADR-0083 apportions a planner's. It goes under the `analyzer` entry of the issues its `activity_log` call named, split evenly between them. An activity that names no issue is charged to the project total that S-0226 adds
- [ ] The story the issue step makes from an issue carries the issue's strategic usage as its own `strategic` entry, summed up to its epic. `flai stats` counts it once in totals: the issue's until a story is made from it, the story's after
- [ ] `flai stats` reports analyzer usage per issue and in totals, apart from the agents' figures. `design/system/metrics.md` and `design/system/continuous-improvement.md` describe it. Tests pin the charge, the split, and the carry-over

## Tasks
- T-0935 An ADR extends ADR-0083 to issues and to the story the issue step makes from one
- T-0940 An issue reads and writes usage with a strategic entry per kind, and front-matter-fields.txt lists the key
- T-0942 An analyzer activity's apportioned usage is charged evenly to the issues it named, or to the project total
- T-0944 The story the issue step makes from an issue carries the issue's strategic usage, summed up to its epic
- T-0948 flai stats reports analyzer usage per issue and in totals, counted once across an issue and its story

## Notes

- Rewritten on TH-0107: the first version restated S-0225's criteria, which S-0225 delivers for every kind. It waits for S-0226 because S-0226 adds the project total.

### Planning

- Declared touches, all kept: `flai/internal/usage`, `flai/internal/serve`, `flai/internal/workitem`, `flai/internal/issues`, `flai/internal/metrics`, `design/system/metrics.md`, `design/system/continuous-improvement.md`, and `design/adrs`.
- Touches added from the code layout: `flai/cmd/issue.go`, `flai/internal/mcpserver/issues.go`, and their tests. Both call `issues.ForStory` and make the story, so the carry-over lands there. The dashboard's `issue.story` runs the CLI and needs no change.
- Also from the code layout: `flai/cmd/stats.go` and `flai/cmd/check_stats_test.go`, the stats report and its tests, as in S-0225 and S-0226.
- Touches added from co-change (`flai touches suggest`): `design/system/flai-cli.md` (33%), `docs/users/flai.md` (31%), and `docs/users/flai-reference.md` (17%, the generated `flai stats` help). The dashboard is left out: the criteria ask only for `flai stats`.
- Topics `analysis` and `planning` added: the tasks reach the analyzer's design and ADR-0083.
- Forecast 50m, up from flai's 18m (86 s per unit of size, size 12). S-0225, the planner's version of this story, took 52m in progress. This one replaces S-0225's dashboard work with a new usage block in the issue schema and the carry-over to the story the issue step makes. It also adds an ADR. Delivery is flai's 19:18Z moved by the 32m added.
- Cost of delay 54.15 USD a week, as `flai cod` gives it after the forecast change: the story's share of E-0016's 1500 USD a week, 50m of 23h5m over the epic's 17 open stories without inputs. The story has no inputs of its own. It replaces planner-E-0016's 60.98, worked out from the old 1h15m forecast.
- Tasks in four layers: T-0935 (the ADR and the issue schema in the design), then T-0940 (the issue's `usage` and the charge function), then T-0942 (the analyzer's charge) beside T-0944 (the carry-over), then T-0948 (`flai stats`, which counts once by the carry-over).
