---
id: T-0938
type: task
nature: feature
title: The design, the user guide, and the flai reference describe the Orchestrator and Analyzer pages and the orchestrator's stop and start
status: done
parent: S-0228
owner: alex
created: 2026-10-05T05:45:26Z
updated: 2026-10-07T00:43:51Z
transitions:
  - to: ready
    at: 2026-10-07T00:36:34Z
    by: agent-S-0228
  - to: in-progress
    at: 2026-10-07T00:36:35Z
    by: agent-S-0228
  - to: done
    at: 2026-10-07T00:43:51Z
    by: agent-S-0228
stream: S-0228
tags: [dashboard]
touches: [design/system/flaiover-dashboard.md, docs/users/flaiover.md, design/system/dashboard-host-channel.md, design/system/flai-cli.md, docs/users/flai-reference.md, design/system/strategic-agents.md, docs/users/flai.md, docs/operators/index.md]
after: [T-0917, T-0921, T-0928, T-0931]
usage:
  source: log
  seconds: 436
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 143
      output: 63425
      cache_read: 9740562
      cache_write: 235782
      cost: 4.5401
---
# T-0938 The design, the user guide, and the flai reference describe the Orchestrator and Analyzer pages and the orchestrator's stop and start

## Work

- `design/system/flaiover-dashboard.md`: the two pages among its views, the shared strategic agent panel, and the API rows `/api/orchestrator`, `/api/analyzer`, and `/api/agent-stream/[story]?role`.
- `docs/users/flaiover.md`: Orchestrator and Analyzer sections after Planner, saying what each page shows, how to stop and start the orchestrator, and how to run the analyzer with a focus.
- `design/system/dashboard-host-channel.md` and `design/system/flai-cli.md`: `agent.stream`'s `role`, the writes `orchestrate.stop` and `orchestrate.start`, and the hold they set, in the method tables.
- `docs/users/flai-reference.md`: `flai serve orchestrate stop|start`.

It waits for every other task of the story, whose behaviour it describes as built.

## Done when

- [ ] Each document describes the pages and the commands as built, with `updated` bumped
- [ ] The markdown lint passes on each

## Notes
