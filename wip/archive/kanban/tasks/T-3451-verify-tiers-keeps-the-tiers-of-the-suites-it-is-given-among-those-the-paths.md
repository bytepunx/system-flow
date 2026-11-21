---
id: T-3451
type: task
nature: feature
title: verify.Tiers keeps the tiers of the suites it is given among those the paths select, and every one under all
status: done
parent: S-0809
owner: alex
created: 2026-11-16T23:49:15Z
updated: 2026-11-18T08:47:20Z
transitions:
  - to: ready
    at: 2026-11-18T08:36:48Z
    by: agent-S-0809
  - to: in-progress
    at: 2026-11-18T08:36:48Z
    by: agent-S-0809
  - to: done
    at: 2026-11-18T08:47:20Z
    by: agent-S-0809
stream: S-0809
tags: [flai]
touches: [flai/internal/verify/select.go, flai/internal/verify/select_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go]
usage:
  source: log
  seconds: 632
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 18
      output: 4658
      cache_read: 1565305
      cache_write: 54406
      cost: 0.7779
---
# T-3451 verify.Tiers keeps the tiers of the suites it is given among those the paths select, and every one under all

## Work

ADR-0312 §6, criterion 1. Give `verify.TierOptions` (`flai/internal/verify/story.go`) the suites a run is for: a list of suite names, or `all`, with the zero value meaning `all` so that every caller that gives none runs what it runs today. In `Tiers`, after `SelectStory` (`flai/internal/verify/select.go`) has chosen the tiers the paths select, keep only those whose `Suite` is in the list; under `all`, keep every one. Filter after the selection and its `covers` handling, as S-0805's `Suite` filter in `selectTiers` does for `flai test --suite`, so that a tier covered only by a tier the filter leaves out still runs when its own suite is listed, and a `manual` tier still never runs. Reuse S-0805's filter where it fits rather than write a second. Record on the report the suites the run was for, so that a caller can name them. The `Tier.Suite` field is S-0805's, which this story reaches through S-0806. This task waits for nothing.

## Done when

- [ ] A test in `select_test.go` or `story_test.go` shows that given the suites `behavior` and `smoke`, `Tiers` runs only the tiers of those suites that the paths select, and given `all`, or no suites, exactly what it ran before this story.
- [ ] A test shows that a tier covered only by a tier of a suite left out runs, and a `manual` tier never runs, whatever the suites.
- [ ] The report names the suites the run was for, and a test checks it.
- [ ] `flai test` on the changed paths passes.

## Notes
