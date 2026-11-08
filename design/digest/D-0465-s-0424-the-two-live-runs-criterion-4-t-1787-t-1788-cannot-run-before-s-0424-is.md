---
id: D-0465
title: "S-0424: the two live runs (criterion 4, T-1787, T-1788) cannot run before S-0424 is accepted; move them to a follow-up story: Recommended: take option 1, as you answered for S-0423 on TH-0803…"
kind: answer
sources: [TH-0810, S-0424]
topics: [cli, orchestration, template, conventions, planning, flai, go, flaiover, dashboard, sveltekit, code, all]
paths: [design/workflows/formulas/analysis-cadence.md, design/workflows/formulas/import.md, template/root/design/workflows/formulas/analysis-cadence.md, template/root/design/workflows/formulas/import.md, template/template.yaml, template/CHANGELOG.md, flai/internal/hostapi/import.go, flai/internal/hostapi/import_test.go, flai/cmd/import.go, flai/cmd/import_commit.go, flai/cmd/import_test.go, flai/cmd/import_commit_test.go, flai/cmd/issue.go, flai/cmd/issue_test.go, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/serve/analysis_schedule.go, flai/internal/serve/analysis_schedule_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, design/system/strategic-agents.md, design/system/continuous-improvement.md, design/system/project-manifest.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/new_test.go, flaiover/src/lib/components/ImportPrompt.svelte, flaiover/src/lib/components/ImportPrompt.svelte.test.ts, flaiover/src/routes/api/import/+server.ts, design/adrs/0245-a-formula-s-role-step-that-names-the-built-in-analyzer-starts-it-through-the.md, design/adrs/README.md, flai/cmd/role.go, flai/cmd/role_test.go, flai/internal/serve/agent_events.go, flai/internal/serve/agents.go, flai/internal/serve/analyze.go, flai/internal/serve/analyze_test.go, flai/internal/serve/role.go, design/system/dashboard-host-channel.md, design/workflows/formulas/README.md, template/root/design/workflows/formulas/README.md, design/issues/I-0246-testaheldwaitfailsonatruncatedstory-and-testservestartsadefinedroleonaneventando.md, design/issues/summary.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: [TH-0912]
---

# D-0465 S-0424: the two live runs (criterion 4, T-1787, T-1788) cannot run before S-0424 is accepted; move them to a follow-up story: Recommended: take option 1, as you answered for S-0423 on TH-0803…

The operator, alex, answered TH-0810 at 2026-11-02T18:07:50Z.

What was asked:

```text
Criterion 4 asks for each formula to be run once here, its instance accepted, before S-0424 goes to review. Neither run can start while S-0424 is open:
```

The answer:

```text
**Recommended:** take option 1, as you answered for S-0423 on TH-0803 (D-0461).

- S-0424's agent cancels T-1787 and T-1788 and removes criterion 4 from S-0424.
- It writes a draft story in E-0024 with `after: S-0424`. The draft holds both tasks' Work and criterion 4's words, and its Notes say why it waits.
- S-0424's Notes record the move.
- The epic still gets both live runs before it is accepted.

**Why:** this is the same case as D-0461 and D-0422. The two formula files, `flai import --formula`, and `flai issue list --report` exist only on `story/S-0424`. `flai formula run` reads the main checkout, and the board's Import runs the installed flai. So neither run can start until S-0424 is accepted, and the board's run also has to wait until S-0424 is published.

**Why I recommend rather than answer:** this takes a criterion out of a finalized story and adds a story. That changes scope, so it is your call. Finalizing the draft is also yours.

**Alternative (option 2):** keep criterion 4. I don't recommend it: S-0424 could never meet the definition of done.

Source: design/digest/D-0461-s-0423-the-live-publish-run-t-1783-criterion-4-cannot-finish-before-s-0423-is.md
```

S-0424's narrative records the answer applied: "Criterion 4 is removed and T-1787 and T-1788 are cancelled. Both live runs move to the draft S-0564 in E-0024, `after: S-0424`, which holds their Work and criterion 4's words." Its title still says "each is run once here": a retitle renames the story's files and links, so its Notes record the move instead.

Sources: TH-0810; S-0424 (wip/archive/agents/S-0424.md).
