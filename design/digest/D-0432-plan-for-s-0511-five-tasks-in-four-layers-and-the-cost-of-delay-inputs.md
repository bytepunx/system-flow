---
id: D-0432
title: "Plan for S-0511: five tasks in four layers, and the cost of delay inputs: Recommended for the operator, again: keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1)"
kind: answer
sources: [TH-0674, S-0511]
topics: [cli, flai, go, code, all]
paths: [flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/call.go, flai/internal/guard/guard.go, flai/cmd/guard.go, docs/users/flai-reference.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, design/system/flai-cli.md, design/system/agent-adapters.md, design/adrs/README.md, design/adrs, design/issues/I-0210-flai-guard-lets-a-session-flai-serve-started-run-the-flai-serve-sub-commands-that-write-the-host-configuration.md, design/issues/summary.md]
date: 2026-10-30
status: absorbed
absorbed_by: "design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md#decision"
asked_again: []
cited: []
---

# D-0432 Plan for S-0511: five tasks in four layers, and the cost of delay inputs: Recommended for the operator, again: keep `flai config set` and `flai template use` in S-0511's guard rule (assumption 1)

The operator decided that S-0511's guard rule refuses `flai config set` and `flai template use` to every sub-agent and every session `flai serve` started, beside the `flai serve` sub-commands that write the host configuration.

What was asked: the planner of S-0511 assumed on TH-0674 (assumption 1) that the rule covers the two commands, though I-0210 names only `flai serve`. The orchestrator recommended keeping them, since they write the same `~/.flai/config.json` through the same gap (`agentEnv` drops `FLAI_CONFIG`), and a guard that refuses one path to that file and leaves the other open does not close the risk. The alternative was `flai serve` alone and a separate issue.

The answer: on 2026-10-30T16:48:25Z the operator, alex, confirmed the recommendation. S-0511 built it, and ADR-0220 records it: `flai config set` and `flai template use` are refused, while `flai serve project add` and `remove`, which write the project registry, stay out.

Sources: TH-0674, S-0511, ADR-0220.
