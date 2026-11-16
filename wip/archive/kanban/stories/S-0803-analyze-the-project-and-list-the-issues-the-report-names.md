---
id: S-0803
type: story
nature: improvement
title: Analyze the project and list the issues the report names
status: done
owner: alex
created: 2026-11-16T06:00:53Z
updated: 2026-11-16T06:24:19Z
transitions:
  - to: ready
    at: 2026-11-16T06:00:54Z
    by: flai-serve
  - to: in-progress
    at: 2026-11-16T06:01:10Z
    by: agent-S-0803
  - to: review
    at: 2026-11-16T06:19:27Z
    by: agent-S-0803
  - to: done
    at: 2026-11-16T06:24:19Z
    by: orchestrator
tags: []
touches: [wip/kanban/stories/S-0803-analyze-the-project-and-list-the-issues-the-report-names.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1107
  priced_by: harness
  turns:
    - day: 2026-11-16
      ceremony: 1
      work: 31
  models:
    - model: claude-opus-5-5
      input: 68
      output: 9418
      cache_read: 4834413
      cache_write: 168070
      cost: 2.5001
  strategic:
    - kind: orchestrator
      seconds: 41
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 12
          output: 123
          cache_read: 357191
          cache_write: 66595
          cost: 0.1296
        - model: claude-sonnet-5-5
          input: 2
          output: 27
          cache_read: 0
          cache_write: 13926
          cost: 0.0111
formula: analysis-cadence
---
# S-0803 Analyze the project and list the issues the report names

## Goal

Run the analyzer on its focus, or on all three, and list the open issues its report names, for the operator to choose stories from at acceptance.

## Acceptance criteria

- [x] Run the analyzer and write its report
- [x] List the open issues the analyzer's report names

## Tasks
- T-3426 Run the analyzer and write its report
- T-3427 List the open issues the analyzer's report names

## Notes

Made by `flai formula run analysis-cadence` from `design/workflows/formulas/analysis-cadence.md`, with the inputs `focus=`.

The analyzer's run on S-0803 (focus all) ended with its work done and wrote `design/analysis/2026-11-16-all.md`. `flai issue list --report latest` names these open issues, none yet linked from a story:

| Issue | Class | Cost | Title |
|-------|-------|------|-------|
| I-0386 | efficiency | 30h | A story in review whose close-out fails after main moved waits for the operator to send it back, while its after chain holds the backlog and the board starves |
| I-0387 | efficiency | 20m | Stories of one epic rewrite the same section of one design document, so acceptances send back on same-line conflicts the section merge cannot settle |
| I-0388 | efficiency | 8m | flai serve's auto-upgrade restart interrupts the close-out jobs it was running, which start again from scratch |
| I-0389 | efficiency | 1m | flai serve opens a new board-starvation thread at each flip of the board, 48 in five days, each a thread to the operator |
| I-0390 | defect | 15m | Design sentences the code contradicts: ci.md's race detector on macOS, markdown.md's yaml in flaiover, metrics.md's WIP series and unnamed keys, flai-cli.md's missing dashboard logs row, the analyzer's MultiEdit |
| I-0391 | impression | 5m | Narrative, index, and notice-log writes take no lock, and docedit's atomic write does not sync before its rename |
| I-0392 | impression | 5m | flai file show follows a symlink under design or docs out of the repository: docedit/file.go confines by string prefix and os.Stat, unlike hostapi/docs.go's EvalSymlinks |
| I-0393 | impression | 5m | A formula input's value is substituted raw into a step's command line, and only ADR-0273's wording keeps the orchestrator's judgement shell-safe |
| I-0394 | impression | 5m | CI runs no vulnerability scan, and scripts/install-tools.sh installs npm tools without a lockfile or integrity check |

The operator chooses at acceptance which of them become stories.

### Accepted by the orchestrator

- Verified: ada0fbab71cb44ecb3db631663869cdf5cc13c94
- Batch: 20261116-062411
- Enqueued verified: ada0fbab71cb44ecb3db631663869cdf5cc13c94
- At: 2026-11-16T06:24:19Z

Verdict: S-0803 meets both criteria at ada0fbab71cb44ecb3db631663869cdf5cc13c94.
- 1: T-3426 done
- 2: T-3427 done
