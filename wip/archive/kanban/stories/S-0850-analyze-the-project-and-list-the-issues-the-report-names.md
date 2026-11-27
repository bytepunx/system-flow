---
id: S-0850
type: story
nature: improvement
title: Analyze the project and list the issues the report names
status: done
owner: alex
created: 2026-11-23T06:00:22Z
updated: 2026-11-27T22:51:38Z
transitions:
  - to: ready
    at: 2026-11-23T06:00:23Z
    by: flai-serve
  - to: in-progress
    at: 2026-11-27T22:34:01Z
    by: agent-S-0850
  - to: review
    at: 2026-11-27T22:50:03Z
    by: agent-S-0850
  - to: done
    at: 2026-11-27T22:51:38Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0850-analyze-the-project-and-list-the-issues-the-report-names.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 985
  priced_by: harness
  turns:
    - day: 2026-11-27
      ceremony: 4
      work: 53
  models:
    - model: claude-opus-5-5
      input: 118
      output: 8373
      cache_read: 9571875
      cache_write: 185557
      cost: 3.5668
  strategic:
    - kind: orchestrator
      seconds: 79
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 24
          output: 178
          cache_read: 989316
          cache_write: 79336
          cost: 0.3396
        - model: claude-sonnet-5-5
          input: 2
          output: 27
          cache_read: 0
          cache_write: 13936
          cost: 0.0114
formula: analysis-cadence
---
# S-0850 Analyze the project and list the issues the report names

## Goal

Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance.

## Acceptance criteria

- [x] Run the analyzer and write its report
- [x] List the open issues the analyzer's report names

## Tasks
- T-3611 Run the analyzer and write its report
- T-3612 List the open issues the analyzer's report names

## Notes

Made by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=`.

### Accepted by the orchestrator

- Verified: e38a7e6cf74fff43677290229abc7432e69fffd1
- At: 2026-11-27T22:51:38Z

Verdict: pass
- 1: T-3611 done
- 2: T-3612 done
