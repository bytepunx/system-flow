---
id: D-0534
title: "S-0606 criterion 3: the walk misses 30 ms with only the archive on the slow tier: Recommended: option 1. Put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the…"
kind: answer
sources: [TH-0908, S-0606]
topics: [server-side, back-end, flai, cli, go, code, all]
paths: [flai/internal/watch/watch.go, flai/internal/serve/serve.go, flai/internal/issues/issues.go, flai/internal/adr/adr.go, flai/internal/hostapi/hostapi.go, design/system/server-performance.md, design/system/flai-cli.md]
date: 2026-11-09
status: absorbed
absorbed_by: "design/system/server-performance.md#fix"
asked_again: []
cited: []
---

# D-0534 S-0606 criterion 3: the walk misses 30 ms with only the archive on the slow tier: Recommended: option 1. Put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the…

The operator, alex, answered TH-0908 at 2026-11-09T07:32:32Z by confirming the orchestrator's recommendation, option 1.

What was asked: S-0606's criterion 3 wants a walk of the watched folders under 30 ms. With only `wip/archive` on the watcher's slow tier, a walk still cost about 45 ms of CPU (75 ms on main), because the fast walk statted 2,695 files outside the archive: threads (908), digest entries (526), conversations (370), issues (297), and ADRs (264). With those five folders on the slow tier too, a walk measured 22 to 23 ms of CPU (p90 27 to 30 ms).

The answer:

```text
**Recommended:** option 1. Put `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on the slow tier. Make flai's own issue and ADR writes atomic (write and rename). Keep criterion 3 as written.
```

The reason given: every flai write to those folders still reaches the dashboard at the next fast look, because a file added, removed, or replaced by a rename moves its folder's modification time. Threads, conversations, and digest entries were written by rename already; issues (`flai/internal/issues`) and ADRs (`flai/internal/adr`) were written in place and are now written by rename too. A `git` merge on main replaces files, so an accepted story's ADRs and issues are seen at once. Only an in-place hand edit in those five folders waits up to the slow tier's 10 s.

So flai serve's watcher keeps `wip/archive`, `wip/threads`, `wip/messages`, `design/digest`, `design/issues`, and `design/adrs` on its slow tier, and any flai code that writes a file in one of them writes it by write and rename, never in place, or the dashboard sees the change up to 10 s late. S-0606 applies it in T-2699 (issue and ADR writes by rename), T-2700 (the five folders on the slow tier, with a test that a flai issue and ADR write are seen at the next fast look), and T-2701 (the walk measured again, recorded under cause 16 in `design/system/server-performance.md`, the tiers in `design/system/flai-cli.md`).

Sources: TH-0908, S-0606.
