---
id: T-1499
type: task
nature: remediation
title: Each side names its own missing stamp on every proven connection, and the image lines say an unlisted image is dialled under the allowance
status: done
parent: S-0239
owner: alex
created: 2026-10-12T21:56:34Z
updated: 2026-10-12T22:06:48Z
transitions:
  - to: ready
    at: 2026-10-12T21:56:36Z
    by: agent-S-0239
  - to: in-progress
    at: 2026-10-12T21:56:37Z
    by: agent-S-0239
  - to: done
    at: 2026-10-12T22:06:48Z
    by: agent-S-0239
stream: S-0239
tags: []
touches: [flai/internal/channel/channel.go, flai/internal/channel/channel_test.go, flaiover/src/lib/server/release.ts, flaiover/src/lib/server/release.test.ts, flaiover/src/lib/server/agent.test.ts, flai/internal/serve/imagecheck.go, flai/internal/serve/serve.go, flai/internal/serve/serve_test.go, flai/cmd/serve.go, flai/cmd/serve_test.go, flai/cmd/dashboard.go, flai/cmd/dashboard_test.go, design/system/release-signing.md, docs/operators/index.md, design/system/flai-cli.md, docs/users/flai.md, flaiover/src/lib/server/agent.ts]
after: [T-1265, T-1266]
usage:
  source: log
  seconds: 611
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 104
      output: 39653
      cache_read: 5795045
      cache_write: 224503
      cost: 3.2263
---
# T-1499 Each side names its own missing stamp on every proven connection, and the image lines say an unlisted image is dialled under the allowance

## Work

Two gaps T-1266's review found. It waits for T-1265, whose status lines it changes, and T-1266, whose docs describe the gaps.

- Each side records its own missing stamp on every proven connection, whatever its own allowance: a connection proven with an unstamped side was let through by the peer's allowance, so it is never silent. flai's `Client.unsignedSides` gives `flai` whenever its own release is the zero value; the dashboard's `unsignedSides` gives `dashboard` whenever its own stamp is missing. The peer's missing stamp is recorded only when this side's allowance admitted it, as now. A signed pair still records nothing.
- `ImageVerdict` records that an unlisted image was dialled under the allowance, and `flai serve status`'s image line and `flai dashboard status`'s `flai serve:` line say it is dialled as unsigned is allowed, in place of "flai does not dial it" and "refused it".
- `design/system/release-signing.md`, `design/system/flai-cli.md`, and `docs/operators/index.md` say both as built.

## Done when

- A dev flai connected to a dashboard started with `FLAIOVER_ALLOW_UNSIGNED=1`, with flai's own allowance off, shows `unsigned allowed: flai` in flai's statuses and warns; the mirror holds on the dashboard.
- The image lines of the two statuses say an unlisted image is dialled when the allowance is in force.
- The tests pass with `flai test` on the changed paths.

## Notes

Added by the story's agent after T-1266's review.
