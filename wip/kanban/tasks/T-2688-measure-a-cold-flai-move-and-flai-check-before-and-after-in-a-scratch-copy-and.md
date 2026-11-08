---
id: T-2688
type: task
nature: improvement
title: Measure a cold flai move and flai check before and after in a scratch copy, and record them in server-performance.md
status: in-progress
parent: S-0600
owner: alex
created: 2026-11-07T23:51:43Z
updated: 2026-11-08T01:02:17Z
transitions:
  - to: ready
    at: 2026-11-08T01:02:16Z
    by: agent-S-0600
  - to: in-progress
    at: 2026-11-08T01:02:17Z
    by: agent-S-0600
stream: S-0600
tags: [cli, docs]
touches: [design/system/server-performance.md]
after: [T-2682, T-2685]
---
# T-2688 Measure a cold flai move and flai check before and after in a scratch copy, and record them in server-performance.md

## Work

Build flai from `main` and from the story's branch, as S-0596 did (`## A second pass at 596 stories (S-0596)` › `### How it was measured` in `design/system/server-performance.md`). In a scratch copy of the main checkout under the session's scratch folder, with a `.flai-cache` in it, time each binary three times, alternating, on:

- a `flai move` of a task to and fro, with no kept file (the first run) and with one;
- a `flai check`, with no kept file and with one.

Take the median wall time, and note the host's load, since S-0596 found load inflated wall times two to three times. Write a section after S-0531's, `## Parses and lint findings kept between processes (S-0600)`, with the method, a table of before and after, and what is left; mark cause 10's row as remedied by S-0600. Never write a time you did not observe.

Third layer: it waits for T-2682 and T-2685, since both caches must be in the binary it measures.

## Done when

- [ ] `design/system/server-performance.md` records the before and after of a cold `flai move` and a cold `flai check`, each with and without a kept file.
- [ ] Each cold figure with a kept file is at least 1 s less than before, or the story's agent has said on a thread on S-0600 why it is not.
- [ ] The markdown lint and `flai check --strict` pass on the document.

## Notes
