---
id: T-1832
type: task
nature: feature
title: digest.schedule and digest.compact_after are manifest settings with their check, and the template's manifest documents them
status: done
parent: S-0432
owner: alex
created: 2026-10-16T06:06:47Z
updated: 2026-10-24T13:41:26Z
transitions:
  - to: ready
    at: 2026-10-24T13:36:40Z
    by: agent-S-0432
  - to: in-progress
    at: 2026-10-24T13:36:40Z
    by: agent-S-0432
  - to: done
    at: 2026-10-24T13:41:26Z
    by: agent-S-0432
stream: S-0432
tags: [cli, template]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, flai/internal/check/digest_test.go, template/root/system-flow.yaml.tmpl, docs/operators/settings.md]
after: [T-1831]
usage:
  source: log
  seconds: 282
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 58
      output: 368
      cache_read: 2981730
      cache_write: 128562
      cost: 1.0465
---
# T-1832 digest.schedule and digest.compact_after are manifest settings with their check, and the template's manifest documents them

## Work

- `flai/internal/manifest/manifest.go`: the `digest` block gains `schedule` (a five-field cron expression in UTC, or `daily`, parsed as `analysis.schedule` is) and `compact_after` (a duration in days such as `30d`, default 30 days), with `DigestSchedule` and `DigestCompactAfter`; `manifest.digest` reports either when it does not parse.
- `settings.go`: both settings for `flai manifest set` and the dashboard's settings panel.
- `template/root/system-flow.yaml.tmpl`: both, commented out, with what they do.

Waits for T-1831, which fixes their meaning and default. Runs beside the absorb and compact task.

## Done when

The manifest tests cover a valid and an invalid value of each and the default; `flai test flai/internal/manifest` passes.

## Notes

`digest.agent` came with S-0431 in the same block.
