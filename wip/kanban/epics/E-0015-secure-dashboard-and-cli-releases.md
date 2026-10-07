---
id: E-0015
type: epic
nature: feature
title: Secure Dashboard and CLI Releases
status: in-progress
owner: alex
created: 2026-10-01T11:06:14Z
updated: 2026-10-07T22:13:38Z
transitions:
  - to: ready
    at: 2026-10-01T11:15:15Z
    by: alex
  - to: in-progress
    at: 2026-10-05T00:30:17Z
    by: alex
tags: [dashboard, cli]
topics: [releases]
usage:
  source: sum
  seconds: 3118
  estimated: true
  turns:
    - day: 2026-10-07
      ceremony: 1
      hand_edits: 1
      work: 21
  models:
    - model: claude-fable-5-1
      input: 2082
      output: 80320
      cache_read: 8413158
      cache_write: 202190
      cost: 10.1839
    - model: claude-haiku-4-5-20251001
      input: 266
      output: 9362
      cache_read: 1564698
      cache_write: 78430
      cost: 0.3016
    - model: claude-opus-5-5
      input: 266
      output: 25292
      cache_read: 10833355
      cache_write: 402288
      cost: 5.6431
    - model: claude-sonnet-5-5
      input: 48
      output: 14014
      cache_read: 566139
      cache_write: 128752
      cost: 0.5753
  strategic:
    - kind: orchestrator
      seconds: 823
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 94
          output: 1657
          cache_read: 17636947
          cache_write: 35717
          cost: 4.3563
cost_of_delay:
  inputs:
    penalty_per_week: 25
    by: planner-E-0015
    at: 2026-10-07T22:13:25Z
  value: 25
  by: planner-E-0015
  at: 2026-10-07T22:13:38Z
---
# E-0015 Secure Dashboard and CLI Releases

## Outcome

The dashboard and CLI are both cryptographically signed in CI with a private key and the CLI verifies new releases during self-upgrade with the public key.

Determine how each component can use this same mechanism to ensure they are only talking to a signed version of the other such that flai terminates the connection to an unsigned flaiover and flaiover terminates the connection from an unsigned flai.

## Stories
- S-0192 Add a way to create a sibling story
- S-0193 Explore ways to sign and verify components
- S-0232 The release key signs flai's checksums.txt in CI and both components carry the public key
- S-0233 flai self-upgrade, flai host upgrade, and install.sh verify the release's signature before installing it
- S-0234 The flaiover image's digest list is signed in CI and published on a flaiover GitHub release
- S-0235 A signed release stamp is built into flai and into the flaiover image
- S-0236 flai dashboard resolves the image from the signed digest list and runs it by digest
- S-0237 flai and flaiover exchange their release stamps in hello and refuse an unsigned peer with close code 4403
- S-0238 flai measures the container's image through Docker before it dials and refuses one no signed list names
- S-0239 dashboard.allow_unsigned lets a development build connect, shown on every page and in every status

## Notes
