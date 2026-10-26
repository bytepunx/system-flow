---
id: T-2128
type: task
nature: improvement
title: The manifest's threads.remind_after and threads.remind_every are declared, defaulted, validated by flai check, settable with flai manifest set, and documented
status: done
parent: S-0498
owner: alex
created: 2026-10-26T04:13:50Z
updated: 2026-10-26T04:32:54Z
transitions:
  - to: ready
    at: 2026-10-26T04:19:50Z
    by: agent-S-0498
  - to: in-progress
    at: 2026-10-26T04:19:51Z
    by: agent-S-0498
  - to: done
    at: 2026-10-26T04:32:54Z
    by: agent-S-0498
stream: S-0498
tags: [flai, template]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/check/check.go, docs/operators/settings.md, design/system/project-manifest.md, template/root/system-flow.yaml.tmpl, template/CHANGELOG.md, docs/users/flai-reference.md, template/template.yaml]
after: [T-2127]
usage:
  source: log
  seconds: 783
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 62
      output: 20904
      cache_read: 3572528
      cache_write: 106586
      cost: 1.7471
---
# T-2128 The manifest's threads.remind_after and threads.remind_every are declared, defaulted, validated by flai check, settable with flai manifest set, and documented

## Work

Add a `threads` block to the manifest, as the ADR of T-2127 decides, following how `digest.compact_after` is declared:

- `flai/internal/manifest/manifest.go`: a `Threads` struct on `Manifest` (`yaml:"threads,omitempty"`), with `RemindAfter` (a duration, unset by default: no reminders) and `RemindEvery` (a duration, default `24h`); `RemindAfterDuration()` and `RemindEveryDuration()` parsers as `CompactAfterDuration` is; `Threads.Errors()` for a value that is not a duration, a negative one, or a `remind_every` of zero, called from `Manifest.Errors`.
- `flai/internal/manifest/settings.go`: two `Setting` entries of `KindDuration`, `threads.remind_after` and `threads.remind_every`, with their defaults and one-line descriptions, so `flai manifest set` and the host API's `settings.manifest` take them.
- `flai/internal/check/check.go`: report `Threads.Errors()` as `manifest.threads`, as `manifest.digest` is reported.
- Tests in `manifest_test.go` and `settings_test.go`: the defaults, a set value, each refused value, and that `flai manifest set threads.remind_after 4h` round-trips.
- `docs/operators/settings.md`: one row per key under `## Project manifest`; `TestSettingsIndexIsComplete` fails without them.
- `design/system/project-manifest.md`: the two keys in the example and a paragraph saying what they do, beside `dashboard.notify_url`.
- `template/root/system-flow.yaml.tmpl`: the two keys commented out with their defaults, as `digest.compact_after` is; `template/CHANGELOG.md`: an entry.

This task waits for T-2127 so the names and defaults are the ADR's.

## Done when

- `flai manifest set threads.remind_after 4h` and `threads.remind_every 12h` write the keys, `flai check` reports `manifest.threads` on `remind_after: soon`, on a negative value, and on `remind_every: 0`, and a manifest without the block parses with `RemindAfter` unset and `RemindEvery` 24h.
- `flai test` on the changed paths passes, `TestSettingsIndexIsComplete` among them.
- The template's manifest and the design document name both keys.

## Notes

Layer 2, with T-2129 (the blocking of an asked story), which shares no path with it. T-2130 (the host channel) waits for this task because both add rows to `docs/operators/settings.md`.
