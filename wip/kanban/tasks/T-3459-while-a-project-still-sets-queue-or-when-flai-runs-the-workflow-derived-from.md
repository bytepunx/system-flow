---
id: T-3459
type: task
nature: feature
title: While a project still sets queue or when, flai runs the workflow derived from them, and once its file names batch suites it runs the file and ignores the keys
status: backlog
parent: S-0813
owner: alex
created: 2026-11-16T23:49:42Z
updated: 2026-11-16T23:49:42Z
transitions: []
stream: S-0813
tags: [flai]
touches: [flai/internal/workflow/load.go, flai/internal/workflow/load_test.go, flai/internal/verify/manifest.go, flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
after: [T-3452]
---
# T-3459 While a project still sets queue or when, flai runs the workflow derived from them, and once its file names batch suites it runs the file and ignores the keys

## Work

Criterion 3. After S-0807, S-0808, and S-0809, the close-out, the acceptance's gates, and the batch read their suites from the workflow. This task makes the keys reach them only through the workflow. It waits for T-3452's derivation.

- Where flai loads the project's workflow for verify, the gates, and the batch (`flai/internal/workflow/load.go`, or the caller S-0807 to S-0809 settle on), apply T-3452's derivation while the manifest sets `queue` or `when` and the file names no suites for `batch`. Once the file names them, load the file as it stands.
- In `flai/internal/verify/select.go`, retire the separate `queue` and `when` selection, `SelectStoryLeavingQueue` and its `when` match, in favour of the suites the derived or the project's workflow lists. Keep the `queued` step and its words for a tier a step leaves to the batch, as ADR-0174 shows it.
- Keep `flai/internal/verify/manifest.go`'s copy of the two keys only as far as the derivation needs them. In `flai/internal/verify/story.go`, keep `QueueTiersNotRun` reading the derived suites, so the acceptance preview's `queue_tiers` blocker answers as before.
- Tests show both halves of criterion 3: a project with the keys and the baseline file runs what it ran before, and the existing ADR-0174, ADR-0221, and ADR-0233 tests pass unchanged. A project with no keys whose file names suites runs what the file says.

## Done when

- The tests above pass, and `flai/internal/verify/select_test.go` and `flai/internal/verify/story_test.go` pass with their existing cases unchanged.
- `flai test` passes for `flai/internal/workflow` and `flai/internal/verify`, and for the packages that import them.

## Notes
