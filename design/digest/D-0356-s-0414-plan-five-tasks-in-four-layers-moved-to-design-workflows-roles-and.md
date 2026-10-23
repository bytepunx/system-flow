---
id: D-0356
title: "S-0414 plan: five tasks in four layers, moved to design/workflows/roles and generated definitions per TH-0471: Recommended for the operator: take the planner's answers to questions 2 and 3"
kind: answer
sources: [TH-0567]
topics: [cli, planning, orchestration, template, conventions, flai, go, code, all]
paths: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/workflows/roles/analyzer.md, design/workflows/roles/orchestrator.md, design/workflows/roles/planner.md, design/workflows/roles/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/roles_test.go, flai/internal/roles/conditions.go, flai/internal/roles/conditions_test.go, flai/internal/roles/roles.go, flai/internal/roles/roles_test.go, flai/internal/upgrade/upgrade.go, flai/internal/upgrade/upgrade_test.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md, design/README.md, template/root/design/README.md.tmpl, CLAUDE.md, template/root/CLAUDE.md.tmpl]
date: 2026-10-21
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0356 S-0414 plan: five tasks in four layers, moved to design/workflows/roles and generated definitions per TH-0471: Recommended for the operator: take the planner's answers to questions 2 and 3

The operator, alex, answered TH-0567 at 2026-10-21T12:44:16Z.

What was asked:

```text
S-0414 is replanned on your TH-0471 answers. It keeps its five tasks and four layers, and now has a forecast and a cost of delay.
```

The answer:

```text
Recommended for the operator: take the planner's answers to questions 2 and 3.

- **Question 2: generation lands in S-0419.** That story already moves the built-in prompts into role files, and it would generate the definitions and check them against today's `.claude/agents`. flai serve can pass a generated definition through `--agents`, so it needs no write under `.claude/`. A symlink is then needed only for sessions started by hand. This changes S-0419's scope, so it is yours to decide. The alternative is a new story right after S-0414.
- **Question 3: write your TH-0471 figure on the epic.** Run `flai edit E-0024 --time-lost-per-cycle 4h`. `flai cod` then values every E-0024 story from it, in place of the hand-worked 66.76 on S-0414. Setting an epic's input is yours: my permission covers a backlog story's inputs only, and a story's input would stop the epic from taking yours.

Source: design/system/strategic-agents.md § Cost of delay
```

Source: TH-0567.
