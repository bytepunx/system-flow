---
id: T-2268
type: task
nature: improvement
title: A cold stats.get and flai stats --json are measured again on this repository with the kept file, and server-performance.md records the figures
status: done
parent: S-0531
owner: alex
created: 2026-10-28T18:24:40Z
updated: 2026-10-28T18:36:16Z
transitions:
  - to: ready
    at: 2026-10-28T18:26:34Z
    by: agent-S-0531
  - to: in-progress
    at: 2026-10-28T18:31:17Z
    by: agent-S-0531
  - to: done
    at: 2026-10-28T18:36:16Z
    by: agent-S-0531
stream: S-0531
tags: [flai, performance]
touches: [design/system/server-performance.md]
after: [T-2267]
usage:
  source: log
  seconds: 299
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 10605
      cache_read: 1724141
      cache_write: 66861
      cost: 1.0111
---
# T-2268 A cold stats.get and flai stats --json are measured again on this repository with the kept file, and server-performance.md records the figures

## Work

Measure with S-0500's method, in `design/system/server-performance.md` § The stories' commits kept between calls (S-0500): two binaries built the same hour, `main` and the story branch, both run from the main checkout with `FLAI_CONFIG` and `FLAI_CACHE_DIR` as `scripts/env.sh` sets them. Time a cold `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'` and `flai stats --json`, three runs each, alternating, with no branch tip moved between runs. Give the median wall time and the `exec.git.log` phase, the load average, and the commit each binary was built from. Check that the two binaries' `claims` in `flai stats --json` are identical.

Add a section for this story after S-0500's, saying what the kept file removes and what still reads the whole log: the first call after a tip moved. Correct S-0500's sentence that a fresh process has nothing kept. Never write a figure you did not observe.

This task waits for T-2267: it measures that change.

## Done when

- The section gives the cold `stats.get` and `flai stats --json` figures of both binaries, measured on this repository, with the `exec.git.log` phase of each.
- The two binaries' `claims` were compared and the section says whether they matched.
- `flai test design/system/server-performance.md` passes.

## Notes
