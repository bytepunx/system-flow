---
id: D-0433
title: flai's embedded baseline role files under flai/internal/roles/baseline must equal the template's role files byte for byte, so a change to a template role file is copied there in the same story
kind: decision
sources: [S-0419]
topics: [planning, orchestration, analysis, template, go, code]
paths: [flai/internal/roles/baseline/planner.md, flai/internal/roles/baseline/orchestrator.md, flai/internal/roles/baseline/analyzer.md, flai/internal/roles/baseline/curator.md, flai/internal/roles/baseline_test.go, template/root/design/workflows/roles/planner.md, template/root/design/workflows/roles/orchestrator.md, template/root/design/workflows/roles/analyzer.md, template/root/design/workflows/roles/curator.md]
date: 2026-10-31
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0433 flai's embedded baseline role files under flai/internal/roles/baseline must equal the template's role files byte for byte, so a change to a template role file is copied there in the same story

S-0419 (T-1754) found that S-0416 had already embedded the four baseline role files in flai, but held them to the template's only by their front matter, and the orchestrator's body had drifted. It decided: "This task now holds the embedded copies to the template's files whole, not only their front matter. The orchestrator's body had drifted and was copied over."

Since S-0419 the embedded copy matters more: `roles.Builtin` runs a built-in role from the project's file, else from flai's baseline copy, and the claude-code adapter passes that file's body as the role's definition prompt. A project without the file runs the embedded text.

`TestBaselineIsTheTemplatesRoles` in `flai/internal/roles/baseline_test.go` fails when an embedded file differs from `template/root/design/workflows/roles/<name>.md` in front matter or body, and names the fix:

```sh
cp template/root/design/workflows/roles/<name>.md flai/internal/roles/baseline/<name>.md
```

So a story that changes a built-in role's baseline prompt or front matter changes the template's role file and its embedded copy together. Until the generation story (S-0553) lands, it also changes this repository's `design/workflows/roles/<name>.md` and the `.claude/agents/<name>.md` pair, as `design/system/strategic-agents.md` § Role files (S-0414) says.

Source: S-0419 (wip/archive/agents/S-0419.md, Decisions, T-1754).
