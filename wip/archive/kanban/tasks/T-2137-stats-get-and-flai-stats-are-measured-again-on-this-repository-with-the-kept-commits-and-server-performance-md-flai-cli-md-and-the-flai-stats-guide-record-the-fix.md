---
id: T-2137
type: task
nature: improvement
title: stats.get and flai stats are measured again on this repository with the kept commits, and server-performance.md, flai-cli.md, and the flai stats guide record the fix
status: done
parent: S-0500
owner: alex
created: 2026-10-26T04:26:57Z
updated: 2026-10-26T04:42:28Z
transitions:
  - to: ready
    at: 2026-10-26T04:35:22Z
    by: agent-S-0500
  - to: in-progress
    at: 2026-10-26T04:35:22Z
    by: agent-S-0500
  - to: done
    at: 2026-10-26T04:42:28Z
    by: agent-S-0500
stream: S-0500
tags: [flai, docs]
touches: [design/system/server-performance.md, design/system/flai-cli.md, docs/users/flai.md]
after: [T-2136]
usage:
  source: log
  seconds: 426
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 72
      output: 20366
      cache_read: 3417624
      cache_write: 127134
      cost: 1.9424
---
# T-2137 stats.get and flai stats are measured again on this repository with the kept commits, and server-performance.md, flai-cli.md, and the flai stats guide record the fix

## Work

Measure `stats.get` cold and warm and `flai stats --json` on this repository with the parameters the charts page sends, as § A regression in the chart data: the hold replay (S-0471) › Measured on this repository in `design/system/server-performance.md` did: a `main` binary and a `story/S-0500` binary built the same hour, `flai hostapi --timing stats.get '{"since":"30d","type":"story","bucket":"day"}'` three runs each for cold, a throw-away program built with `go build -overlay` that calls the host API's method table six times in one process for warm, with the median of calls 2 to 6 and the phases from the `request answered` event, and the load average at the time. Check that both binaries give the same `claims` section.

Record it in `design/system/server-performance.md` as a new section after the S-0471 one, with the cause, the fix (the kept commits keyed by `git show-ref --heads`, and when the log is read again), the table of figures, and what is left of a warm call. In `design/system/flai-cli.md` § Internal structure, describe the kept commits beside the paragraph on `release.Pending`'s kept history: the key, what is kept, and that every caller of `storygit.Committed` shares it. In `docs/users/flai.md`, in the `flai stats` section's paragraph that says the touches drift reads git, say that a long-running flai reads the log again only when a branch tip moves. Bump each document's `updated`.

This task waits for T-2136, since it measures and describes the fix T-2136 builds.

## Done when

- `design/system/server-performance.md` has a section on I-0184 with the cause, the fix, the cold and warm figures of `stats.get` and `flai stats --json` for `main` and `story/S-0500` measured on this repository, with the date, the load, and the commits of the two binaries, and the `exec.git.log` phase is absent from the story binary's warm calls
- `design/system/flai-cli.md` § Internal structure describes the kept commits, their key, and that every caller of `storygit.Committed` shares them
- `docs/users/flai.md`'s `flai stats` section says when a long-running flai reads git again for the touches drift
- `flai check --strict` and the markdown lint pass on the three documents

## Notes

The measurement follows S-0471's method so the figures compare. The design docs are written in the same story as the behaviour change, as `documentation.md` asks; the metrics are unchanged, so `design/system/metrics.md` and `docs/users/flai-reference.md` stay as they are.
