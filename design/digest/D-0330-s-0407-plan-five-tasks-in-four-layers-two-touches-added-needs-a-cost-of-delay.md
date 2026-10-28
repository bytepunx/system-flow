---
id: D-0330
title: "S-0407 plan: five tasks in four layers, two touches added; needs a cost of delay input (recommend penalty 500 USD a week on E-0023): Recommendation: answer TH-0467 question 1 with `penalty_per_week…"
kind: answer
sources: [TH-0470]
topics: [agents, cli, template, flai, go, conventions, code, all]
paths: [flai/internal/guard/opencode.go, flai/internal/guard/opencode_test.go, flai/internal/guard/guard.go, flai/internal/guard/subagents.go, flai/internal/guard/subagents_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md, flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_test.go, ".opencode/plugins/flai-guard.js", template/root/.opencode/plugins/flai-guard.js, flai/tests/integration/opencode_test.go, design/system/agent-adapters.md, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md]
date: 2026-10-20
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0657]
---

# D-0330 S-0407 plan: five tasks in four layers, two touches added; needs a cost of delay input (recommend penalty 500 USD a week on E-0023): Recommendation: answer TH-0467 question 1 with `penalty_per_week…

The operator, alex, answered TH-0470 at 2026-10-20T21:13:05Z.

What was asked:

```text
The E-0023 planner drafted S-0407's five tasks. I revisited them against the story. They are kept as written, with two touches added. The cost of delay input is still missing.
```

The answer:

```text
**Recommendation:** answer TH-0467 question 1 with `penalty_per_week: 500` on E-0023, as both planners recommend. Every E-0023 story then shares it by duration, and this planner writes S-0407's value on its next run.

## Why I do not set an input on S-0407 myself

- Under `plan_backlog_stories` I may set inputs only on the story. The input recommended here belongs on the epic. Its alternative, 500 on S-0407 alone, would value this one story at the whole epic's figure.
- flai refuses inputs once a story or its epic has any. An input on S-0407 now could keep you from setting E-0023's. The epic-level input is money and yours, and it is already asked on TH-0467.

## The plan: approved

- **The tasks fit.** Five tasks in four layers. T-1691 and T-1692 share no path. The two added touches (`flai/cmd/guard.go` and the reference) are right, because the guard's help changes.
- **The forecast is flai's.** 45m kept against flai's 39m, with the reason stated.
- **Proposal 1: yes.** Add a `startAllowed` test with an opencode story that has roles to T-1692's `## Work`, so criterion 3 is shown in a test. This changes another planner's task, so it is approved here.
- **Proposal 2: leave it.** T-1694 waits for T-1693, so the docs describe tested behaviour.

The thread stays open for the cost of delay input only.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision
```

Source: TH-0470.
