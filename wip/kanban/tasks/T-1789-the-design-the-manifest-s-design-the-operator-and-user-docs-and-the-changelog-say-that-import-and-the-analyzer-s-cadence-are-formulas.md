---
id: T-1789
type: task
nature: improvement
title: The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas
status: backlog
parent: S-0424
owner: alex
created: 2026-10-16T05:49:12Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0424
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/continuous-improvement.md, design/system/project-manifest.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flaiover.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
after: [T-1785, T-1786]
---
# T-1789 The design, the manifest's design, the operator and user docs, and the changelog say that import and the analyzer's cadence are formulas

## Work

- In `design/system/strategic-agents.md` § The analysis schedule and `continuous-improvement.md` § Cadence, describe the cadence formula and `analysis.formula`.
- Add `analysis.formula` to `project-manifest.md` and to `docs/operators/settings.md`.
- Add the import formula and the cadence formula to `docs/operators/index.md`.
- In `docs/users/flaiover.md`, say that the board's Import shows its instance as a story.
- Name both in `template/CHANGELOG.md`.

Waits for both formulas. Runs beside the two runs.

## Done when

The markdown lint, the settings-index test, and `flai check --strict` pass.

## Notes
