---
id: S-0086
type: story
nature: remediation
title: The host flai connection banner says 'not connected' even when a flai is connected but too old
status: done
parent: E-0003
owner: alex
created: 2026-09-20T23:52:48Z
updated: 2026-09-22T23:52:42Z
transitions:
  - to: ready
    at: 2026-09-20T23:53:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-22T23:51:21Z
    by: system-flow
  - to: review
    at: 2026-09-22T23:51:22Z
    by: system-flow
  - to: done
    at: 2026-09-22T23:52:42Z
    by: alex
tags: [dashboard]
touches: [flaiover/src]
---
# S-0086 The host flai connection banner says 'not connected' even when a flai is connected but too old

## Goal
`hostFlai.usable` is `connected && !error`, and `HostFlaiBanner.svelte` shows the same banner ("No flai on the host is connected, so the project cannot be shown") whether flai never connected at all or a flai is genuinely connected but too old for the dashboard's required methods. Found 2026-09-20: the operator's `flai serve` kept running an old in-memory version after `flai self-upgrade` replaced the binary on disk (self-upgrade cannot restart an already-running process), and the dashboard reported it as if nothing had ever connected, when `/api/agent` itself already carried the true story (`connected: true`, `error: "flai 1.7.0 ... upgrade flai on the host, then run flai serve stop and flai dashboard"`).

## Acceptance criteria
- [x] A connected flai that is missing required methods (an old version) shows its own message: that flai is connected, its version, what it lacks, and the exact commands to fix it (already in `status.error`) — never "no flai on the host is connected"
- [x] The true not-connected case (`connected: false`) and the "never configured" case (`!status.configured`) keep their own distinct messages, unchanged
- [x] The header badge (`HostFlai.svelte`) also tells the three states apart, not just connected/not connected
- [x] A test for each of the three states

## Tasks
- T-0331 HostFlai badge and banner tell connected-but-outdated apart from not-connected, each with its own test

## Notes
Found while diagnosing a report from the operator ("the dashboard is reporting that no flai serve instance has ever connected to it") that turned out to be this: `flai serve` was still running 1.7.0 in memory after the binary on disk was upgraded to 1.9.0; restarting it (`flai serve stop`, `flai serve start`) fixed the actual connection, but the banner's wording would have said the same misleading thing to anyone in this situation, including someone who has never touched flai serve directly and would not know to look at `flai serve status` themselves.
