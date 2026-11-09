---
id: T-2773
type: task
nature: improvement
title: Run TestReferenceIsCurrent and TestFlagIndexIsCurrent under -short
status: done
parent: S-0634
owner: alex
created: 2026-11-09T21:58:15Z
updated: 2026-11-09T22:12:01Z
transitions:
  - to: ready
    at: 2026-11-09T22:03:23Z
    by: agent-S-0634
  - to: in-progress
    at: 2026-11-09T22:03:24Z
    by: agent-S-0634
  - to: done
    at: 2026-11-09T22:12:01Z
    by: agent-S-0634
stream: S-0634
tags: [flai, testing]
touches: [flai/cmd/reference_test.go]
usage:
  source: log
  seconds: 517
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 13
      output: 1334
      cache_read: 1705869
      cache_write: 20070
      cost: 0.5285
---
# T-2773 Run TestReferenceIsCurrent and TestFlagIndexIsCurrent under -short

## Work

In `flai/cmd/reference_test.go`, `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` skip under `-short` because they read the monorepo's docs. Together they take 0.04 s. Replace each `testing.Short()` skip with a skip when the page it reads is absent, `../../docs/users/flai-reference.md` or `settingsPage`, as S-0592 did for `TestSettingsIndexHasEveryKey`. Say why in each test's comment and cite I-0308.

The `go-test` tier runs the tests of the changed packages and their importers under `-short`. `flai/cmd` imports `flai/internal/manifest`, so once the two tests run under `-short`, a change to `flai/internal/manifest/settings.go` runs them in `flai test`. So does a change to any command's flags. This task waits for nothing.

## Done when

- `go test -short -count=1 -run '^(TestReferenceIsCurrent|TestFlagIndexIsCurrent)$' -v ./cmd` in `flai/` runs both tests and they pass, with neither skipped.
- Reproduced once by hand and reverted: change a setting's `Default` in `flai/internal/manifest/settings.go` without `make flai-reference`. `flai test flai/internal/manifest/settings.go` then fails on `TestReferenceIsCurrent`, naming `make flai-reference`. Write what you saw in this task's Notes.
- `flai test flai/cmd/reference_test.go` passes.

## Notes
