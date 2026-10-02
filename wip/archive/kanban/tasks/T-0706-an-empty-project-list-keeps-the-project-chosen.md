---
id: T-0706
type: task
nature: remediation
title: An empty project list keeps the project chosen
status: done
parent: S-0178
owner: arobson
created: 2026-10-02T16:46:44Z
updated: 2026-10-02T16:50:51Z
transitions:
  - to: ready
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: in-progress
    at: 2026-10-02T16:47:08Z
    by: agent-S-0178
  - to: done
    at: 2026-10-02T16:50:51Z
    by: agent-S-0178
stream: S-0178
tags: []
touches: [flaiover/src/lib/project.svelte.ts, flaiover/src/lib/project.svelte.test.ts]
usage:
  source: log
  seconds: 223
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 78
      output: 31845
      cache_read: 2665332
      cache_write: 114981
      cost: 1.8918
---
# T-0706 An empty project list keeps the project chosen

## Work

`projectState.load()` picks no project when `/api/projects` answers a list without the current one, and an empty list does that too, which happens while flai serve is away. The layout keys the page on `projectState.current`, so the page remounts and the window goes to the top, and again when the list comes back. Keep the chosen project when the list is empty, as the catch already does when the request fails; a non-empty list without it still drops it.

Waits for nothing: no other task touches the project state.

## Done when

- A test shows an empty list after one with the chosen project keeps `current`; it fails without the change.
- A non-empty list without the chosen project still drops it, and the existing tests pass.

## Notes
