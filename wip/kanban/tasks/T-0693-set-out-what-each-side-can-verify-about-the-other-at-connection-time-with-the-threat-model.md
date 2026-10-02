---
id: T-0693
type: task
nature: research
title: Set out what each side can verify about the other at connection time, with the threat model
status: done
parent: S-0193
owner: arobson
created: 2026-10-02T12:16:03Z
updated: 2026-10-02T12:24:21Z
transitions:
  - to: ready
    at: 2026-10-02T12:22:35Z
    by: claude-fable-5-1
  - to: in-progress
    at: 2026-10-02T12:22:35Z
    by: claude-fable-5-1
  - to: done
    at: 2026-10-02T12:24:21Z
    by: claude-fable-5-1
stream: S-0193
tags: []
touches: [design/system/release-signing.md]
after: [T-0692]
usage:
  source: log
  seconds: 106
  estimated: true
  models:
    - model: claude-fable-5-1
      input: 64
      output: 5
      cache_read: 281189
      cache_write: 10330
      cost: 0
---
# T-0693 Set out what each side can verify about the other at connection time, with the threat model

## Work

Add to `design/system/release-signing.md` a `## Verifying the peer` section: what the epic asks (flai closes the connection to an unsigned flaiover, flaiover closes the connection from an unsigned flai), what a running process can and cannot prove about its own build to a peer without hardware attestation, and the designs that reach as far as the facts allow: a signed release manifest embedded at build time and exchanged in `hello`, flai checking the running container's image digest against the signed digest list with `docker inspect` before it dials, and the agent credential as the thing that already authenticates the pair. State plainly the threat each one stops and the one it does not. Cover the development case: this repository runs flai from source and builds the image locally with `flai dashboard --build`, so say how an unsigned build is allowed on purpose and shown. End with a recommendation. Waits for T-0692 because it edits the same file and builds on the signing recommendation.

## Done when

- The document has a `## Verifying the peer` section with the threat model, one subsection per design, what each guarantees and does not, the development case, and a recommendation.

## Notes
