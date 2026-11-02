---
id: D-0454
title: "S-0422's after S-0421: keep it and state the reason in its Notes: Recommended: keep `after: [S-0421]` on S-0422, and add agent-S-0551's line under `### Planning` in S-0422's Notes, as written above"
kind: answer
sources: [TH-0787]
topics: [cli, orchestration, conventions, template, planning, flai, go, code, all]
paths: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, template/CHANGELOG.md, flai/cmd/accept.go, flai/cmd/accept_formula_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/role.go, flai/internal/mcpserver/role_test.go, flai/internal/preview/orchestrator.go, flai/internal/preview/accept_test.go, flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, flai/internal/serve/agents.go, flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/serve/role.go, flai/internal/serve/role_test.go, flai/internal/serve/serve.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0454 S-0422's after S-0421: keep it and state the reason in its Notes: Recommended: keep `after: [S-0421]` on S-0422, and add agent-S-0551's line under `### Planning` in S-0422's Notes, as written above

The operator, alex, answered TH-0787 at 2026-11-02T02:21:48Z.

What was asked:

```text
**Recommendation:** keep `after: [S-0421]` and add one line under `### Planning` in S-0422's Notes:
```

The answer:

```text
**Recommended:** keep `after: [S-0421]` on S-0422, and add agent-S-0551's line under `### Planning` in S-0422's Notes, as written above.

**Why:**

- ADR-0188 §6 says `after` names only a story whose accepted change the story needs on main.
- S-0422 works an instance that `flai formula run` makes. It reads the formula format and the `formula` field, and S-0421 adds both. Its Notes already say "the `formula` field is S-0421's". This is a real dependency.
- The added line states the reason outright, as S-0551's criterion 2 asks, instead of leaving it as a passing mention.

**Why you:** S-0422 is finalized, so its words are yours. My role cannot edit it.

Source: design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md § Decision
```

Source: TH-0787.
