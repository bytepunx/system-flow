---
id: T-1824
type: task
nature: feature
title: The curate role is known to the conventions, the strategic pack, the activity documents, and the manifest's digest.agent
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:05:32Z
updated: 2026-10-24T00:00:58Z
transitions:
  - to: ready
    at: 2026-10-23T23:47:48Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-23T23:47:49Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T00:00:57Z
    by: agent-S-0431
stream: S-0431
tags: [cli]
touches: [flai/internal/conventions/conventions.go, flai/internal/conventions/conventions_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/workitem/activity.go, flai/internal/workitem/activity_test.go, flai/internal/workitem/usage_test.go, flai/internal/usage/usage.go, flai/internal/mcpserver/activity.go, flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/check/check.go, flai/internal/check/digest_agent_test.go, flai/internal/mcpserver/activity_test.go]
after: [T-1823]
usage:
  source: log
  seconds: 788
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 94
      output: 518
      cache_read: 4083152
      cache_write: 124634
      cost: 1.4316
---
# T-1824 The curate role is known to the conventions, the strategic pack, the activity documents, and the manifest's digest.agent

## Work

- `flai/internal/conventions/conventions.go`: `curate` joins `StrategicRoles`, so a convention's `roles` may list it.
- `flai/internal/context/strategic.go`: `roleTopic` maps `curate` to its topic, and `flai prime --role curate` builds its pack with no item.
- `flai/internal/workitem/activity.go`: `ActivityCurator` (`curator`) joins `ActivityKinds`, so `activity_log`, the activity document, and `flai stats`' strategic use take it as they take the analyzer.
- `flai/internal/manifest/manifest.go` and `settings.go`: a `digest` block with `agent`, merged over the project's `agent` as `analysis.agent` is (`Manifest.DigestAgent`), its check (`manifest.digest`), and the setting `digest.agent` for `flai manifest set`.

Waits for T-1823, which names the role, its topic, and its settings.

## Done when

Each file's tests cover the role or the setting, `flai prime --role curate --json` gives a pack in a test, and `flai test flai/internal/conventions flai/internal/context flai/internal/workitem flai/internal/manifest` passes.

## Notes

S-0432 adds `digest.schedule` and `digest.compact_after` to the same block.
