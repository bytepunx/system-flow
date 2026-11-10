---
id: D-0562
title: "Plan for S-0680: one guard task, then the docs and the issue's close in one layer: Recommended: keep S-0680's scope as built"
kind: answer
sources: [TH-0991]
topics: [guard, orchestration, flai, cli, go, code, all]
paths: [flai/internal/guard/guard_tables.go, flai/internal/guard/guard.go, flai/internal/guard/hostconfig.go, flai/internal/guard/guard_test.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/orchestrate_test.go, design/system/strategic-agents.md, docs/users/flai.md, design/issues/I-0318-flai-guard-refuses-the-orchestrator-flai-serve-status-a-read-because-its-reads.md, design/issues/summary.md, design/issues/I-0335-flai-stream-sync-reports-a-path-as-undone-by-the-branch-when-main-gains-a.md]
date: 2026-11-10
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0562 Plan for S-0680: one guard task, then the docs and the issue's close in one layer: Recommended: keep S-0680's scope as built

The operator, alex, answered TH-0991 at 2026-11-10T09:07:50Z.

What was asked:

```text
S-0680 is planned as three tasks in two layers.
```

The answer:

```text
**Recommended:** keep S-0680's scope as built.

- The read forms of `flai serve` that ADR-0220 lists stay open to every agent the guard holds to its reads. That covers the orchestrator, the planner, the analyzer, the curator, defined roles, and read-only sub-agents.
- Keep the container harness's `--image`, `--network`, and `--login` among the refused writes. The agent added them in T-2901 (27cc458).

**Why:**

- Your answer on the scope is still owed. agent-S-0680's entry cleared my earlier recommendation, so I am posting it again. That lets it be confirmed, or answered otherwise.
- The scope is unchanged from my 08:58 recommendation. ADR-0220 § Decision names these forms as reads. `askForHost` already tells every one of these agents that the forms are open. Opening them to the orchestrator alone would leave the guard and that message disagreeing for everyone else.
- The harness flags fall under ADR-0220's rule. It refuses every command that writes the host configuration, in any form. `flai serve agent harness` with a name and a setting is in its table of writes. `--image`, `--network`, and `--login` are such settings, added after the ADR was written. Counting them as writes closes a gap under the same decision, and the change sits inside T-2901's touches. It refuses more than before and opens nothing.
- All three tasks are done.

**Alternative:** open the reads to the orchestrator alone, through `design/workflows/roles/orchestrator.md`. agent-S-0680 narrows it before review if you answer that here.

Source: design/adrs/0220-flai-guard-refuses-every-sub-agent-and-every-session-flai-serve-started-the.md § Decision
```

Source: TH-0991.
