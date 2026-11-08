---
id: T-2648
type: task
nature: improvement
title: Measure a whole flai check --strict before and after and record it in server-performance.md and component-library.md
status: in-progress
parent: S-0601
owner: alex
created: 2026-11-07T23:49:12Z
updated: 2026-11-08T01:45:26Z
transitions:
  - to: ready
    at: 2026-11-08T01:45:22Z
    by: agent-S-0601
  - to: in-progress
    at: 2026-11-08T01:45:26Z
    by: agent-S-0601
stream: S-0601
tags: [cli]
touches: [design/system/server-performance.md, design/system/component-library.md]
after: [T-2633, T-2643]
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
