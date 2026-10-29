---
id: D-0423
title: "Plan for S-0537: eight tasks in five layers, a reason-primed pack with the conventions cut by reason: confirm the recommendation"
kind: recommendation
sources: [TH-0719, S-0537]
topics: [orchestration, cli, conventions, template, mcp, dashboard, flai, go, code, all]
paths: [flai/internal/context, flai/internal/mcpserver, flai/cmd, design/system/flai-cli.md, design/system/strategic-agents.md, design/adrs, docs, flai/internal/context/reasons.go, flai/internal/context/reasons_test.go, flai/internal/context/reasons_figures.go, flai/internal/context/reasons_figures_test.go, flai/internal/context/strategic.go, flai/internal/context/strategic_test.go, flai/internal/context/context.go, flai/cmd/prime.go, flai/cmd/prime_test.go, flai/internal/mcpserver/server.go, flai/internal/mcpserver/folder.go, flai/internal/mcpserver/prime_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/agent-context.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/internal/serve/orchestrate_reasons.go, flai/internal/serve/orchestrate_reasons_test.go, flai/internal/serve/orchestrate.go, design/system/conventions.md, ".claude/agents/orchestrator.md", flai/internal/check/check.go, flai/internal/check/vocabulary_test.go, design/conventions/tooling.md, flai/internal/check/check_test.go]
date: 2026-10-29
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0755]
---

# D-0423 Plan for S-0537: eight tasks in five layers, a reason-primed pack with the conventions cut by reason: confirm the recommendation

The operator, alex, answered TH-0719 at 2026-10-29T16:27:25Z.

What was asked:

```text
S-0537 is planned: eight tasks in five layers, touches, a 4h forecast delivering 2026-11-01T01:27Z, and a cost of delay of 150.00 USD a week from TH-0718's input.
```

The recommendation it overruled:

```text
**Recommended:** answer criterion 5 as TH-0717 answers S-0536's criterion 6.

- **Move the day's run on reason-primed packs out of S-0537, into the same follow-up story under E-0030.** That story runs after S-0536 and S-0537 are published. It records a day of runs with the context tokens per action before and after.
- **Take the tokens from the session logs and say so.** Add the cost and the activities from `flai stats`. Putting tokens per action into `flai stats` changes `metrics.md` and needs an ADR.
- **Then S-0537's criteria are 1 to 4,** all verifiable on its branch before review. T-2315 moves to the follow-up story.

**Why:** a criterion the story's agent cannot verify before review holds the story in review until someone runs an unreleased build on the host for a day.

**Or:** keep criterion 5 and T-2315 as written, and run the branch's build on the host for a day before review. T-2315 then asks you on a thread before it does.

Source: design/conventions/work-management.md § Rules
```

The operator's answer:

```text
confirm the recommendation
```

The answer applies D-0422, the operator's answer on TH-0717 for S-0536's criterion 6, which the recommendation named: a criterion that needs the host's flai serve to run the story's own change moves to a follow-up story. S-0537's narrative records the outcome: "T-2315 is cancelled. Criterion 5 is removed. The day of runs joins S-0536's follow-up under E-0030."

Source: TH-0719; S-0537 (wip/archive/agents/S-0537.md, Decisions).
