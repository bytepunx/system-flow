---
id: T-3471
type: task
nature: feature
title: The baseline and the template's board give verify, runs, and verify_at as today's behaviour, and stay equal
status: done
parent: S-0806
owner: alex
created: 2026-11-16T23:50:33Z
updated: 2026-11-17T08:48:20Z
transitions:
  - to: ready
    at: 2026-11-17T08:32:22Z
    by: agent-S-0806
  - to: in-progress
    at: 2026-11-17T08:32:22Z
    by: agent-S-0806
  - to: done
    at: 2026-11-17T08:48:20Z
    by: agent-S-0806
stream: S-0806
tags: [flai, template]
touches: [flai/internal/workflow/baseline/board.md, template/root/design/workflows/board.md, flai/internal/workflow/baseline_test.go, template/CHANGELOG.md]
after: [T-3469]
usage:
  source: log
  seconds: 958
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 9430
      cache_read: 2500369
      cache_write: 104971
      cost: 1.3693
---
# T-3471 The baseline and the template's board give verify, runs, and verify_at as today's behaviour, and stay equal

## Work

Write ADR-0312 §8 into `flai/internal/workflow/baseline/board.md` and, byte for byte, into `template/root/design/workflows/board.md`:

- `verify: all` at `warn` on the transition from `in-progress` to `review`;
- `verify: all` on the transition with the `accept` effect;
- `verify_at: {batch: all, ci: all}`, with no `publish`;
- `runs: {kluster: [deploy, checks]}` on the `review` state.

Say in the body of both files, in a sentence, what the keys mean and that nothing reads them to decide a run yet. Assert each value in `baseline_test.go`, and keep `TestBaselineIsTheTemplatesBoard` passing. Add a line under `## Unreleased` in `template/CHANGELOG.md`.

It waits for T-3469, whose parser must take the keys before the baseline that flai embeds can carry them.

## Done when

- `baseline_test.go` asserts the four values and that `verify_at` has no `publish`.
- `TestBaselineIsTheTemplatesBoard` passes, and `flai test` on the changed paths passes.

## Notes
