---
id: T-1773
type: task
nature: feature
title: The template and this repository have design/workflows/formulas with its README, and the design, the convention, and the docs say what a formula holds and how to run one
status: backlog
parent: S-0421
owner: alex
created: 2026-10-16T05:47:33Z
updated: 2026-10-21T12:40:44Z
transitions: []
stream: S-0421
tags: [cli, template]
touches: [template/root/design/workflows/formulas/README.md, template/root/design/workflows/README.md, design/workflows/formulas/README.md, design/workflows/README.md, design/system/repository-layout.md, design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, template/template.yaml, design/system/flai-cli.md, design/system/project-manifest.md]
after: [T-1771, T-1772]
---
# T-1773 The template and this repository have design/workflows/formulas with its README, and the design, the convention, and the docs say what a formula holds and how to run one

## Work

- Add `design/workflows/formulas/README.md`, template and here, saying what the folder is for; and `design/workflows/README.md` in both when S-0414 left none, else a line for formulas in it.
- Add a section on formulas to `design/system/strategic-agents.md`, or start a design document of its own if that section outgrows it, linking the ADR.
- Add the folder to `repository-layout.md`.
- Say in the convention, template and here, that a repeatable process is a formula.
- Describe the commands in `docs/users/flai.md`, and add how to write one to `docs/operators/index.md`.
- Update `template/CHANGELOG.md` and bump the template's version.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes

- Title and first Work bullet rewritten by planner-S-0421 to TH-0471's answer 2, as approved on TH-0581.
