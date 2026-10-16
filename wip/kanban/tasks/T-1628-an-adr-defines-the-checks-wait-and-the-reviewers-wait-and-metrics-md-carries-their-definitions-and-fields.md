---
id: T-1628
type: task
nature: feature
title: An ADR defines the checks wait and the reviewers wait, and metrics.md carries their definitions and fields
status: backlog
parent: S-0396
owner: alex
created: 2026-10-16T04:38:56Z
updated: 2026-10-16T04:38:56Z
transitions: []
stream: S-0396
tags: [cli]
touches: [design/adrs, design/system/metrics.md]
---
# T-1628 An ADR defines the checks wait and the reviewers wait, and metrics.md carries their definitions and fields

## Work

- Write the ADR with `flai adr new`, refining ADR-0091 and ADR-0114. It defines the two waits as the story's criterion 1 says, from S-0394's kept events, within the story's `review` intervals.
- It also covers:
  - how a push of a new head during a wait restarts the checks wait
  - how a story sent back ends its wait
  - the fields `items[].wait_checks_seconds` and `items[].wait_reviewers_seconds`, `waiting.weeks[].checks` and `.reviewers`, and the longest waits' kinds `checks` and `reviewers`
- Write the same into `design/system/metrics.md` § Waiting and § Charts, linking the ADR.
- First in its layer: it waits for no task. The definition is the contract the code and the chart follow.

## Done when

The ADR is accepted, `metrics.md` defines every new field, and `flai check --strict` passes.

## Notes
