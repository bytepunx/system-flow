---
id: T-3120
type: task
nature: research
title: Find what wrote PLACEHOLDER into T-3079's file and correct I-0345's account of it
status: done
parent: S-0717
owner: alex
created: 2026-11-11T09:07:17Z
updated: 2026-11-11T22:33:14Z
transitions:
  - to: ready
    at: 2026-11-11T22:32:57Z
    by: system-flow
  - to: in-progress
    at: 2026-11-11T22:32:57Z
    by: system-flow
  - to: done
    at: 2026-11-11T22:33:14Z
    by: agent-S-0717
stream: S-0717
tags: [flai]
touches: [design/issues/I-0345-wait-for-events-fails-with-no-front-matter-while-another-story-s-checked.md]
usage:
  source: log
  seconds: 17
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 536
      cache_read: 206579
      cache_write: 5098
      cost: 0.0887
---
# T-3120 Find what wrote PLACEHOLDER into T-3079's file and correct I-0345's account of it

## Work

I-0345 says that a checked creation "leaves a placeholder in place, as ADR-0272 describes". That is not what the code does. `itemnew.Create` (`flai/internal/itemnew/itemnew.go`) saves the whole item through `Repo.Create` and `Save`, which write it with `atomicfile.WriteFile`, a temp file renamed into place (`flai/internal/workitem/itemlock.go`, `flai/internal/atomicfile/atomicfile.go`). No Go code under `flai/` holds the string `PLACEHOLDER`. I-0344 records the same event and says its writer is not known.

The times are known:

- T-3079's `created` is 2026-11-10T11:03:52Z.
- S-0708 and S-0667 saw the file hold 11 bytes, `PLACEHOLDER`, from about 11:04:00Z to 11:04:25Z.
- agent-S-0665 moved the task to ready at 11:04:25Z.

So the file held a whole item first, then the placeholder, then the whole item again. Find the writer from S-0665's run log, the serve journal (`.flai-cache/serve/serve.log`), `wip/archive/agents/S-0665.md`, and git. The candidates are the agent's own file write tool, a `flai edit` or `item_edit` (a checked edit writes, checks, and restores without the checked-write lock, as ADR-0272's consequences say), or something else.

Write the writer and its evidence into I-0345's instance. Correct the issue's description where it names ADR-0272's checked creation. If the writer is flai code, add a task to this story for the fix when it is five minutes' work. Otherwise open a thread on S-0717 that names the code path. If the logs no longer hold the answer, say so in the issue and stop there.

This task waits for no other: it reads logs and writes only the issue file.

## Done when

- I-0345's instance names what wrote the file and the evidence, or says that the logs no longer show it.
- I-0345's description no longer says that the checked creation of ADR-0272 writes a placeholder, unless the evidence shows it does.
- When the writer is flai code, a task or a thread on S-0717 carries its fix.

## Notes
