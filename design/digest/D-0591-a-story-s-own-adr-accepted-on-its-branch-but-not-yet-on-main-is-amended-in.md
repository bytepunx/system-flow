---
id: D-0591
title: A story's own ADR, accepted on its branch but not yet on main, is amended in place when the story's work changes it, rather than refined by a second ADR
kind: decision
sources: [S-0702, S-0767]
topics: [conventions, design, flai]
paths: [design/adrs, design/system/documentation-standard.md]
date: 2026-11-11
status: active
absorbed_by: ""
asked_again: []
cited: [S-0767]
---

# D-0591 A story's own ADR, accepted on its branch but not yet on main, is amended in place when the story's work changes it, rather than refined by a second ADR

An accepted ADR is immutable (`design/system/documentation-standard.md`, ADR-0223), but a story that writes an ADR with status `accepted` on its branch may still find, before acceptance, that the decision needs another clause.

S-0702 met this when an integration failure called for a rule its own ADR-0286 did not state. It amended the ADR in place: "Amended ADR-0286 point 2, which is S-0702's own ADR and not yet on main, so it names the absent-file rule the integration failure called for, rather than adding a second ADR to refine one that has not been accepted onto main" (wip/archive/agents/S-0702.md, Decisions).

S-0767 did the same with ADR-0297 when its rule changed from words in order to a multiset of words, and amended the title too: "ADR-0297 is amended in place, title included, since it is this story's own and not yet on main (D-0591). Its file name keeps the slug flai adr new gave, which the new title still begins with" (wip/archive/agents/S-0767.md, Decisions). The departure from the first plan went into the ADR's Alternatives considered.

So immutability starts once the ADR is on main. Before then it is the story's own draft: amend its text and its title in place, keep its file name, and record the rejected first version under Alternatives considered, rather than add a second ADR.

Sources: S-0702, S-0767.
