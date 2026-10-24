---
id: T-1847
type: task
nature: experiment
title: The results document states the hypothesis and the two windows and gives the measured figures for each, as flai stats reports them
status: done
parent: S-0435
owner: alex
created: 2026-10-16T06:08:54Z
updated: 2026-10-24T23:48:14Z
transitions:
  - to: ready
    at: 2026-10-24T23:45:49Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T23:45:50Z
    by: system-flow
  - to: done
    at: 2026-10-24T23:48:14Z
    by: agent-S-0435
stream: S-0435
tags: [cli]
touches: [design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md]
usage:
  source: log
  seconds: 144
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 9224
      cache_read: 1479541
      cache_write: 38122
      cost: 0.7855
---
# T-1847 The results document states the hypothesis and the two windows and gives the measured figures for each, as flai stats reports them

## Work

- Find S-0430's acceptance time with `flai show S-0430 --json`. The before window is the 14 days before it, the after window the 14 days after.
- Run `scripts/flai.sh stats --json` over each window and take `digest.weeks[]`, `digest.entries`, and `waiting.weeks[].threads` as flai gives them.
- Write `design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md` from `design/experiments/template.md`: the hypothesis, the windows by date, the method with the commands, and one table per window of threads opened to the operator, stories completed, threads per completed story, threads that asked again, and the thread waits' total and mean.

First layer: the analysis task reads these figures.

## Done when

The document holds the hypothesis, the method, and the table with every figure taken from flai's output, and `flai check --strict` reports nothing on it.

## Notes

If the after window is not yet whole when the story starts, block it with `flai block --reason` until it is, rather than measuring a short window.
