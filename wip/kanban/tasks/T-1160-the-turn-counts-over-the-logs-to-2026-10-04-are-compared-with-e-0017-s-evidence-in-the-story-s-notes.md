---
id: T-1160
type: task
nature: feature
title: The turn counts over the logs to 2026-10-04 are compared with E-0017's evidence in the story's notes
status: done
parent: S-0293
owner: alex
created: 2026-10-07T09:28:30Z
updated: 2026-10-07T09:44:20Z
transitions:
  - to: ready
    at: 2026-10-07T09:28:46Z
    by: agent-S-0293
  - to: in-progress
    at: 2026-10-07T09:37:06Z
    by: agent-S-0293
  - to: done
    at: 2026-10-07T09:44:20Z
    by: agent-S-0293
stream: S-0293
tags: []
touches: [wip/kanban/stories/S-0293-flai-stats-classifies-a-story-run-s-tool-calls-so-the-ceremony-turns-e-0017-removes-are-measured-per-story-and-over-time.md]
after: [T-1157]
usage:
  source: log
  seconds: 434
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 44
      output: 20479
      cache_read: 2905493
      cache_write: 86197
      cost: 1.5519
---
# T-1160 The turn counts over the logs to 2026-10-04 are compared with E-0017's evidence in the story's notes

## Work

Build flai from the branch and measure every story flai serve kept a log for, without writing, against the operator's serve folder (`flai serve agent usage --all --json --config ~/.flai/config.json`). Sum each class over the days to 2026-10-04 and compare with the evidence in E-0017 (about 2,900 ceremony turns, 1,449 hand test runs, 650 empty wakes, 253 hand edits, 9% of turns pure ceremony). Record the comparison in the story's `## Notes`, with what each of the epic's figures counted where it differs from a turn. Waits for T-1157, whose classification it runs; it changes no file of the branch.

## Done when

- The story's notes hold a table of each class, the count to 2026-10-04, the epic's figure, and the difference in percent, and say why any difference over ten percent remains.

## Notes
