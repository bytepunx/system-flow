---
id: T-0886
type: task
nature: feature
title: flai guard lets the orchestrator publish only through release_publish, and only while the publish permission is on
status: done
parent: S-0222
owner: alex
created: 2026-10-05T04:46:25Z
updated: 2026-10-06T11:58:34Z
transitions:
  - to: ready
    at: 2026-10-06T11:53:20Z
    by: agent-S-0222
  - to: in-progress
    at: 2026-10-06T11:53:21Z
    by: agent-S-0222
  - to: done
    at: 2026-10-06T11:58:34Z
    by: agent-S-0222
stream: S-0222
tags: [flai]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 313
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 35
      output: 172
      cache_read: 1162657
      cache_write: 83352
      cost: 0.789
---
# T-0886 flai guard lets the orchestrator publish only through release_publish, and only while the publish permission is on

## Work

In `flai/internal/guard/guard.go`, add `release_publish` to the orchestrator role's rules that S-0218 writes. Pass it while `orchestration.permissions.publish` is on. Refuse it otherwise, naming `publish` as the permission that would allow it, as S-0218's refusals do. Refuse it to every other role and to sub-agents.

Refuse the orchestrator's other routes to the remote whatever its permissions, so that `release_publish` is the only way it publishes:

- `flai release --pending`, `--apply`, and `--dry-run`
- `flai push`
- `git push` and `git tag`

`flai release --evaluate` stays a read, as S-0217's T-0813 makes it.

This task waits for nothing in this story. It needs S-0218's orchestrator role and permission check, which the story's `after` brings in. It runs with the `whole_epics` task, whose paths it does not share.

## Done when

- guard tests pass `release_publish` for the orchestrator with `publish` on
- guard tests refuse it with `publish` off, naming the permission, and refuse it for the planner and a sub-agent
- guard tests refuse the orchestrator's `flai release --pending`, `flai push`, `git push`, and `git tag` with `publish` on, and pass its `flai release --evaluate`
- `go test ./internal/guard/` passes

## Notes
