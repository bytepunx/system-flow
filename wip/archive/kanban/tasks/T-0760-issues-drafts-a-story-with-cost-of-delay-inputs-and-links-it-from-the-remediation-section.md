---
id: T-0760
type: task
nature: improvement
title: issues drafts a story with cost of delay inputs and links it from the Remediation section
status: done
parent: S-0203
owner: alex
created: 2026-10-03T18:07:41Z
updated: 2026-10-03T18:16:35Z
transitions:
  - to: ready
    at: 2026-10-03T18:08:34Z
    by: agent-S-0203
  - to: in-progress
    at: 2026-10-03T18:11:26Z
    by: agent-S-0203
  - to: done
    at: 2026-10-03T18:16:35Z
    by: agent-S-0203
stream: S-0203
tags: []
touches: [flai/internal/issues]
usage:
  source: log
  seconds: 309
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 39
      output: 14808
      cache_read: 1604387
      cache_write: 57768
      cost: 0.9716
---

# T-0760 issues drafts a story with cost of delay inputs and links it from the Remediation section

## Work

In `flai/internal/issues/stories.go`, `ForStory` takes the creation time and the project's planning cycle and returns, besides title, nature, and body, `Draft: true` and a `*workitem.CostOfDelay` (nil when the issue gives nothing):

- With `cost` and `count`, `inputs.time_lost_per_cycle` is cost × count ÷ cycles, where cycles is the time from `first_reported` to now in planning cycles, at least one. Rounded to the minute, or to the second under a minute; written as `normalise` writes durations.
- An issue with a `## Impact` section carries its `revenue_per_week`, `penalty_per_week`, and `time_lost_per_cycle` lines (`- key: value`) over; a time lost there wins over the derivation. A value that does not parse is left out and said so.
- `by: flai`, `at` the creation time.
- The derivation, or the figures carried over, goes in the story body's `## Notes`.

Add `LinkStory(is, story, now)`, which appends a line naming the story to the issue's `## Remediation` section (adding the section when it has none), sets `updated`, and saves. Waits for nothing: the first layer.

## Done when

- `ForStory` and `LinkStory` behave as above, with behavior tests in `stories_test.go` for: draft set, derivation within one cycle and over several, no cost gives no inputs, an Impact section's amounts and time lost carried over, a bad Impact value left out, and the Remediation line written.
- `go test -race -short ./internal/issues/` passes.

## Notes
