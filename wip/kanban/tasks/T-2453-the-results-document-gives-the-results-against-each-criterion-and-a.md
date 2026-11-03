---
id: T-2453
type: task
nature: experiment
title: The results document gives the results against each criterion and a recommendation, and the experiments index lists it
status: backlog
parent: S-0565
owner: alex
created: 2026-11-03T01:25:27Z
updated: 2026-11-03T01:25:27Z
transitions: []
stream: S-0565
tags: [docker]
touches: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/experiments/README.md]
after: [T-2451, T-2452]
---
# T-2453 The results document gives the results against each criterion and a recommendation, and the experiments index lists it

## Work

Waits for T-2451 and T-2452: it sums up what both runs recorded.

Finish the results document: the results against each of S-0565's three criteria, with the stop and restart that T-2450 wrote in S-0565's Notes linked rather than copied; the issues filed, by ID, each with the story of E-0029 that carries it on; and a recommendation to adopt, adapt, or drop running this repository's stories on `harness: container` by default, with what must change first. Add its row to the table in `design/experiments/README.md`, as S-0505's row has it. Tick each criterion of S-0565 you verified with `flai criteria tick`.

## Done when

- The results document passes `flai check` and the markdown lint, and has its results and its recommendation.
- `design/experiments/README.md` lists it with its recommendation.
- Each criterion of S-0565 the runs met is ticked, and any not met says why in S-0565's Notes.

## Notes
