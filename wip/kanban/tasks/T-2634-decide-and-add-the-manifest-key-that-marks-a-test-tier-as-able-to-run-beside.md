---
id: T-2634
type: task
nature: feature
title: Decide and add the manifest key that marks a test tier as able to run beside others
status: backlog
parent: S-0610
owner: alex
created: 2026-11-07T23:48:16Z
updated: 2026-11-07T23:48:16Z
transitions: []
stream: S-0610
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/verify/verify.go, flai/internal/verify/manifest.go, flai/internal/verify/manifest_test.go, docs/operators/settings.md, design/system/project-manifest.md, template/root/system-flow.yaml.tmpl, design/adrs/README.md]
---
# T-2634 Decide and add the manifest key that marks a test tier as able to run beside others

## Work

Write the ADR that says how a tier in the manifest's `tests` says it may run beside other tiers, as the story's Notes ask, and add the key. This task waits for no other: the other tasks read the key it adds.

- Choose the key. Two shapes fit the goal: a flag, such as `parallel: true`, where every flagged tier selected in one run starts at once and an unflagged tier waits for those before it in the list; or `after: [<tier>]`, naming the tiers a tier waits for. The second lets `smoke` wait for `template` and `markdown` only, whose verdicts it reads to skip their steps (the smoke task). Record the choice, the one left out, and why in the ADR.
- Say in the ADR what stays as a serial run gives it: tiers are listed and started cheapest first, `covers` still skips a covered tier, the run stops at the first tier in list order that fails, a tier after it that was running is stopped and listed as not reached, and the verdict and `stopped_at` are those of a serial run.
- Add the field to `manifest.TestTier` (`flai/internal/manifest/manifest.go`), refuse what makes no sense in `TestErrors` (an unknown tier name, a cycle, or a `manual` tier, as the key's shape needs), serialise and compare it in `flai/internal/manifest/settings.go`, and copy it in `verify.FromManifest` (`flai/internal/verify/manifest.go`) onto `verify.Tier` (`flai/internal/verify/verify.go`).
- Document it: a `tests[].<key>` row in `docs/operators/settings.md`, which `TestSettingsIndex` requires; the `tests` key in `design/system/project-manifest.md`; the commented example in `template/root/system-flow.yaml.tmpl`. Say whether an older flai ignores the key, as it ignores `queue`, since flai serve's installed flai reads the manifest a close-out runs.
- The new ADR file is named when it is written; `flai task done` widens the task's touches to it.

## Done when

- The ADR is written and listed in `design/adrs/README.md`.
- The manifest loads, validates, and writes the key, with tests for a valid value and each refused one.
- `verify.Tier` carries it, and `flai test` on the changed paths passes, `settings-index` included.

## Notes
