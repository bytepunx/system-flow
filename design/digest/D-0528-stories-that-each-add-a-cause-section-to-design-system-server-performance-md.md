---
id: D-0528
title: Stories that each add a cause section to design/system/server-performance.md change only their own section, after the S-0596 Stories table, and the one that lands second keeps both in cause order
kind: decision
sources: [S-0607, E-0033]
topics: [server-side, back-end]
paths: [design/system/server-performance.md]
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

Sources: S-0607, E-0033.
