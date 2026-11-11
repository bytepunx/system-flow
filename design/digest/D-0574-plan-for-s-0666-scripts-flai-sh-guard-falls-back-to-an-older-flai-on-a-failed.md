---
id: D-0574
title: "Plan for S-0666: scripts/flai.sh guard falls back to an older flai on a failed build and exits 127 with none, five tasks in three layers: Recommended: confirm the plan as written, T-3069 included"
kind: answer
sources: [TH-1045, S-0666]
topics: [harness, cli, flai, go, code, all]
paths: [design/adrs, design/adrs/README.md, scripts/flai.sh, flai/cmd/guard_test.go, ".opencode/plugins/flai-guard.js", flai/tests/integration/opencode_test.go, scripts/README.md, docs/users/flai.md, design/system/agent-adapters.md, docs/operators/index.md, design/system/strategic-agents.md, design/issues/I-0275-this-repository-s-guard-hooks-let-a-call-through-when-scripts-flai-sh-exits.md, design/issues/summary.md, design/issues/I-0229-two-stories-accepted-apart-leave-main-red-s-0535-reworded-to-backlog-s-condition-and-s-0416-s-testadefinedroleisheldtoitsfile-still-expects-the-old-words.md]
date: 2026-11-10
status: absorbed
absorbed_by: "design/adrs/0280-scripts-flai-sh-guard-runs-an-older-bin-flai-or-else-the-installed-flai-when.md#decision"
asked_again: []
cited: []
---

# D-0574 Plan for S-0666: scripts/flai.sh guard falls back to an older flai on a failed build and exits 127 with none, five tasks in three layers: Recommended: confirm the plan as written, T-3069 included

The operator, alex, answered TH-1045 at 2026-11-10T19:24:09Z.

What was asked:

```text
This is the first plan for S-0666. The start prompt said the operator had written on a thread called "asked", but no such thread exists, and S-0666 had no thread, touches, forecast, or tasks before this run. I planned it as new.
```

The answer:

```text
**Recommended:** confirm the plan as written, T-3069 included.

- ADR-0235 decided that a failed build of this tree fails open. The new ADR reverses that for `guard` alone, so a guarded session never passes a call unguarded. Every other command still fails on a failed build, so the agent working on flai still sees it.
- The fallback runs an older `bin/flai` or the installed flai. ADR-0235 already accepts an older flai as the guard where `go` is absent.
- Keep T-3069: this repository's OpenCode plugin passes an exit 127 today. Both guards should judge 126 and 127 the same way.
- The tasks in three layers, the touches, and the 45m forecast fit the story. The input is set: `time_lost_per_cycle: 5m` from TH-1039. The planner works out the value on its next run.

It reverses a recorded decision, so it is yours to confirm.

Source: design/adrs/0235-the-guard-hook-refuses-a-call-when-flai-guard-cannot-start-exit-126-or-127-and.md
```

Source: TH-1045.

S-0666's agent recorded in its narrative that the recommendation reversed part of ADR-0235, so it was the operator's to confirm. Work went on once the operator moved S-0666 to ready with the plan on it, and the new ADR, [ADR-0280](../adrs/0280-scripts-flai-sh-guard-runs-an-older-bin-flai-or-else-the-installed-flai-when.md), was written proposed so the operator still decides at acceptance (`wip/archive/agents/S-0666.md` § Decisions).
