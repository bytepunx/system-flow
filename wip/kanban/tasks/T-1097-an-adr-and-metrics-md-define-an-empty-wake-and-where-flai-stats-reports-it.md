---
id: T-1097
type: task
nature: improvement
title: An ADR and metrics.md define an empty wake and where flai stats reports it
status: backlog
parent: S-0272
owner: alex
created: 2026-10-06T22:53:06Z
updated: 2026-10-06T22:53:06Z
transitions: []
stream: S-0272
tags: [metrics]
touches: [design/adrs, design/system/metrics.md]
---
# T-1097 An ADR and metrics.md define an empty wake and where flai stats reports it

## Work

Criterion 4, the contract. `design/system/metrics.md` is the contract between `flai stats` and the dashboard, and it changes only with an ADR.

- Write the ADR with `flai adr new`. It defines an empty wake as a main-agent `wait_for_events` call whose result had `timed_out` true, no events and no changed paths. A call that answered `end` is not an empty wake. Empty wakes are counted per story from the run logs that `usage` already reads, and reported by `flai stats` per story and in total for the window, as text and `--json`.
- Say how this relates to S-0293, which classifies every turn of a run and has empty wakes as one of its classes. This count is that class, computed the same way, so S-0293 takes it over rather than defining it again.
- Add the definition to `metrics.md`, linking the ADR.

It waits for no task. The task that counts empty wakes in code waits for this one.

## Done when

- The ADR is written, proposed, with the definition, the source, and where it is reported.
- `metrics.md` defines the empty wake and links the ADR.
- The markdown lint and `flai check --strict` are clean.

## Notes

Drafted by the planner. `design/adrs` is a folder touch because the ADR's file name is not known until `flai adr new` numbers it.
