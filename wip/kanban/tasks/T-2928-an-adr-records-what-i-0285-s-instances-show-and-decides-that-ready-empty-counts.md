---
id: T-2928
type: task
nature: improvement
title: An ADR records what I-0285's instances show and decides that ready-empty counts only ready stories not held on after and that a move to ready warns of an after hold
status: backlog
parent: S-0673
owner: alex
created: 2026-11-10T08:49:04Z
updated: 2026-11-10T08:49:04Z
transitions: []
stream: S-0673
tags: [flai, adr]
touches: [design/adrs/README.md]
---
# T-2928 An ADR records what I-0285's instances show and decides that ready-empty counts only ready stories not held on after and that a move to ready warns of an after hold

## Work

Read the story's `### Planning` notes, I-0285, and Finding 2 of `design/analysis/2026-11-07-all.md`. Check the planner's reading of the instances against `flai stats --json` (`items[].held_seconds`, `claims.weeks[].held_seconds`, `claims.days[]`) and the archived stories' transitions, and change the remedy where the evidence says otherwise.

Write the ADR with `flai adr new`, refining ADR-0211 §1 and ADR-0230. Its Context gives the instances: the holds on stories planned before ADR-0230, the hand moves to ready of stories held on `after`, and `readyDrained` counting held stories so that no `ready-empty` was told while lanes stood empty. Its Decision: `ready-empty` counts the ready stories not held on `after`, and a story held on overlap alone still counts; a move of a story to ready warns, and does not refuse, when a story it names in `after` is not done. Its Alternatives name the ones left: refusing the move, counting no held story, and changing the ready WIP limit's count, which `metrics.md` defines.

This task waits for nothing; the other three build what it decides.

## Done when

- The ADR is written and accepted, and `design/adrs/README.md` lists it.
- `flai check --strict` is clean.

## Notes
