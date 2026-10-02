---
id: T-0694
type: task
nature: research
title: Report the options to the operator and record the decision
status: done
parent: S-0193
owner: arobson
created: 2026-10-02T12:16:03Z
updated: 2026-10-02T12:35:27Z
transitions:
  - to: ready
    at: 2026-10-02T12:24:21Z
    by: claude-fable-5-1
  - to: in-progress
    at: 2026-10-02T12:24:21Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:35:27Z
    by: claude-fable-5-1
stream: S-0193
tags: []
touches: [design/system/release-signing.md, design/adrs]
after: [T-0693]
usage:
  source: log
  seconds: 666
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 512
      output: 208
      cache_read: 2496090
      cache_write: 17985
      cost: 0
---
# T-0694 Report the options to the operator and record the decision

## Work

Open a thread on S-0193 with `thread_open` that leads with the recommendation and lists the options and the questions the operator must answer: the signing tool and key, whether the image is signed by digest list or by cosign, what the connection check does on failure, and how a development build is allowed. Wait for the answer with `wait_for_events`. Record the decision at once: an ADR with `flai adr new` for the technology and the contract between the two components, and a `## Decision` section in `design/system/release-signing.md` that links it. Waits for T-0693 because the thread reports the whole document.

## Done when

- The thread is answered and its answer is recorded in the narrative's `## Decisions`.
- An ADR records the decision and `design/system/release-signing.md` has a `## Decision` section linking it.
- The story's first criterion is checked.

## Notes
