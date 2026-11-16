---
id: D-0650
title: "An ADR number a story's planned touches name is not reserved: flai adr new takes the next free number past every story branch, so a story whose number another took writes the next one and owes no message"
kind: decision
sources: [S-0800, S-0798]
topics: [conventions, design, planning, flai]
paths: [design/adrs, design/adrs/README.md]
date: 2026-11-16
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0650 An ADR number a story's planned touches name is not reserved: flai adr new takes the next free number past every story branch, so a story whose number another took writes the next one and owes no message

A planner predicts a story's touches file by file, and for a story that will write an ADR it guesses the ADR's number from the highest on main. Two stories planned in parallel can name the same number. S-0797, S-0798, and S-0800 all named ADR-0309 in their touches.

S-0800 took ADR-0309 first and decided: "S-0797 and S-0798 planned ADR-0309 in their touches; flai adr new numbers from every story branch, so their ADRs will take the next free number. Their planners' guesses need no message" (wip/archive/agents/S-0800.md, Decisions). S-0798 then wrote "ADR-0310, not 0309: S-0800 took 0309 first" (wip/archive/agents/S-0798.md, Decisions).

So a planned ADR number is a guess, not a claim. Run `flai adr new`, take the number it gives, and let `flai task done` narrow the touches to the file the branch changes. Do not message the other story about it. Why the numbers do not collide is D-0210, absorbed by `design/system/flai-cli.md` § Commands.

Sources: S-0800, S-0798.
