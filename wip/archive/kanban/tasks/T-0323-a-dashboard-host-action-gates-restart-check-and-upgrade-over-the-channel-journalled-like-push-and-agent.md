---
id: T-0323
type: task
nature: feature
title: A dashboard host action gates restart, check, and upgrade over the channel, journalled like push and agent
status: done
parent: S-0081
owner: alex
created: 2026-09-22T21:12:43Z
updated: 2026-09-22T21:36:53Z
transitions:
  - to: ready
    at: 2026-09-22T21:22:49Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T21:22:50Z
    by: system-flow
  - to: review
    at: 2026-09-22T21:36:53Z
    by: system-flow
  - to: done
    at: 2026-09-22T21:36:53Z
    by: system-flow
stream: S-0081
tags: []
---
# T-0323 A dashboard host action gates restart, check, and upgrade over the channel, journalled like push and agent

## Work
Add `ActionDashboard = "dashboard"` to `flai/internal/hostapi/writes.go`'s `Actions`, describing what it lets a dashboard do. Add methods: `dashboard.status` (read, wraps `flai dashboard status --json`), `dashboard.check` (read, wraps `flai dashboard check --json`; a read because it changes no container or registration, matching `push.pending`'s precedent of a read that does real work), `dashboard.restart` and `dashboard.upgrade` (`action: ActionDashboard`, `progress: true` so the dashboard sees the pull and swap happen; wrap the new subcommands, `upgrade` passing `--tag` through when given), `dashboard.stop` (`action: ActionDashboard`, wraps the existing `flai dashboard stop --json`, so the host action does exactly what the operator's own shell does).

Update `design/system/dashboard-host-channel.md` (the host actions list) and `design/system/flai-cli.md` (the hostapi method table) to describe the new action and methods. Update `docs/operators/index.md`'s host actions section (`flai serve enable dashboard`) alongside `push` and `agent`.

## Done when
`flai/internal/hostapi/writes_test.go` and `contract_test.go` cover the five new methods (good and refused cases, per their existing per-method table convention); `go test -race -short ./internal/hostapi/...` and `golangci-lint run ./...` clean; docs updated.

## Notes
No ADR: a fourth instance of the host-actions pattern ADR-0029 already established (push, agent, now dashboard), not a new pattern. Depends on T-0322's subcommands existing.
