---
id: T-1763
type: task
nature: feature
title: The host's settings methods list and write a role's agent, schedule, and declared permissions through flai manifest set
status: done
parent: S-0420
owner: alex
created: 2026-10-16T05:46:38Z
updated: 2026-11-02T00:47:11Z
transitions:
  - to: ready
    at: 2026-11-02T00:20:51Z
    by: agent-S-0420
  - to: in-progress
    at: 2026-11-02T00:20:53Z
    by: agent-S-0420
  - to: done
    at: 2026-11-02T00:47:11Z
    by: agent-S-0420
stream: S-0420
tags: [cli]
touches: [flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/settings.go, flai/internal/hostapi/settings_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, flai/internal/serve/role_schedule.go, flai/internal/serve/role_schedule_test.go, docs/operators/settings.md, flaiover/src/lib/server/agent.test.ts]
after: [T-1762]
usage:
  source: log
  seconds: 1578
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 92
      output: 34212
      cache_read: 7254056
      cache_write: 167500
      cost: 3.1028
---
# T-1763 The host's settings methods list and write a role's agent, schedule, and declared permissions through flai manifest set

## Work

- Make the settings host methods list, for each defined role, `roles.<name>.agent`, its schedule when its trigger has one, and the permissions its file declares, with what each lets the role do.
- Write them through `flai manifest set` while the `settings` host action is on, and show them read-only otherwise.
- The three built-in roles keep their keys: `planning.*`, `orchestration.*`, and `analysis.*`.
- Add `roles.<name>.schedule` to the manifest's catalog, a cron expression or `daily` as `analysis.schedule` is. It overrides the schedule in the role file's trigger, and `role_schedule.go` reads it. Give it a row in `docs/operators/settings.md`.

Waits for the status task: both change `flai/cmd/serve_actions.go`, which builds `agent.status` and the settings block.

## Done when

`go test ./internal/hostapi/... ./internal/manifest/... ./internal/serve/... ./cmd/...` passes.

## Notes

`flai/internal/hostapi/settings_test.go` may not exist yet. Create it if so, or put the tests in `hostapi_test.go` and widen this task's touches.

- Replanned 2026-10-21 (TH-0580): the orchestrator approved `roles.<name>.schedule` in this task, so the Work names it, the wait on T-1762, and the wider Done when.
