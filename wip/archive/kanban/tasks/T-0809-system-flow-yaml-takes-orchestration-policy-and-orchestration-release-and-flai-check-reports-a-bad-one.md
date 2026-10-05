---
id: T-0809
type: task
nature: feature
title: system-flow.yaml takes orchestration.policy and orchestration.release, and flai check reports a bad one
status: done
parent: S-0217
owner: alex
created: 2026-10-04T19:08:50Z
updated: 2026-10-05T06:06:11Z
transitions:
  - to: ready
    at: 2026-10-05T05:59:51Z
    by: agent-S-0217
  - to: in-progress
    at: 2026-10-05T05:59:51Z
    by: agent-S-0217
  - to: done
    at: 2026-10-05T06:06:11Z
    by: agent-S-0217
stream: S-0217
tags: [flai]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/orchestration_test.go, design/system/project-manifest.md, docs/operators/settings.md]
usage:
  source: log
  seconds: 380
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 54
      output: 300
      cache_read: 1484825
      cache_write: 91127
      cost: 1.3578
    - model: claude-sonnet-5
      input: 6
      output: 10
      cache_read: 47259
      cache_write: 37480
      cost: 0.0352
---
# T-0809 system-flow.yaml takes orchestration.policy and orchestration.release, and flai check reports a bad one

## Work

Add an `Orchestration` block to `manifest.Manifest` beside `Planning` (`flai/internal/manifest/manifest.go`, lines 116-175 hold `Planning` and its `Errors()`):

- `orchestration.policy`: `cod`, `wsjf`, `throughput`, or `fifo`, the same names as `flai order --by`. Unset means `fifo`, which leaves the operator's order alone.
- `orchestration.release.policy`: `judgement`, `threshold`, or `theme`. Unset means `judgement`, which is never met by itself.
- For `threshold`: `value`, the unreleased cost of delay per week in the project's currency, and `count`, the number of accepted stories not yet released. Either or both may be set.
- For `theme`: `epic`, an epic ID, or `tag`, a tag. One of them.

Give `Orchestration` an `Errors()` like `Planning.Errors()`. It names the key and says how to fix it: a policy outside its list, a threshold with neither figure or a negative one, or a theme with neither or both of epic and tag. `flai check` reports each error as `manifest.orchestration` (`flai/internal/check/check.go`, next to `manifest.planning`).

Leave `orchestration.permissions` out. It is S-0218's. Describe the keys in `design/system/project-manifest.md` and in `docs/operators/settings.md` under `## Project manifest`.

This task waits for nothing. It runs with T-0810, whose paths it does not share.

## Done when

- the manifest parses each policy and release block, and `Errors()` names each bad value in a test
- `flai check` reports a bad block as `manifest.orchestration`, and a fixture test pins it
- `project-manifest.md` and `settings.md` describe both keys, their values, and their defaults
- `go test ./internal/manifest/ ./internal/check/` passes

## Notes
