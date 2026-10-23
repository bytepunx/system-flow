---
id: T-2063
type: task
nature: feature
title: I-0148 is closed once a close-out here leaves integration and smoke to the merge queue
status: done
parent: S-0466
owner: alex
created: 2026-10-23T03:22:33Z
updated: 2026-10-23T04:45:20Z
transitions:
  - to: ready
    at: 2026-10-23T04:45:12Z
    by: agent-S-0466
  - to: in-progress
    at: 2026-10-23T04:45:13Z
    by: agent-S-0466
  - to: done
    at: 2026-10-23T04:45:20Z
    by: agent-S-0466
stream: S-0466
tags: []
touches: [design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md, design/issues/summary.md]
after: [T-2061, T-2062]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 3
      output: 857
      cache_read: 231617
      cache_write: 9537
      cost: 0.1336
---

# T-2063 I-0148 is closed once a close-out here leaves integration and smoke to the merge queue

## Work

## Done when

## Notes
