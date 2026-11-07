---
id: I-0304
title: flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out
class: efficiency
status: closed
count: 2
cost: 8m
first_reported: 2026-11-07T11:41:27Z
last_reported: 2026-11-07T11:47:42Z
updated: 2026-11-07T22:24:24Z
---

# I-0304 flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out

## Description
flai test on a change to flai's host API runs none of flaiover's server tests that drive the real flai, so a changed gate in hostapi first fails in the close-out

## Instances

### 2026-11-07T11:41:27Z
Story: S-0587.
S-0587 changed the gate of settings.agent in flai/internal/hostapi/writes.go so that it is per project unless everywhere is true. `flai test` on the hostapi and flaiover paths passed. The close-out's flaiover tier then failed at flaiover/src/lib/server/writes.test.ts:578, which runs the built flai and still expected a 403 without --all-projects. Its vitest `related` run picks tests by the TypeScript import graph, so no change under flai/internal/hostapi selects it. It cost one close-out round: 2m29s, plus the fix and a new request.

### 2026-11-07T11:47:42Z
Story: S-0587.
The same story's second close-out stopped at the integration tier on `TestSettingsIndexIsComplete` (flai/cmd/settings_doc_test.go:41): docs/operators/settings.md named `agent.projects.<root>`, not each key `agent.projects.<name>.command`, `.name`, and `.harnesses.<name>.*`. That test skips under `-short`, so `flai test` on docs/operators/settings.md and flai/internal/config ran it nowhere. T-2580's Done when asked for the settings-index test to pass, and its sub-agent reported that it did. A change to the host configuration's shape meets the settings index first in the close-out.

## Remediation

Story S-0592 remediates this issue, created from it at 2026-11-07T15:15:58Z.
Closed 2026-11-07T22:24:24Z: S-0592: a change under flai/internal/hostapi/ or to flai/cmd/hostapi.go now selects the flaiover-flai tier, which runs vitest related on src/lib/server/testing.ts, so writes.test.ts and the other tests that drive the built flai run in flai test. TestSettingsIndexHasEveryKey checks the configuration and manifest key rows under -short, so go-test catches a schema change, and the settings-index tier runs it when docs/operators/settings.md changes. Repository-manifest tests in flai/internal/verify/select_test.go pin both tiers.
