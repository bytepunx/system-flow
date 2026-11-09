---
id: D-0509
title: A formula's instance enters ready with no touches and is held by the story in progress that started it, so that story's agent clears the hold with flai touches on the instance's own story file
kind: decision
sources: [S-0564, I-0282, S-0620]
topics: [orchestration, planning, formulas, claims]
paths: [design/workflows/formulas, flai/cmd/formula.go]
date: 2026-11-07
status: absorbed
absorbed_by: "design/system/strategic-agents.md#formulas-s-0421"
asked_again: []
cited: [S-0620, TH-0933]
---

# D-0509 A formula's instance enters ready with no touches and is held by the story in progress that started it, so that story's agent clears the hold with flai touches on the instance's own story file

`flai formula run --ready` made an instance story that declared no touches. A ready story with no touches is held (no-touches) by every story in progress, the one that started it included, until those move to review. A story that waits on its instance's acceptance therefore deadlocked with it.

S-0564's agent ran `flai formula run analysis-cadence --ready --autocommit`; the instance S-0584 was held by S-0562 and S-0564. It cleared the hold with `flai touches S-0584 --add` on S-0584's own story file, and recorded the fault as I-0282. S-0585, S-0562's publish instance, was held the same way.

S-0564's narrative, `## Decisions`: "S-0584 given a touch of its own. It entered ready with no touches and was held (no-touches) by S-0564 itself, a deadlock. `flai touches S-0584 --add` on its own story file cleared the hold."

S-0620 fixed I-0282 on 2026-11-09 the way this entry advised. Its narrative, `wip/archive/agents/S-0620.md` § Decisions: "`formulaRun` gives the instance story its own story file as its touch, before `--ready` moves it. A claim on its own file overlaps no other story, so the `no-touches` hold (ADR-0046) never applies." The hand step is no longer needed with a flai that has S-0620.

Sources: S-0564 (`wip/archive/agents/S-0564.md`), I-0282, S-0620 (`wip/archive/agents/S-0620.md`).
