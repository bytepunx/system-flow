---
id: T-2645
type: task
nature: improvement
title: "settings.get passes its context to Host.Settings and records its phases: the configuration, the served projects, and the project's own settings"
status: done
parent: S-0605
owner: alex
created: 2026-11-07T23:49:06Z
updated: 2026-11-08T00:48:52Z
transitions:
  - to: ready
    at: 2026-11-08T00:35:40Z
    by: agent-S-0605
  - to: in-progress
    at: 2026-11-08T00:35:45Z
    by: agent-S-0605
  - to: done
    at: 2026-11-08T00:48:51Z
    by: agent-S-0605
stream: S-0605
tags: [cli]
touches: [flai/internal/hostapi/writes.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/cmd/serve_settings.go]
after: [T-2639]
usage:
  source: log
  seconds: 786
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 28
      output: 7163
      cache_read: 3847611
      cache_write: 92821
      cost: 1.5189
---
# T-2645 settings.get passes its context to Host.Settings and records its phases: the configuration, the served projects, and the project's own settings

## Work

- Change `Host.Settings` (`flai/internal/hostapi/writes.go`) from `func(root string) any` to take the request's context, and have the `settings.get` handler (`flai/internal/hostapi/hostapi.go`) pass it, as `agent.status` passes its own to `host.Agent`.
- In `hostSettings` (`flai/cmd/serve_actions.go`), record its steps with `perf.Track`, one phase each for the host configuration, the harnesses and providers, `servedView` (`settings.projects`, say), and the project's own settings from its manifest and MCP state, named as the other methods name theirs (`repo.open`, `manifest.load`).
- Update the callers: `TestSettingsGetSaysWhatMayChange` (`flai/internal/hostapi/writes_test.go`) and the `hostSettings(root)` calls in `flai/cmd/serve_actions_test.go`. Add a test that `settings.get` under a `perf` recorder records the phases.
- Waits for T-2639, which changes `servedView` in `flai/cmd/serve_actions.go` too, so the phase around it wraps the new code.

## Done when

- `flai hostapi --timing settings.get` logs `settings.get`'s phases on stderr.
- The new phase test, `TestSettingsGetSaysWhatMayChange`, and the `hostSettings` tests in `flai/cmd/serve_actions_test.go` pass.

## Notes
