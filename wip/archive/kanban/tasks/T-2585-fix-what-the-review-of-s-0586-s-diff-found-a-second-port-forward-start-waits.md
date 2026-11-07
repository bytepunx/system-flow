---
id: T-2585
type: task
nature: improvement
title: "Fix what the review of S-0586's diff found: a second port-forward start waits until kubectl forwards, a broken sentence on the Klusters page, and a misplaced bullet in the operators guide"
status: done
parent: S-0586
owner: alex
created: 2026-11-07T10:59:16Z
updated: 2026-11-07T10:59:38Z
transitions:
  - to: ready
    at: 2026-11-07T10:59:30Z
    by: agent-S-0586
  - to: in-progress
    at: 2026-11-07T10:59:31Z
    by: agent-S-0586
  - to: done
    at: 2026-11-07T10:59:38Z
    by: agent-S-0586
stream: S-0586
tags: [kluster]
touches: [flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_forward_test.go, flaiover/src/routes/ops/klusters/+page.svelte, docs/operators/index.md]
usage:
  source: log
  seconds: 7
  estimated: true
  priced_by: harness
  models: []
---
# T-2585 Fix what the review of S-0586's diff found: a second port-forward start waits until kubectl forwards, a broken sentence on the Klusters page, and a misplaced bullet in the operators guide

## Work

The verifier's review of S-0586's diff, after the close-out passed at b509acf, found three things to fix. It waits for every other task of S-0586, whose work it reviews.

- `flai/internal/hostapi/kluster_forward.go`: a second start of a port whose first start still waits for kubectl answered at once, with a local port not listening yet. It now waits until kubectl says it forwards, or says the forward ended.
- `flaiover/src/routes/ops/klusters/+page.svelte`: step 2's sentence "naming only the config file, and that config file, and commits both" reads "naming only the config file, writes that config file, and commits both".
- `docs/operators/index.md`: the bullet on `flai kluster up` and the rest goes back with the bullets above the two-step table, not after the Port-forwards paragraph.

## Done when

- `TestASecondStartWaitsUntilTheFirstForwards` passes, and fails without the change.
- `flai test` passes on the four paths.

## Notes
