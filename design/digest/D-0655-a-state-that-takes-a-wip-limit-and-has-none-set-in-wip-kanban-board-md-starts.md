---
id: D-0655
title: A state that takes a WIP limit and has none set in wip/kanban/board.md starts with the default of its kind, a queue 5, an active state 2, a waiting state 3, which gives the baseline's ready 5, in-progress 2, review 3
kind: decision
sources: [S-0739]
topics: [flai, go, code, planning]
paths: [flai/internal/workitem/board.go, wip/kanban/board.md, design/system/workflow.md]
date: 2026-11-17
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0655 A state that takes a WIP limit and has none set in wip/kanban/board.md starts with the default of its kind, a queue 5, an active state 2, a waiting state 3, which gives the baseline's ready 5, in-progress 2, review 3

ADR-0301's question 2 keeps the WIP limits' values in `wip/kanban/board.md` and has the workflow file say only which states take one (`limit: true`). S-0739 set the defaults those states start with, by kind and not by name. Its narrative's Decisions (`wip/archive/agents/S-0739.md`):

> WIP limits stay in `wip/kanban/board.md` (ADR-0301 question 2); the defaults are by kind of a limited state (queue 5, active 2, waiting 3), which gives today's `ready: 5, in-progress: 2, review: 3` on the baseline.

The code is `defaultLimitByKind` and `defaultLimits(w)` in `flai/internal/workitem/board.go`. A board starts with these defaults for its workflow, and `DefaultWIPLimits` is the baseline's. So a workflow file that adds a limited `active` state, such as a `testing` state, gives it 2 until `flai board limit` sets another value.

The design says it differently. `design/system/workflow.md`'s Columns and WIP limits table gives the defaults by the baseline's column names. ADR-0301's `limit` field says a state whose limit `wip/kanban/board.md` does not set "has no limit, as today". The code follows the decision above. A story that relies on what an added state's limit is reads `board.go`, not the ADR's words.

Source: S-0739's narrative, `wip/archive/agents/S-0739.md`, § Decisions.
