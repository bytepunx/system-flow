---
id: D-0186
title: A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts
kind: cause
sources: [I-0088]
topics: [orchestration, planning, analysis, flai, cli, go, flaiover, dashboard, sveltekit, template, conventions, code, all]
paths: [flai/internal/harness, flai/internal/hostapi, flai/internal/preview, flai/internal/mcpserver, design/adrs, flai/cmd/accept.go, flai/internal/guard, flaiover/src/lib/components/Review.svelte, flaiover/src/routes/items, design/system/workflow.md, design/system/strategic-agents.md, docs/users/flai.md, docs/users/flaiover.md, flai/cmd/accept_orchestrator_test.go, flai/cmd/guard.go, flaiover/src/lib/components/Review.svelte.test.ts, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai-reference.md, docs/operators/settings.md]
date: 2026-10-06
status: absorbed
absorbed_by: "design/system/workflow.md#what-waits-on-an-acceptance-s-0296"
asked_again: []
cited: []
---

# D-0186 A story waits in review for the operator while nothing else can start, so the board stands idle until a person accepts

I-0088 was closed on 2026-10-06.

The cause:

```text
Only the operator accepts a story. While it waits in review it keeps its claim, and when every ready story overlaps it (I-0087) nothing else can start. The board then does no work until a person acts, however long that is. The agents were finished and idle; the wait was for a person.
```

The fix:

```text
S-0221 lets the orchestrator accept a story in review under orchestration.permissions.accept_reviews (ADR-0093), and S-0295 made a story in review hold nothing by overlap (ADR-0096). What still waited on a person, a story naming it in after: and every ready story while review is full, now starts once the orchestrator accepts: TestTheOrchestratorsAcceptanceEndsTheBoardsWaitOnTheOperator in flai/internal/serve/review_wait_test.go shows it, and workflow.md's 'What waits on an acceptance' and the operators' guide say so (S-0296). This project has not turned accept_reviews on yet, and S-0286's rule keeping a story that changes a path Claude Code protects for the operator is still to come.
```

Source: I-0088.
