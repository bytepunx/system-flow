---
id: T-2789
type: task
nature: remediation
title: workflow.md says an automatic start waits on the hold and story_start names the touches that clear a no-touches hold, and I-0301 is closed
status: done
parent: S-0633
owner: alex
created: 2026-11-09T22:01:51Z
updated: 2026-11-09T23:20:51Z
transitions:
  - to: ready
    at: 2026-11-09T23:19:43Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T23:19:44Z
    by: system-flow
  - to: done
    at: 2026-11-09T23:20:51Z
    by: agent-S-0633
stream: S-0633
tags: [docs, issues]
touches: [design/system/workflow.md, design/issues/I-0301-flai-serve-starts-a-story-s-agent-for-a-ready-story-the-board-holds-no-touches.md, design/issues/summary.md]
after: [T-2786, T-2787]
usage:
  source: log
  seconds: 67
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 21
      output: 4305
      cache_read: 3246702
      cache_write: 142391
      cost: 1.517
---
# T-2789 workflow.md says an automatic start waits on the hold and story_start names the touches that clear a no-touches hold, and I-0301 is closed

## Work

Update the **Held.** and **Skip-ahead.** bullets of `design/system/workflow.md` to describe what T-2786 and T-2787 built. flai serve judges the hold again at the start of every automatic start, and does not start a story it holds then. `story_start`'s refusal of a story held `no-touches` by its own empty claim names `flai touches` as what clears it. Name the path T-2786's Notes found. Then close the issue with `flai issue close I-0301 --reason "<what fixed it>"`, naming the start path and the two changes, and check that `design/issues/summary.md` lists it closed.

Waits for T-2786 and T-2787: it describes what they changed and closes the issue only once both are done.

## Done when

- `design/system/workflow.md` says an automatic start judges the hold at the start, and what `story_start` tells an agent whose story declares no touches
- `flai issue close I-0301` has run with a reason naming what fixed it, and `design/issues/summary.md` shows I-0301 closed
- `flai check --strict` and the markdown lint pass on the changed files

## Notes
