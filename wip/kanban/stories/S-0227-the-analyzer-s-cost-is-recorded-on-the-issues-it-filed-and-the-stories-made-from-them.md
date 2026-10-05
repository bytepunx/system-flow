---
id: S-0227
type: story
nature: improvement
title: The analyzer's cost is recorded on the issues it filed and the stories made from them
status: backlog
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:18Z
updated: 2026-10-05T04:38:16Z
transitions: []
tags: [flai]
touches: [flai/internal/usage, flai/internal/serve, flai/internal/workitem, design/system/metrics.md, flai/internal/issues, flai/internal/metrics, design/system/continuous-improvement.md, design/adrs]
after: [S-0223, S-0226]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 60.98
  by: planner-E-0016
  at: 2026-10-04T04:48:23Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T20:47:00Z
  basis: "Its own forecast of 1h15m; 18th in the pull order with an in-progress limit of 3, behind S-0257, S-0244, S-0258, S-0276, S-0212, S-0213, S-0214, S-0215, S-0216, S-0217, S-0218, S-0219, S-0220, S-0221, S-0222, S-0223, S-0224 and S-0226."
  by: flai
  at: 2026-10-05T04:38:16Z
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

## Notes

- Rewritten on TH-0107: the first version restated S-0225's criteria, which S-0225 delivers for every kind. It waits for S-0226 because S-0226 adds the project total.
