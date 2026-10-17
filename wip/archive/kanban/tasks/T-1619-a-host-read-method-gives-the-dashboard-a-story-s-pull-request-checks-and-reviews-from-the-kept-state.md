---
id: T-1619
type: task
nature: feature
title: A host read method gives the dashboard a story's pull request, checks, and reviews from the kept state
status: done
parent: S-0394
owner: alex
created: 2026-10-16T04:36:58Z
updated: 2026-10-17T16:28:25Z
transitions:
  - to: ready
    at: 2026-10-17T16:15:38Z
    by: agent-S-0394
  - to: in-progress
    at: 2026-10-17T16:15:39Z
    by: agent-S-0394
  - to: done
    at: 2026-10-17T16:28:25Z
    by: agent-S-0394
stream: S-0394
tags: [cli, dashboard]
touches: [flai/internal/hostapi/reads.go, flai/internal/hostapi/reads_test.go, flaiover/src/lib/server/agent.ts]
after: [T-1615]
usage:
  source: log
  seconds: 766
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 15
      output: 6240
      cache_read: 893336
      cache_write: 26986
      cost: 0.4562
---
# T-1619 A host read method gives the dashboard a story's pull request, checks, and reviews from the kept state

## Work

- In `flai/internal/hostapi/reads.go`, add the read method `forge.status`, taking a story ID. It answers with T-1615's state for the story: the pull request's URL, its checks, and its reviews, or none.
- It needs no host action: it reads only what flai serve kept.
- Add the method to the contract test if the package has one, `contract_test.go`.
- Waits for T-1615: it reads the kept state.

## Done when

`forge.status` answers for a story with state and one without, and the hostapi tests pass.

## Notes
