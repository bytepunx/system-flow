---
id: D-0523
title: A story's agent may give a task to a fresh sub-agent with a full brief rather than a fork, when the brief holds all the task needs and a fork would carry the session's large pack
kind: decision
sources: [S-0596, I-0307]
topics: [all, code]
paths: [design/conventions/delegation.md]
date: 2026-11-07
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0523 A story's agent may give a task to a fresh sub-agent with a full brief rather than a fork, when the brief holds all the task needs and a fork would carry the session's large pack

`delegation.md` § Tasks worked by sub-agents says to use a fork, which inherits the conversation, where the harness offers one. S-0596 did not, and said why in its narrative's `## Decisions` (`wip/archive/agents/S-0596.md`):

```text
Task sub-agents were fresh general-purpose agents with full briefs, not forks: the briefs held everything a task needed, and each fork would have carried this session's 300 KB pack.
```

So a story whose pack is large, and whose tasks can be briefed whole (here three read-only measurement tasks, T-2607 to T-2609), weighs the fork's inherited context against its size. The same decision records that `flai guard` refused those sub-agents `flai hostapi` and `git clone` (I-0307); how they measured instead is in `design/system/server-performance.md` § A second pass at 596 stories (S-0596) › How it was measured.

Source: S-0596.
