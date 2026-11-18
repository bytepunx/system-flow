---
id: T-3462
type: task
nature: feature
title: On the baseline workflow a story's close-out runs and its move to review requires what they did before this story
status: done
parent: S-0807
owner: alex
created: 2026-11-16T23:49:59Z
updated: 2026-11-18T11:09:24Z
transitions:
  - to: ready
    at: 2026-11-18T10:48:42Z
    by: agent-S-0807
  - to: in-progress
    at: 2026-11-18T10:48:42Z
    by: agent-S-0807
  - to: done
    at: 2026-11-18T11:09:24Z
    by: agent-S-0807
stream: S-0807
tags: [flai]
touches: [flai/cmd/verify_baseline_test.go]
after: [T-3457, T-3460]
usage:
  source: log
  seconds: 1242
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 48
      output: 19912
      cache_read: 3935230
      cache_write: 181327
      cost: 2.1748
---
# T-3462 On the baseline workflow a story's close-out runs and its move to review requires what they did before this story

## Work

Criterion 4 and ADR-0312 §8. It waits for T-3457 and T-3460 because it tests the rule and the selection together on the baseline.

- In a test repository on the baseline workflow (`verify: all` at `warn` on the move to review), assert that `flai verify <S-nnnn>` on a story in progress selects the same tiers as before this story. Check it for a story that changes `flai/**`, where `when` runs the queue tiers, and for one that changes only markdown.
- Assert that the move to review with no passing verify at the head warns and is not refused. Assert that the warning is gone once the close-out has passed at the head.
- Run the same selection on this repository's manifest and compare it to the tiers `flai verify` ran before this story. Use a fixed tier list in the test rather than the live manifest, so that S-0805's labels do not change the answer.

## Done when

- `flai/cmd/verify_baseline_test.go` passes, and it fails when either selection or the move's severity changes on the baseline.
- `flai test` passes on the changed path.

## Notes
