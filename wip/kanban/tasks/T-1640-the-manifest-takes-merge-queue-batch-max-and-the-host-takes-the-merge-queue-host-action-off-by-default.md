---
id: T-1640
type: task
nature: feature
title: The manifest takes merge_queue.batch_max and the host takes the merge-queue host action, off by default
status: backlog
parent: S-0399
owner: alex
created: 2026-10-16T05:00:47Z
updated: 2026-10-16T05:00:47Z
transitions: []
stream: S-0399
tags: [cli]
touches: [flai/internal/manifest/manifest.go, flai/internal/manifest/settings.go, flai/internal/manifest/manifest_test.go, flai/internal/config/config.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go]
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
