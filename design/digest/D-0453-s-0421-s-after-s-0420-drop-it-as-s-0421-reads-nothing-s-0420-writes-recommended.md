---
id: D-0453
title: "S-0421's after S-0420: drop it, as S-0421 reads nothing S-0420 writes: Recommended: drop S-0420 from S-0421's `after`, as agent-S-0551 proposes"
kind: answer
sources: [TH-0786]
topics: [cli, orchestration, template, conventions, planning, flai, go, code, all]
paths: [design/adrs, design/adrs/README.md, design/conventions/strategic-agents.md, design/workflows/formulas/README.md, design/workflows/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/system/work-hierarchy.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md, flai/cmd/formula.go, flai/cmd/formula_test.go, flai/cmd/root.go, flai/internal/check/check.go, flai/internal/check/formulas_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go, flai/internal/workitem/fields_test.go, flai/internal/workitem/front-matter-fields.txt, flai/internal/workitem/item.go, template/CHANGELOG.md, template/root/design/conventions/strategic-agents.md, template/root/design/workflows/formulas/README.md, template/root/design/workflows/README.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md]
date: 2026-11-02
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0453 S-0421's after S-0420: drop it, as S-0421 reads nothing S-0420 writes: Recommended: drop S-0420 from S-0421's `after`, as agent-S-0551 proposes

The operator, alex, answered TH-0786 at 2026-11-02T02:08:04Z.

What was asked:

```text
**Recommendation:** drop S-0420 from S-0421's `after`:
```

The answer:

````text
**Recommended:** drop S-0420 from S-0421's `after`, as agent-S-0551 proposes:

```bash
flai edit S-0421 --after ""
```

**Why:**

- ADR-0188 §6 says `after` names only a story whose accepted change the story needs on main, and never an order of preference.
- S-0421 needs nothing from S-0420. S-0421 adds formula files, `flai formula run`, and their check. S-0420 builds the dashboard's role page. A formula's role step runs through `role.run`, and S-0417 (done) added that.
- S-0421's Notes name S-0420 only as a story that was replanned the same way. That is a shared history, not a dependency.
- While the `after` stays, S-0421 waits for S-0420's acceptance with no reason of its own.

**Why you:** S-0421 is finalized, so its words are yours. My role cannot edit it.

Source: design/adrs/0188-a-task-takes-no-folder-touch-the-planner-s-folder-touch-on-a-story-needs-the.md § Decision
````

Source: TH-0786.
