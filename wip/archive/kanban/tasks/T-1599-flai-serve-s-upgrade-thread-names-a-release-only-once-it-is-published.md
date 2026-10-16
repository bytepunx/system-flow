---
id: T-1599
type: task
nature: feature
title: flai serve's upgrade thread names a release only once it is published
status: done
parent: S-0391
owner: alex
created: 2026-10-15T05:09:36Z
updated: 2026-10-15T05:21:44Z
transitions:
  - to: ready
    at: 2026-10-15T05:09:57Z
    by: agent-S-0391
  - to: in-progress
    at: 2026-10-15T05:09:57Z
    by: agent-S-0391
  - to: done
    at: 2026-10-15T05:21:44Z
    by: agent-S-0391
stream: S-0391
tags: []
touches: [flai/internal/serve/upgrade.go, flai/internal/serve/upgrade_test.go, flai/internal/serve/serve.go, flai/cmd/serve.go, flai/cmd/serve_actions.go, design/system/flai-cli.md, design/adrs/README.md, design/adrs/0161-flai-serve-s-upgrade-thread-names-a-flai-release-only-once-it-is-published-and.md]
usage:
  source: log
  seconds: 707
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 51
      output: 13507
      cache_read: 3326059
      cache_write: 212795
      cost: 2.4578
---

# T-1599 flai serve's upgrade thread names a release only once it is published

## Work

With `auto-upgrade` off, `flai serve`'s ADR-0158 thread names the newest `flai/vX.Y.Z` tag in the project's history and tells the operator to run `flai host upgrade --version X.Y.Z`. A tag comes minutes before its binaries are published, and the host refuses an unpublished release, so the thread can name a command that fails. Give the upgrader a check that a release is published (from `selfupgrade.List`, as `flai host versions` and the MCP `versions` tool read it), and have `ask` open the thread only once the newest release is published; until then it asks again on the next minute's look. A version found published is not checked again. When the check cannot answer, the thread opens as before. Record it as an ADR refining ADR-0158 with `flai adr new` (no commit), and update ADR-0158's description in `design/system/flai-cli.md`. Waits for nothing: it shares no path with T-1600.

## Done when

- `ask` opens no thread for a release the check says is unpublished, and opens it on a later look once it is published.
- A test in `flai/internal/serve/upgrade_test.go` reproduces the unpublished case and fails without the change.
- The ADR is written and indexed; `design/system/flai-cli.md` says it.
- `flai test` on the changed paths passes.

## Notes
