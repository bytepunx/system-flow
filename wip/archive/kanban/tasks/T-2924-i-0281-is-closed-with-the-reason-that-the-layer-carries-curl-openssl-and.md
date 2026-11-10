---
id: T-2924
type: task
nature: remediation
title: I-0281 is closed with the reason that the layer carries curl, openssl, and GoReleaser and scripts/agent-image.sh --check passes
status: done
parent: S-0672
owner: alex
created: 2026-11-10T08:48:30Z
updated: 2026-11-10T10:10:28Z
transitions:
  - to: ready
    at: 2026-11-10T10:10:15Z
    by: agent-S-0672
  - to: in-progress
    at: 2026-11-10T10:10:16Z
    by: agent-S-0672
  - to: done
    at: 2026-11-10T10:10:28Z
    by: agent-S-0672
stream: S-0672
tags: [docker]
touches: [design/issues/I-0281-the-system-flow-image-layer-lacks-curl-and-goreleaser-so-the-smoke-tier-s.md, design/issues/summary.md]
after: [T-2919]
usage:
  source: log
  seconds: 12
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 5
      output: 1003
      cache_read: 665554
      cache_write: 9545
      cost: 0.2295
---
# T-2924 I-0281 is closed with the reason that the layer carries curl, openssl, and GoReleaser and scripts/agent-image.sh --check passes

## Work

Close I-0281 with `flai issue close I-0281 --reason "<reason>"`. The reason names what fixed it: the layer's curl, openssl, and GoReleaser from T-2917, and `scripts/agent-image.sh --check` from T-2919, which ran the smoke tier inside the layer. It names that check's passing run and the commit it ran at. flai rewrites `design/issues/summary.md`.

It waits for T-2919, whose passing check is what the reason cites.

## Done when

- `flai issue show I-0281` reads closed, with that reason.
- `design/issues/summary.md` no longer lists I-0281 among the open issues.

## Notes
