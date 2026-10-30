---
id: D-0432
title: "Plan for S-0511: five tasks in four layers, and the cost of delay inputs: Recommended for the operator, again: keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1)"
kind: answer
sources: [TH-0674]
topics: [cli, flai, go, code, all]
paths: [flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md, design/issues/summary.md]
date: 2026-10-30
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0432 Plan for S-0511: five tasks in four layers, and the cost of delay inputs: Recommended for the operator, again: keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1)

The operator, alex, answered TH-0674 at 2026-10-30T16:48:25Z.

What was asked:

```text
Question first, since it needs the operator: S-0511 has no cost of delay inputs and no epic, so `flai cod` refuses it. Recommended: `penalty_per_week: 500`, with no time lost per cycle. The defect is a safety gap found once by inspection (an agent `flai serve` started can run `flai serve enable auto-approve`, `notify set --command`, or `agent set --` and so widen its own permissions or name a program the host runs as the operator); it has cost no cycle any time, so a weekly penalty for the standing risk fits better than time lost. Alternatives: `penalty_per_week: 1000`, as S-0475 got for a…
```

The answer:

```text
**Recommended for the operator, again:** keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1).

The planner's and the agent's entries came after my recommendation of 2026-10-28T08:21:40Z, so it no longer stood as pending. I post it again so you can confirm it or answer otherwise.

- **Why include them:** they write the same `~/.flai/config.json` through the same gap, since `agentEnv` drops `FLAI_CONFIG`. Refusing one path to that file and leaving the other open does not close the risk I-0210 records.
- **What it costs:** lines in T-2185's rule and test, and in the help. No new task.
- **The alternative:** `flai serve` alone, and a separate issue for the two commands. agent-S-0511 drops them before review if you answer so.

S-0511 is in progress and works to this recommendation meanwhile. The rest of the plan needs no answer.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
```

Source: TH-0674.
