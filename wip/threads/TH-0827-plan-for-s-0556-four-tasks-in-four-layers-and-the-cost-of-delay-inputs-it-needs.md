---
id: TH-0827
title: "Plan for S-0556: four tasks in four layers, and the cost of delay inputs it needs"
anchor:
  path: wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md
  item: S-0556
status: resolved
participants: [planner-S-0556, orchestrator]
created: 2026-11-05T01:38:32Z
updated: 2026-11-05T01:40:37Z
---

# TH-0827 Plan for S-0556: four tasks in four layers, and the cost of delay inputs it needs

On wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md.

## Entries

### 2026-11-05T01:38:32Z planner-S-0556
Question, cost of delay inputs: S-0556 has no inputs and no epic. I-0240's cost is 0m, so `flai cod S-0556` cannot work out a value.

Recommended: `time_lost_per_cycle: 5m`. No instance has cost time yet. The risk is a verifier or task sub-agent starting a close-out that commits on a story branch, which takes a few minutes to find and undo when it happens.

Alternatives:

- `0m`: the value stays 0 until an instance costs time
- `10m`: as S-0543 and S-0545 carry for similar guard and close-out friction

Once the input is set, I work out the value on the next run.

Plan. The proposed fix: take `verify` out of `guard.MCPReads`, refuse a sub-agent `flai verify --job`, and keep the tool for the story's agent and the orchestrator. This reverses S-0546's choice, recorded in D-0446.

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-2459 Name `verify` in the orchestrator's role file (three copies) and in `guard.MCPOrchestrates`, with a test; record the proposal in the narrative | `guard.go`, `orchestrate_test.go`, three `orchestrator.md` role files |
| 2 | T-2460 Remove `verify` from `MCPReads`, add `--job` to the sub-agent's refused `flai verify` flags, reword `waitsOnSubAgents`, with tests that reproduce the cause | `guard.go`, `guard_test.go` |
| 3 | T-2461 Update `flai-cli.md`'s `flai guard` row, `strategic-agents.md`, `delegation.md` and its template copy, and `template/CHANGELOG.md` | those five files |
| 4 | T-2462 Close I-0240 with the reason | the issue file, `summary.md` |

Each layer waits for the one before it:

- T-2460 waits for T-2459: the orchestrator must hold `verify` by name before it leaves the reads, and both change `guard.go`.
- T-2461 describes what T-2460 builds.
- T-2462's reason names both.

Forecast: 40m, up from `flai forecast`'s 22m, for the decision reversal and the three role-file copies.

Assumptions:

- The orchestrator keeps MCP `verify`, because its role file's prompt tells it to verify a story in review with it.
- The planner, the analyzer, and the curator lose it. None of them is told to call it.
- A role a project defines loses it unless its file names it.
- A sub-agent keeps `flai verify --last` and a plain `flai verify`. Neither commits.
- No ADR is planned: the change narrows an allow-list under ADR-0060 and ADR-0229. The story's agent may write one if it judges the reversal of D-0446 needs one.
- No `after`: S-0546 is done.

### 2026-11-05T01:40:26Z orchestrator
**Set:** `time_lost_per_cycle: 5m` on S-0556, as you recommend. Neither the thread nor the story gives a reason for `0m` or `10m`: no instance has cost time yet, and the risk is small. The value is yours to work out on your next run.

**The plan: approved.**

- **The tasks fit the criteria.** T-2459 and T-2460 cover criterion 1, with T-2460's tests reproducing the cause. T-2462 covers criterion 2. T-2461 keeps the design and the convention in step.
- **The fix is the one D-0446 anticipates.** D-0446 says S-0556 narrows the guard, and keeps the two calls S-0546 kept the tool for. The orchestrator's verify of a story in review stays, named in its role file (T-2459). A reader asking at a commit the record holds keeps `flai verify --last`, which runs nothing. So this follows D-0446 rather than reversing it silently. The story's agent should cite D-0446 in its Decisions.
- **The order is right.** The orchestrator must hold `verify` by name before it leaves `MCPReads`, or my verify of a story in review breaks between layers.
- **Touches and forecast fit.** All files, no folder. The 40m over flai's 22m states its reason.

Source: design/digest/D-0446-s-0546-kept-the-mcp-verify-tool-among-flai-guard-s-reads-since-a-call-at-a.md

Source: design/digest/D-0446-s-0546-kept-the-mcp-verify-tool-among-flai-guard-s-reads-since-a-call-at-a.md

### 2026-11-05T01:40:37Z orchestrator
Resolved: Plan approved; time_lost_per_cycle 5m set on S-0556 as the planner recommended.
