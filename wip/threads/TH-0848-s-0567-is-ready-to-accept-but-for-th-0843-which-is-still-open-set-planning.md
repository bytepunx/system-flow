---
id: TH-0848
title: "S-0567 is ready to accept but for TH-0843, which is still open: set planning.agent back to the host harness, then resolve TH-0843"
anchor:
  path: wip/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md
  item: S-0567
status: resolved
participants: [orchestrator, alex]
created: 2026-11-07T05:47:34Z
updated: 2026-11-07T06:08:55Z
---

# TH-0848 S-0567 is ready to accept but for TH-0843, which is still open: set planning.agent back to the host harness, then resolve TH-0843

On wip/kanban/stories/S-0567-flai-serve-starts-the-planner-in-a-container-on-harness-container-and-it-plans.md.

## Entries

### 2026-11-07T05:47:34Z orchestrator
**Recommended, alex:** do what TH-0843 still asks, then resolve TH-0843. I then accept S-0567 on my next run, at 8ab0cd02052d.

```sh
flai manifest set planning.agent='{"model":"claude-opus-5-5"}'
flai thread resolve TH-0843 --reason "planner back on the host harness"
```

**The one blocker:** `flai accept S-0567 --by orchestrator --dry-run` refuses because `thread_open` TH-0843 is not resolved. agent-S-0567 opened it, so I cannot resolve it. The step it still waits on is yours: `system-flow.yaml` in the main checkout still has `harness: container` under `planning.agent`, as an uncommitted change.

**Everything else passes:**

- `flai verify S-0567` passed at 8ab0cd02052d, the head of `story/S-0567`.
- The verifier matched all four criteria. Criteria 1 to 3 are met by the Notes' `### Live run` and the 05:09:49Z entry in `wip/agents/planner.md`. Criterion 4 is met by I-0274, I-0275, I-0276, `design/issues/summary.md`, and `design/system/strategic-agents.md`.
- There is no other blocker.

**Alternative:** accept it yourself with `flai accept S-0567 --by alex`. That resolves TH-0843 with the archive.

### 2026-11-07T06:08:55Z alex
Resolved: S-0567 was accepted
