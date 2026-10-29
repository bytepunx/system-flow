---
id: D-0421
title: A refusal of a hand start in flai serve sits in the exported start function, not its Check function, since the replanner and the schedules call the Check functions too and a refusal there would drop their item
kind: decision
sources: [S-0534, S-0418]
topics: [orchestration, go]
paths: [flai/internal/serve, flai/internal/serve/budget.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [S-0418, T-1748]
---

# D-0421 A refusal of a hand start in flai serve sits in the exported start function, not its Check function, since the replanner and the schedules call the Check functions too and a refusal there would drop their item

S-0534 made `flai serve` refuse a start asked for by hand while a provider's budget is held. Its narrative records where that refusal goes (wip/archive/agents/S-0534.md, Decisions, T-2285):

> hand-start refusals sit in the exported start functions, not the `*Check` functions, since the replanner and schedules call those too and a refusal there would drop their item.

The queued starts wait instead: the replanner's queue keeps its order, and the analysis and curate schedules keep a due run until the next look. `design/system/flai-cli.md` § The host's budget state lists which hand starts are refused, but not this placement. A later story that adds another refusal of hand starts puts it in the same place, so queued and scheduled work is held, not lost. The same story gated `RunRole`, the start of a project-defined role, in the same way.

S-0418 kept to it when it added event and schedule starts for a project's roles: "Queued and scheduled starts are held, not dropped, by the budget. D-0421: the hand start refuses a held budget, while the replanner and the schedules keep their due work. `startRole` leaves that refusal to `RunRole`" (wip/archive/agents/S-0418.md, Decisions).

Sources: S-0534 (T-2285), S-0418.
