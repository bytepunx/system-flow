---
id: T-2776
type: task
nature: improvement
title: Have the settings-index tier run the reference tests and select docs/users/flai-reference.md
status: done
parent: S-0634
owner: alex
created: 2026-11-09T21:58:27Z
updated: 2026-11-09T22:22:13Z
transitions:
  - to: ready
    at: 2026-11-09T22:12:14Z
    by: agent-S-0634
  - to: in-progress
    at: 2026-11-09T22:12:15Z
    by: agent-S-0634
  - to: done
    at: 2026-11-09T22:22:13Z
    by: agent-S-0634
stream: S-0634
tags: [flai, testing]
touches: [system-flow.yaml, flai/internal/verify/select_test.go]
after: [T-2773]
usage:
  source: log
  seconds: 28
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 8
      output: 44
      cache_read: 1064835
      cache_write: 4675
      cost: 0.391
---
# T-2776 Have the settings-index tier run the reference tests and select docs/users/flai-reference.md

## Work

`make flai-reference` writes two generated pages: `docs/users/flai-reference.md` and the flag index in `docs/operators/settings.md`. Today a change to either page selects only the markdown tier, or `settings-index` without the flag index test. In `system-flow.yaml`'s `tests`, add `docs/users/flai-reference.md` to the `settings-index` tier's `paths`. Widen its `-run` so that it also runs `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent`, such as `^(TestSettingsIndex|TestReferenceIsCurrent$|TestFlagIndexIsCurrent$)`.

In `flai/internal/verify/select_test.go`, update `TestRepositoryManifestChecksTheSettingsIndexWhenItChanges` to the new argv. Pin that `docs/users/flai-reference.md` selects the tier too, and that under `--all` the integration tier still covers it.

This task waits for T-2773: the tier runs under `-short`, where the two tests skip until T-2773 is done.

## Done when

- `flai test docs/users/flai-reference.md` runs `settings-index`, and the tier's output lists `TestReferenceIsCurrent` and `TestFlagIndexIsCurrent` as passed.
- `flai test docs/operators/settings.md` runs `settings-index` with the same tests, the `TestSettingsIndex` tests among them.
- `go test -count=1 -run '^TestRepositoryManifestChecksTheSettingsIndexWhenItChanges$' ./internal/verify` in `flai/` passes.
- `flai test system-flow.yaml flai/internal/verify/select_test.go` passes.

## Notes
