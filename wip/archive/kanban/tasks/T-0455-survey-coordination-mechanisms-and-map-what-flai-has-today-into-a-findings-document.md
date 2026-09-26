---
id: T-0455
type: task
nature: research
title: Survey coordination mechanisms and map what flai has today into a findings document
status: done
parent: S-0124
owner: alex
created: 2026-09-26T07:18:55Z
updated: 2026-09-26T07:24:23Z
transitions:
  - to: ready
    at: 2026-09-26T07:18:58Z
    by: agent-S-0124
  - to: in-progress
    at: 2026-09-26T07:18:59Z
    by: agent-S-0124
  - to: done
    at: 2026-09-26T07:24:23Z
    by: agent-S-0124
stream: S-0124
tags: []
touches: [design/system]
---
# T-0455 Survey coordination mechanisms and map what flai has today into a findings document

## Work

Survey how multi-agent coding tools, research, merge queues, conflict prediction, and concurrency control decide whether two pieces of work can run in parallel, and map what flai already has (touches, WIP limit, pull order, worktrees, flai serve start-up). Write the findings as a design document in `design/system/` with the options ranked.

## Done when

- A findings document in `design/system/` names each candidate mechanism with how it works, its fit for flai, and its cost, and cites its sources.
- It says what flai does today and where a hold would plug in.

## Notes
