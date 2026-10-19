---
id: T-1640
type: task
nature: feature
title: The manifest takes merge_queue.batch_max and the host takes the merge-queue host action, off by default
status: done
parent: S-0399
owner: alex
created: 2026-10-16T05:00:47Z
updated: 2026-10-19T19:20:50Z
transitions:
  - to: ready
    at: 2026-10-19T19:00:25Z
    by: agent-S-0399
  - to: in-progress
    at: 2026-10-19T19:00:26Z
    by: agent-S-0399
  - to: done
    at: 2026-10-19T19:20:50Z
    by: agent-S-0399
stream: S-0399
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/config/config.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, docs/operators/settings.md, flai/internal/manifest/settings_test.go]
usage:
  source: log
  seconds: 1224
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 114
      output: 52823
      cache_read: 7976799
      cache_write: 292570
      cost: 4.2216
---
# T-1640 The manifest takes merge_queue.batch_max and the host takes the merge-queue host action, off by default

## Work

- Add `merge_queue.batch_max` to the manifest, a positive integer, 4 when unset. The manifest's checks reject zero, a negative number, and anything that is not a number. `flai manifest set` sets it through `settings.go`.
- Add the host action `merge-queue` to the host's actions in `hostapi/writes.go`'s `Actions` and `config.HostActions`, off by default. `flai serve enable` and `disable` take it, and the settings page lists it.
- Test the key's default, its validation, and turning the action on and off.
- First in its layer: it waits for no task.

## Done when

The key and the action are read, validated, and switched as above, and their tests pass.

## Notes
