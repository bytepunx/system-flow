---
id: T-0695
type: task
nature: research
title: Plan the child stories of E-0015 from the decision
status: done
parent: S-0193
owner: arobson
created: 2026-10-02T12:16:03Z
updated: 2026-10-02T12:38:03Z
transitions:
  - to: ready
    at: 2026-10-02T12:35:27Z
    by: claude-fable-5-1
  - to: in-progress
    at: 2026-10-02T12:35:27Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:38:03Z
    by: claude-fable-5-1
stream: S-0193
tags: []
after: [T-0694]
usage:
  source: log
  seconds: 156
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 170
      output: 6555
      cache_read: 686605
      cache_write: 16501
      cost: 0.8311
---
# T-0695 Plan the child stories of E-0015 from the decision

## Work

Create the stories under E-0015 that deliver the decision, in the backlog with `flai story new --epic E-0015`, each with a goal, acceptance criteria as checkboxes, topics, touches, and `--after` where one depends on another; list them in the design document's `## Decision` section. Expected shape, to be confirmed by the decision: the signing key and the CI step that signs flai's checksums; `flai self-upgrade` and `install.sh` verify the signature; the flaiover image is signed in CI and `flai dashboard` verifies before it runs; the release manifest and the connection check on both sides; the operator documentation for keys and the development case. Waits for T-0694 because the stories follow the decision.

## Done when

- Every story the decision needs exists under E-0015 in the backlog with goal and criteria.
- The story's second criterion is checked.

## Notes
