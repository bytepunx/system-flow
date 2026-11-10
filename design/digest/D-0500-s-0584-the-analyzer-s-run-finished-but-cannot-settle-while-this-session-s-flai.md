---
id: D-0500
title: "S-0584: the analyzer's run finished but cannot settle while this session's flai mcp holds its zombie (I-0299): Recommended: answer this thread once flai serve reports an `agent` change for the…"
kind: answer
sources: [TH-0859, S-0632, I-0299]
topics: [all]
paths: [wip/kanban/stories/S-0584-analyze-the-project-and-list-the-issues-the-report-names.md, design/issues/I-0299-a-role-run-the-stdio-mcp-server-starts-through-role-run-is-never-reaped-so-flai.md, design/issues/summary.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: [S-0632]
---

# D-0500 S-0584: the analyzer's run finished but cannot settle while this session's flai mcp holds its zombie (I-0299): Recommended: answer this thread once flai serve reports an `agent` change for the…

The operator, alex, answered TH-0859 at 2026-11-07T08:32:14Z.

What was asked:

```text
Step T-2545 (analyze) has no exit yet. The analyzer did its work. The run that started it cannot be settled.
```

The answer:

```text
**Recommended:** answer this thread once flai serve reports an `agent` change for the analyzer's run on S-0584, ended or failed. Do not answer before then.

**What I saw at 2026-11-07T08:24:52Z:**

- pid 3881796 is no longer in the process table. The zombie was reaped, so flai serve's next orphan check can settle the run.
- No `agent` change on S-0584 had reached my inbox yet. I cannot read flai serve's state myself.

**Why wait.** The formula makes the run's end the analyze step's exit. If the agent restarts before the end is logged, it blocks T-2545 for no reason. Once the change comes, agent-S-0584 closes T-2545 if the run ended with its work done, or blocks it if the run failed. Either way, an answer then is safe.

**If no change comes:** I-0299 is the cause. It needs its own story. Fixing it in S-0584 would widen that story's scope.

Source: design/workflows/formulas/analysis-cadence.md

Source: design/workflows/formulas/analysis-cadence.md
```

Source: TH-0859.

## I-0299 is fixed (S-0632)

S-0632 fixed I-0299, the cause the answer names, so a run's zombie no longer keeps it from settling. From S-0632's Decisions (`wip/archive/agents/S-0632.md`):

```text
`process` now answers whether the process has ended (`/proc/<pid>/stat` state `Z` or `X`). `Alive` and `Owns` count such a process as gone. `Alive` keeps its old answer when the stat cannot be read, so off Linux and on a read error nothing changes.
```

```text
A handed-over launcher waits for its child in a goroutine and records nothing: the serving flai still settles the run in `settleOrphans`.
```

So with a flai that has S-0632, a step whose role run finished settles on flai serve's next orphan check, and its agent need not wait on a zombie being reaped. The answer's rule still holds: answer such a thread only once flai serve reports the run's `agent` change.

Sources: TH-0859; S-0632 (`wip/archive/agents/S-0632.md`, Decisions); I-0299.
