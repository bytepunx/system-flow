---
id: T-0617
type: task
nature: remediation
title: The dashboard says the clone is missing published tags instead of listing published stories as waiting
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T07:51:03Z
updated: 2026-10-01T08:08:01Z
transitions:
  - to: ready
    at: 2026-10-01T07:59:57Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T07:59:57Z
    by: agent-S-0174
  - to: ready
    at: 2026-10-01T08:00:54Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T08:04:56Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:08:01Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [flaiover/src, flai/internal/workitem]
usage:
  source: log
  seconds: 242
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 76
      output: 25433
      cache_read: 6730227
      cache_write: 74133
      cost: 2.448
---
# T-0617 The dashboard says the clone is missing published tags instead of listing published stories as waiting

## Work

- The Publish banner shows the missing tags, the remote, and the command that fetches them when `publish.preview` reports a component behind, and a warning when the remote could not be asked.
- The done lane hides the archived stories it shows only because they look unpublished while the clone lags, and marks none as waiting.

## Done when

- flaiover unit tests cover the banner in each state and the done lane while lagging; `npm run check`, lint, and tests pass.

## Notes
- 2026-10-01T08:00:54Z: moved to ready: the I-0024 tasks feed the banner; worked first
