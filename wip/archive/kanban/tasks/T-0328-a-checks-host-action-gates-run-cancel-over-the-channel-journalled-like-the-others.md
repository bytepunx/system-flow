---
id: T-0328
type: task
nature: feature
title: A checks host action gates run/cancel over the channel, journalled like the others
status: done
parent: S-0082
owner: alex
created: 2026-09-22T22:36:44Z
updated: 2026-09-22T23:01:03Z
transitions:
  - to: ready
    at: 2026-09-22T22:53:56Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T22:53:57Z
    by: system-flow
  - to: review
    at: 2026-09-22T23:01:02Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:01:03Z
    by: system-flow
stream: S-0082
tags: []
---
# T-0328 A checks host action gates run/cancel over the channel, journalled like the others

## Work
Add `ActionChecks = "checks"` to `hostapi.Actions` (`flai/internal/hostapi/writes.go`), describing what it lets a dashboard do. Add methods: `checks.status` (read, wraps `flai checks status <id> --json`), `checks.tail` (read, `progress: true`, wraps `flai checks tail <id> --from <n> --json`, `from` an integer, defaulting to 0), `checks.run` (`action: ActionChecks`, `detachTimeout` set to the configured time limit plus a buffer — read it from config the same way `checks.status` would, or default to a generous ceiling if that is awkward to reach from a `build` closure; wraps `flai checks run <id> --json`), `checks.cancel` (`action: ActionChecks`, wraps `flai checks cancel <id> --json`). All four take `{id}`; `run` and `cancel` also take `request_id` as writes.

Update `design/system/dashboard-host-channel.md` (the host actions list, following S-0081's entry) and `design/system/flai-cli.md` (the hostapi method table) to describe the new action and methods. Update `docs/operators/index.md`'s host actions section (`flai serve enable checks`, what `flai serve checks set` configures) alongside push, agent, and dashboard.

## Done when
`flai/internal/hostapi/writes_test.go` and `contract_test.go` cover the four new methods (good and refused cases, per the existing per-method table convention) and that `checks.run`/`checks.cancel` are refused while the action is off, journalled the same way `push.run` is. `go test -race -short ./internal/hostapi/...` and `golangci-lint run ./...` clean; docs updated.

## Notes
No ADR: a fifth instance of the host-actions pattern ADR-0029 already established, not a new one. Depends on T-0327.
