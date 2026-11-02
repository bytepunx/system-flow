---
id: T-1789
type: task
nature: improvement
title: The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:49:12Z
updated: 2026-11-02T05:31:17Z
transitions: []
stream: S-0424
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/continuous-improvement.md, design/system/project-manifest.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, docs/users/flai.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
after: [T-1785, T-1786]
---
# T-1789 The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas

## Work

- In `design/system/strategic-agents.md` § The analysis schedule and § Formulas, and `continuous-improvement.md` § Cadence, describe the cadence formula and `analysis.formula`.
- In `design/system/flai-cli.md` and `design/system/flaiover-dashboard.md`, describe `flai import --formula`, `flai import test`, `flai issue list --report`, and the board's Import starting the import formula.
- Add `analysis.formula` to `project-manifest.md` and to `docs/operators/settings.md`.
- Add the import formula and the cadence formula to `docs/operators/index.md`.
- In `docs/users/flaiover.md`, say that the board's Import commits the import before its tests and shows its instance as a story, whose steps run the check and the tests.
- In `docs/users/flai.md`, describe `flai import --formula`, `flai import test`, `flai issue list --report`, and `analysis.formula`.
- Name both formulas in `template/CHANGELOG.md` under `## Unreleased`.

Third layer: waits for T-1785 and T-1786, whose behaviour it describes.

## Done when

The markdown lint, the settings-index test, and `flai check --strict` pass.

## Notes

Rewritten by agent-S-0424 on 2026-11-02 for TH-0588's answers.
