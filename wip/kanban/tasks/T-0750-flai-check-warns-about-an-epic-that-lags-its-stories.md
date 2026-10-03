---
id: T-0750
type: task
nature: improvement
title: flai check warns about an epic that lags its stories
status: done
parent: S-0200
owner: alex
created: 2026-10-03T07:11:22Z
updated: 2026-10-03T07:32:01Z
transitions:
  - to: ready
    at: 2026-10-03T07:11:52Z
    by: agent-S-0200
  - to: in-progress
    at: 2026-10-03T07:19:14Z
    by: agent-S-0200
  - to: done
    at: 2026-10-03T07:32:01Z
    by: agent-S-0200
stream: S-0200
tags: []
touches: [flai/internal/check]
after: [T-0745]
usage:
  source: log
  seconds: 767
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 16022
      cache_read: 2772981
      cache_write: 62136
      cost: 1.2549
---
# T-0750 flai check warns about an epic that lags its stories

## Work

Add an advisory warning, `epic.lags-stories`, on an open epic whose stories put it further on than its status, naming the moves that catch it up, or `flai accept` when its stories are all done. It is advisory because the epic is the operator's to move, not a story's agent's, so `--strict` passes over it.

It waits for T-0745, whose `EpicFollows` it uses.

## Done when

- [x] `flai check` reports the warning and `--strict` passes over it
- [x] Tests in `flai/internal/check` cover an epic that lags and one that does not

## Notes
