---
id: D-0528
title: Stories that each add a cause section to design/system/server-performance.md change only their own section, after the S-0596 Stories table, and the one that lands second keeps both in cause order
kind: decision
sources: [S-0607, E-0033, S-0602, S-0599, S-0598, S-0604]
topics: [server-side, back-end]
paths: [design/system/server-performance.md, flai/internal/check/check.go, flai/internal/check/items.go]
date: 2026-11-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0528 Stories that each add a cause section to design/system/server-performance.md change only their own section, after the S-0596 Stories table, and the one that lands second keeps both in cause order

The stories of E-0033 each record their fix as a section of its own in `design/system/server-performance.md`, at the same place, so a trial merge of one at another's base conflicts there and holds a ready story on overlap. The agents settled each such hold with the same split, in their conversations. From S-0607's `## Decisions` (`wip/archive/agents/S-0607.md`):

```text
MS-0378 (S-0606) and MS-0389 (S-0598): the trial merges conflict in `server-performance.md` only because each story adds its own section at the same place. Agreed, settled: whichever lands second keeps both, in cause order.
MS-0363: S-0607 changes only its own cause 17 section of `design/system/server-performance.md`, never S-0599's cause 9 section or the S-0596 tables, and keeps both in cause order if it lands second.
```

The same split is shared on MS-0366, MS-0375, and MS-0391 between S-0599, S-0606, S-0604, and S-0602, and on MS-0367 and MS-0376 for S-0610: each story adds and changes only its own cause section, neither edits the other's or the S-0596 tables, and whichever lands second keeps both sections in cause order. A story asked about such a hold may share the path with this split (ADR-0134).

S-0602, accepted on 2026-11-09, landed under this split with S-0598, S-0599, S-0601, S-0603, S-0604, and S-0606 (`wip/archive/agents/S-0602.md`):

```text
`design/system/server-performance.md` conflicts with every E-0033 story that adds its own section after the S-0596 Stories table. On each conversation (MS-0366, MS-0375, MS-0391, MS-0397, MS-0398) the second to land keeps both sections in cause order.
```

It also left the S-0596 "Where the time goes" table alone, so its task no longer marked its cause's row there as its planner had written.

S-0604, accepted on 2026-11-09, landed under the same split (`wip/archive/agents/S-0604.md`):

```text
MS-0364: S-0604 changes only its own cause 14 section of `design/system/server-performance.md`; S-0599's cause 9 section and the S-0596 tables are left alone (replied settled).
```

S-0599, accepted on 2026-11-09, settled MS-0359 with S-0605 and MS-0355 with S-0598 the same way (`wip/archive/agents/S-0599.md`). It kept a code conflict apart from the design one: T-2632 wrote `check.Items` in a new file, `flai/internal/check/items.go`, rather than in `check.go`, which S-0598 changes too, so the two branches do not conflict there. S-0598, accepted on 2026-11-09, agreed (`wip/archive/agents/S-0598.md`):

```text
S-0599 adds `check.Items` in `check/items.go` for its forecast batch. It is a separate file with a narrower rule set, so there is no conflict. The two could be merged later.
MS-0353 with S-0600: both stories change `checker.markdown()`. We agreed that whichever is accepted second keeps both, my write branch and S-0600's `LintFolder` for the whole run.
```

So `check` now has two narrow runs beside `check.Run`: `check.RunItemWrite` for a checked write (S-0598, ADR-0264) and `check.Items` for the forecast batch (S-0599), which a later story may merge. Where a code conflict could not be avoided, on `flai/internal/workitem/repo.go` with S-0606 (MS-0352), S-0599 proposed that whichever story rebases second puts S-0606's `parsed.stale` line in `saveLocked`, which S-0599 made the one write path of `Save` and `SaveAll`.

Sources: S-0607, E-0033, S-0602, S-0599, S-0598, S-0604.
