---
id: T-2655
type: task
nature: improvement
title: Measure a 20-forecast replan in a scratch copy and record cause 9's before and after
status: done
parent: S-0599
owner: alex
created: 2026-11-07T23:49:25Z
updated: 2026-11-08T01:10:51Z
transitions:
  - to: ready
    at: 2026-11-08T00:44:58Z
    by: system-flow
  - to: in-progress
    at: 2026-11-08T00:45:06Z
    by: system-flow
  - to: done
    at: 2026-11-08T01:10:51Z
    by: agent-S-0599
stream: S-0599
tags: [cli]
touches: [design/system/server-performance.md, design/system/strategic-agents.md]
after: [T-2647]
usage:
  source: log
  seconds: 1545
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 78
      output: 13588
      cache_read: 14241273
      cache_write: 231099
      cost: 4.4662
---
# T-2655 Measure a 20-forecast replan in a scratch copy and record cause 9's before and after

## Work

Copy the main checkout to a scratch folder under `/tmp`, as S-0596 did. In it, time a replan that moves 20 forecasts, first with the replanner from the main branch and then with the story's. Run each three times and take the median. Change nothing in the main checkout.

Record what you measured:

- `design/system/server-performance.md`: a section after § A second pass at 596 stories (S-0596) for cause 9, in the shape of the earlier "After the stories" sections: the cause, the fix, and the before and after measured. Say how loaded the host was.
- `design/system/strategic-agents.md` § The replan policy: the replay writes the moved forecasts through `itemedit.ApplyForecasts` in one pass, with one check over the written items and one index write. Replace the "known cost" paragraph with the measured time.

Waits for T-2647: it measures the replanner that task changes.

## Done when

- The story's replan of 20 forecasts in the scratch copy takes under 5 s, by the median of three runs, and both figures are in `server-performance.md`.
- § The replan policy describes the batch write and no longer says a replan costs about a second a story.
- `flai test design/system/server-performance.md design/system/strategic-agents.md` passes.

## Notes
