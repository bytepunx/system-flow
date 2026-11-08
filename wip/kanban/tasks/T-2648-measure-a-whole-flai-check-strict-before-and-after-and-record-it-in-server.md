---
id: T-2648
type: task
nature: improvement
title: Measure a whole flai check --strict before and after and record it in server-performance.md and component-library.md
status: done
parent: S-0601
owner: alex
created: 2026-11-07T23:49:12Z
updated: 2026-11-08T01:59:18Z
transitions:
  - to: ready
    at: 2026-11-08T01:45:22Z
    by: agent-S-0601
  - to: in-progress
    at: 2026-11-08T01:45:26Z
    by: agent-S-0601
  - to: done
    at: 2026-11-08T01:59:18Z
    by: agent-S-0601
stream: S-0601
tags: [cli]
touches: [design/system/server-performance.md, design/system/component-library.md]
after: [T-2633, T-2643]
usage:
  source: log
  seconds: 832
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 22695
      cache_read: 10018401
      cache_write: 338834
      cost: 4.4171
---
# T-2648 Measure a whole flai check --strict before and after and record it in server-performance.md and component-library.md

## Work

Measure a whole `flai check --strict` on this repository, in the story's worktree, with a flai built from the main branch (before) and one built from the story's branch (after). For each, take the wall time cold, after a first run since the build, and warm, a second run with nothing changed, with `node` on `PATH`. Take several runs of each, and give the host's load average beside them, since server-performance.md says a loaded host inflates wall times two to three times.

Record the figures in `design/system/server-performance.md`, under cause 11 of § A second pass at 596 stories (S-0596) › Where the time goes, or in a section after it, as the sections for causes 1 to 7 record theirs. Say what changed: the rules at once, and `node` started only for a changed component.

In `design/system/component-library.md`, § Skeletons (S-0449) and § Examples (S-0443) › The check say how props are read. Add the kept read: where it lives, its key, and that `node` starts only for a component whose key changed.

This task waits for T-2633 and T-2643: it measures what both change, and documents T-2643's key.

## Done when

- The warm whole check takes under 2 s on this repository, or the task's Notes say what keeps it over and by how much.
- server-performance.md records the before and after, and component-library.md describes the kept read.

## Notes

- The warm whole check is over 2 s on this host: 3.3 to 4.9 s wall and 3.8 to 3.9 s CPU, at a load average of 84 to 103 on 24 cores, against 7.9 to 11.7 s wall and 4.9 to 5.3 s CPU on `main` with S-0600. So it is 1.3 to 2.9 s over. What keeps it over is the critical path the rules at once cannot shorten: the item list (0.66 to 1.10 s), `workItems` alone (about 0.45 to 0.8 s), and the `digest` rule (0.76 to 1.77 s), measured in process at a load of 103 to 114. Recorded as I-0312, with the figures in server-performance.md § Cause 11.
- Measured in the main checkout, not the story's worktree: the worktree has no `flaiover/node_modules` and no `.flai-cache`, so a check there never starts `node` and keeps nothing.
