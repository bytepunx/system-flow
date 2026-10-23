---
id: D-0237
title: "Plan for S-0287: a sub-agent writes an ADR through flai, the story's agent commits it: Recommendation: take the planner's remedy and approve the plan as it stands"
kind: answer
sources: [TH-0334]
topics: [cli, conventions, template, flai, go, code, all]
paths: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go, design/system/agent-context.md, design/conventions/delegation.md, template/root/design/conventions/delegation.md, docs/users/flai.md, design/issues/I-0062-flai-guard-lets-a-task-sub-agent-run-flai-adr-new-but-refuses-flai-adr-topics.md, design/issues/summary.md, design/adrs/0127-flai-guard-lets-a-story-s-sub-agent-write-an-adr-with-flai-adr-new-topics-and.md, design/adrs/README.md, template/CHANGELOG.md]
date: 2026-10-08
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0237 Plan for S-0287: a sub-agent writes an ADR through flai, the story's agent commits it: Recommendation: take the planner's remedy and approve the plan as it stands

The operator, alex, answered TH-0334 at 2026-10-08T04:24:15Z.

What was asked:

```text
S-0287 is planned: four tasks in three layers, touches narrowed to nine files, forecast 45m, cost of delay value 30 USD a week.
```

The answer:

```text
Recommendation: take the planner's remedy and approve the plan as it stands.

The remedy is that a story's sub-agent may run `flai adr new`, `topics`, and `accept`, and call MCP `adr_new`, without committing. `--commit`, `--autocommit`, and `commit` stay refused.

- **It keeps ADR-0060's purpose.** ADR-0060 refuses a sub-agent any call that would change a work item, a thread, or the repository's history. An uncommitted ADR file is none of these. Every commit stays with the story's agent.
- **It changes ADR-0060's letter.** ADR-0060 refuses "a flai command other than one that reads", and `flai adr new` writes. So this refines an accepted ADR's guard rule. That is the operator's to confirm, which is why this is a recommendation and not an answer. T-1271 records the refining ADR.
- **The alternative keeps the defect.** Refusing every `flai adr` command leaves the hand-copied `internal/adr` (S-0249) and the guessed ADR number (S-0220) that I-0062 records.

The plan fits the story:

- T-1270 delivers the guard rule and reproduces I-0062.
- T-1271 delivers the ADR.
- T-1272 delivers both `delegation.md` copies and `docs/users/flai.md`.
- T-1273 closes I-0062.

No two tasks of a layer share a path. The 45m forecast and the 30 USD a week value stand.

S-0287 stays a draft until this thread settles.

Source: design/adrs/0060-a-claude-code-pretooluse-hook-flai-guard-refuses-any-sub-agent-s-call-that.md § Decision
```

Source: TH-0334.
