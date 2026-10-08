---
id: T-1400
type: task
nature: feature
title: settings.harness and the Settings page show and set a harness's guard and deny_protected
status: backlog
parent: S-0355
owner: alex
created: 2026-10-08T08:49:20Z
updated: 2026-10-08T08:49:20Z
transitions: []
stream: S-0355
tags: [cli, dashboard]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
after: [T-1398]
---
# T-1400 settings.harness and the Settings page show and set a harness's guard and deny_protected

## Work

- `settings.harness` in `hostapi/settings.go` takes `guard` and `deny_protected` and passes them to `flai serve agent harness`; the settings read returns them per harness.
- `SettingsPanel.svelte`, in each harness's box: a guard choice (default, `none`) and a `deny_protected` checkbox, each with one line saying what it gives up, saved with the harness's program and arguments.
- Tests: `writes_test.go` for the write and the refused value; `SettingsPanel.svelte.test.ts` for the controls and the payload.

It waits for T-1398, whose flags it calls. It runs together with T-1399, which touches other files.

## Done when

- `flai test flai/internal/hostapi/ flaiover/src/lib/components/SettingsPanel.svelte` passes.

## Notes

Layer 2 of S-0355.
