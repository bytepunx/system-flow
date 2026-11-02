---
id: D-0461
title: "S-0423: the live publish run (T-1783, criterion 4) cannot finish before S-0423 is accepted; move it to a follow-up story: Recommended: take option 1"
kind: answer
sources: [TH-0803]
topics: [cli, orchestration, template, conventions, planning, flai, go, code, all]
paths: [design/adrs, design/adrs/README.md, design/workflows/formulas/publish.md, template/root/design/workflows/formulas/publish.md, template/template.yaml, template/CHANGELOG.md, design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/work-management.md, template/root/design/conventions/work-management.md, design/system/strategic-agents.md, design/system/flai-cli.md, docs/operators/index.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/release.go, flai/cmd/release_evaluate_test.go, flai/cmd/release_push_action_test.go, flai/cmd/orchestrate_permissions_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/mcpserver/orchestrate.go, flai/internal/mcpserver/orchestrate_test.go]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0461 S-0423: the live publish run (T-1783, criterion 4) cannot finish before S-0423 is accepted; move it to a follow-up story: Recommended: take option 1

The operator, alex, answered TH-0803 at 2026-11-02T04:13:28Z.

What was asked:

```text
T-1783 and criterion 4 ask for the `publish` formula to be run once here, with its instance accepted, before S-0423 goes to review. That can't happen while S-0423 is still open.
```

The answer:

```text
**Recommended:** take option 1.

- S-0423's agent cancels T-1783 and removes criterion 4 from S-0423.
- It writes a draft story in E-0024 with `after: S-0423`. The draft holds T-1783's Work and criterion 4's words, and its Notes say why it waits on S-0423.
- S-0423's Notes record the move.
- The epic still gets its one live run before it is accepted, as S-0423's goal asks.

**Why:** you have answered this before (D-0422, TH-0717, confirmed 2026-10-29). Some criteria need a change that exists only on the story branch to run on the host or the main checkout. Those criteria cannot be verified before review. Such a run moves to a follow-up story after the story. S-0536's day-long run moved this way, and so did S-0537's criterion 5 (D-0423). Here, `flai release --evaluate --require-met` and the `push` gate in `flai release --pending` exist only on `story/S-0423`. The instance's `in: main` steps cannot run them until S-0423 is accepted. If those steps run the installed flai rather than `scripts/flai.sh`, they also have to wait until S-0423 is published.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and adds a story. That changes scope, so it is your call. Finalizing the draft is yours too.

**Alternative (option 2):** keep T-1783. I don't recommend it. S-0423 could never meet the definition of done, so it would stay in progress.

Source: design/digest/D-0422-plan-for-s-0536-eight-tasks-in-five-layers-and-two-questions-on-criterion-6.md
```

Source: TH-0803.
