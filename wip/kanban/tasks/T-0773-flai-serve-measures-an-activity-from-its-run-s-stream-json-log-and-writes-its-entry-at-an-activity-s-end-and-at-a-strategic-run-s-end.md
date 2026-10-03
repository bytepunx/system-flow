---
id: T-0773
type: task
nature: feature
title: flai serve measures an activity from its run's stream-json log and writes its entry, at an activity's end and at a strategic run's end
status: done
parent: S-0206
owner: alex
created: 2026-10-03T18:37:41Z
updated: 2026-10-03T18:59:43Z
transitions:
  - to: ready
    at: 2026-10-03T18:38:46Z
    by: agent-S-0206
  - to: in-progress
    at: 2026-10-03T18:46:11Z
    by: agent-S-0206
  - to: done
    at: 2026-10-03T18:59:43Z
    by: agent-S-0206
stream: S-0206
tags: []
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/usage]
after: [T-0772]
usage:
  source: log
  seconds: 812
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 98
      output: 1260
      cache_read: 4678806
      cache_write: 123575
      cost: 1.9505
---
# T-0773 flai serve measures an activity from its run's stream-json log and writes its entry, at an activity's end and at a strategic run's end

## Work

In `flai/internal/serve`, the logs of a strategic agent's runs are `<serve dir>/agents/<key>-<kind>-*.log`. Logging an activity reads them, takes the activity's span from the later of the newest run's start and the document's last entry to now (or to the run's end when the run has ended), apportions the run's usage to the span as a task's is (ADR-0051, `usage.Record.Tasks` with one span), and appends the entry through T-0772's model, accruing the totals. When a strategic run ends, serve logs the time since the last entry as one activity, with the run's last result text as its summary, unless nothing was spent in it. `flai/internal/usage` gains the run's last result text where serve needs it.

Waits for T-0772: it writes through the activity document model.

## Done when

- Tests on fixture logs cover an entry written with its seconds and apportioned cost, totals accrued over two activities in one run, and a resumed run (a second log of the same session reporting cumulative totals) charged only its own share
- A run that ends after its agent logged its last activity adds no empty entry

## Notes
