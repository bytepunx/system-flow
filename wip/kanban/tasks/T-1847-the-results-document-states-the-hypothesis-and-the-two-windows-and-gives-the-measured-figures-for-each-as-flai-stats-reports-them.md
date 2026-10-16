---
id: T-1847
type: task
nature: experiment
title: The results document states the hypothesis and the two windows and gives the measured figures for each, as flai stats reports them
status: backlog
parent: S-0435
owner: alex
created: 2026-10-16T06:08:54Z
updated: 2026-10-16T06:08:54Z
transitions: []
stream: S-0435
tags: [cli]
touches: [design/experiments/S-0435-digest-priming-and-threads-to-the-operator.md]
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
