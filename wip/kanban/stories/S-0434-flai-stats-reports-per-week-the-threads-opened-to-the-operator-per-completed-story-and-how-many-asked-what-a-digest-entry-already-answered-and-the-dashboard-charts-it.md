---
id: S-0434
type: story
nature: feature
title: flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it
status: backlog
parent: E-0025
owner: alex
created: 2026-10-16T06:00:19Z
updated: 2026-10-16T06:00:19Z
transitions: []
tags: [cli, dashboard]
topics: [cli, dashboard, analysis]
after: [S-0433]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0434 flai stats reports per week the threads opened to the operator per completed story and how many asked what a digest entry already answered, and the dashboard charts it

## Goal

The analyzer and the operator can see whether the digest cut the questions that wait on the operator (I-0151). `metrics.md` defines, by an ADR, the digest's measure, and `flai stats` reports it per ISO week of the window: the threads opened to the operator, the stories completed, the threads per completed story, and the threads that asked what an entry already answered, those an entry dated before the thread lists in its `asked_again`. The dashboard draws it as a chart beside the agent-waiting chart.

## Acceptance criteria

- [ ] An ADR, refining ADR-0081 and ADR-0091, defines in `design/system/metrics.md` a thread opened to the operator (opened by an author other than the operator and other than flai's own notices, such as the upgrade thread), the week it counts in (its first entry's), the completed stories of a week (as `waiting.weeks[]` counts them), the rate, and a repeat (a thread an entry lists in `asked_again` whose `date` is before the thread's first entry).
- [ ] `flai stats --json` carries `digest.weeks[]` with `week`, `start`, `threads_to_operator`, `stories_completed`, `per_story` (absent with no completed story), and `asked_again`, and `digest.entries` with the counts of active, compacted, and absorbed entries; `flai stats` prints one line per week under a `digest` heading; tests cover a week with repeats, a week with no completed story, and a project with no digest.
- [ ] The dashboard's charts page has a digest chart of threads per completed story per week with the repeats stacked, and a test renders it from a fixture, including a flai that does not report `digest`.
- [ ] `design/system/metrics.md` § Charts, `docs/users/flai.md`, and `docs/users/flaiover.md` describe the measure and the chart; the convention `strategic-agents.md` § As the analyzer names the measure for a bottlenecks run.

## Tasks

## Notes

- Drafted by the planner for E-0025, after S-0433 in the epic's order. It reads only S-0428's front matter, so it could follow S-0428; the plan thread proposes that, since the weeks before the digest are what the experiment (S-0435) compares against and are already in the threads.
- `metrics.md` is the contract between `flai stats` and the dashboard and changes only with an ADR (CLAUDE.md).
