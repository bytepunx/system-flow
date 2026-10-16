---
id: T-1763
type: task
nature: feature
title: The host's settings methods list and write a role's agent, schedule, and declared permissions through flai manifest set
status: backlog
parent: S-0420
owner: alex
created: 2026-10-16T05:46:38Z
updated: 2026-10-16T05:46:38Z
transitions: []
stream: S-0420
tags: [cli]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go]
---
# T-1763 The host's settings methods list and write a role's agent, schedule, and declared permissions through flai manifest set

## Work

- Make the settings host methods list, for each defined role, `roles.<name>.agent`, its schedule when its trigger has one, and the permissions its file declares, with what each lets the role do.
- Write them through `flai manifest set` while the `settings` host action is on, and show them read-only otherwise.
- The three built-in roles keep their keys: `planning.*`, `orchestration.*`, and `analysis.*`.

First layer: shares no path with the status task.

## Done when

`go test ./internal/hostapi/... ./internal/manifest/...` passes.

## Notes

`flai/internal/hostapi/settings_test.go` may not exist yet. Create it if so, or put the tests in `hostapi_test.go` and widen this task's touches.
