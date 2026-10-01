---
id: T-0625
type: task
nature: remediation
title: The pending plan lists every accepted item it could not plan, with the reason
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T08:00:54Z
updated: 2026-10-01T08:03:31Z
transitions:
  - to: ready
    at: 2026-10-01T08:00:58Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T08:00:58Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T08:03:31Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [flai/internal/release, flai/internal/preview, flai/cmd/release.go]
usage:
  source: log
  seconds: 153
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 28
      output: 9314
      cache_read: 2464815
      cache_write: 27149
      cost: 0.8966
---
# T-0625 The pending plan lists every accepted item it could not plan, with the reason

## Work

- `release.Pending` returns, beside the plans, each accepted item since a component's last tag whose own plan could not be computed (removed, no release rule, or no tag saying which of two components it delivers to), once, with the reason, instead of skipping it with `continue`.
- `flai release --pending` (dry run and real), `publish.preview` (`unplanned`), and the board's done lane keep them in view.

## Done when

- A test with a story touching two components and no component tag shows it listed with the reason, and the rest of the batch still planned; tests and lint pass.

## Notes
