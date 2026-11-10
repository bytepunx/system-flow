---
id: T-2859
type: task
nature: remediation
title: The crash-later fake vites crash only after they have answered, and I-0331 is closed
status: in-progress
parent: S-0652
owner: alex
created: 2026-11-10T08:15:17Z
updated: 2026-11-10T08:15:29Z
transitions:
  - to: ready
    at: 2026-11-10T08:15:29Z
    by: system-flow
  - to: in-progress
    at: 2026-11-10T08:15:29Z
    by: system-flow
stream: S-0652
tags: []
touches: [flai/internal/serve/components_test.go, flai/internal/components/devserver_test.go, design/issues/I-0331-testadevserverthatexitsonitsownisforgottenwithitserror-fails-the-integration.md, design/issues/summary.md]
---
# T-2859 The crash-later fake vites crash only after they have answered, and I-0331 is closed

## Work

S-0652's close-out stopped at integration on `TestADevServerThatExitsOnItsOwnIsForgottenWithItsError` (I-0331). The `crash-later` fake vite in `flai/internal/serve/components_test.go` exits 300 ms after it starts, whether or not the probe has seen it answer. Under load the crash comes first and `Start` fails. The fake in `flai/internal/components/devserver_test.go` has the same race. Each fake now crashes 300 ms after it first answers. A start delay before the fake runs reproduces the failure without load. It waits for nothing in this story: the layer after T-2841 and T-2843, found at the close-out.

## Done when

- Both `crash-later` fakes crash only after they have answered a request, and the serve test passes with a start delay that failed it before.
- I-0331 is closed with `flai issue close --commit`, saying what fixed it.

## Notes
