---
id: T-1568
type: task
nature: improvement
title: flai serve checks whether the host's flai is behind each served project on its reconcile loop, upgrades it with auto-upgrade on, and asks the operator otherwise
status: done
parent: S-0387
owner: alex
created: 2026-10-14T20:54:41Z
updated: 2026-10-15T02:04:23Z
transitions:
  - to: ready
    at: 2026-10-15T01:43:13Z
    by: agent-S-0387
  - to: in-progress
    at: 2026-10-15T01:43:13Z
    by: agent-S-0387
  - to: done
    at: 2026-10-15T02:04:23Z
    by: agent-S-0387
stream: S-0387
tags: [flai, serve, host]
touches: [flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/internal/serve/outdated_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/cmd/serve.go]
after: [T-1567]
usage:
  source: log
  seconds: 1270
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 104
      output: 40943
      cache_read: 5990938
      cache_write: 169399
      cost: 3.0985
---
# T-1568 flai serve checks whether the host's flai is behind each served project on its reconcile loop, upgrades it with auto-upgrade on, and asks the operator otherwise

## Work

Build what T-1567's ADR decides in `flai serve`; it waits for T-1567 because the ADR settles when to upgrade and what to ask.

- Add the host action `auto-upgrade` to `flai/internal/hostapi/writes.go` (its constant, `Actions`, and `shellOnly`), so no dashboard sees it or turns it on.
- In a new `flai/internal/serve/upgrade.go`, called from the reconcile loop in `serve.go`: ask `release.FlaiOutdated` for each served project, at most once a minute. When it answers and the ADR's condition holds, with `auto-upgrade` on, ask the host to upgrade through the host client `serve.Options.Host` already carries (`flai/cmd/serve_actions.go` wires it), once per newer release, and log it; with it off, open the one thread the ADR names, unless one is already open.
- A failed or refused upgrade (no published release yet, a signature refusal) is logged at warn and tried again only for a newer release, never in a loop.
- Tests that reproduce I-0141's second instance: a serve already running when the project's history gains a newer `flai/v*` tag finds it on its loop, upgrades with the action on, and opens one thread with it off. Keep `outdated_test.go`'s start-time warning passing or move it into the new check.

## Done when

- `auto-upgrade` is a host action that only a shell turns on, with a test in `writes_test.go`.
- `upgrade_test.go` shows the loop finds a release tagged after serve started, upgrades once with `auto-upgrade` on, and opens one thread with it off.
- `flai test` passes on the paths this task changes.

## Notes
