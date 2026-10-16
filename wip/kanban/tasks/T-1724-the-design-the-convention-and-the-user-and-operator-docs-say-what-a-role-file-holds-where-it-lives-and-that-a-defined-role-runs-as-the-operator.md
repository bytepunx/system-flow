---
id: T-1724
type: task
nature: feature
title: The design, the convention, and the user and operator docs say what a role file holds, where it lives, and that a defined role runs as the operator
status: backlog
parent: S-0414
owner: alex
created: 2026-10-16T05:43:18Z
updated: 2026-10-16T05:50:22Z
transitions: []
stream: S-0414
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/repository-layout.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, docs/users/flai.md, docs/operators/index.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/project-manifest.md]
after: [T-1723]
---
# T-1724 The design, the convention, and the user and operator docs say what a role file holds, where it lives, and that a defined role runs as the operator

## Work

- Add a section on role files to `design/system/strategic-agents.md`, linking the ADR.
- Add `design/roles/` to `design/system/repository-layout.md`.
- Say in the convention `strategic-agents.md`, template and here, above the marker in the template, that a role is defined by its file and runs as the operator. Point to the safety note under `docs/operators/index.md` § Starting an agent when a story becomes ready.
- Describe `flai role list` and `flai role show` in `docs/users/flai.md`.
- Add how to define a role to `docs/operators/index.md`.
- Name the role files and the new findings in `template/CHANGELOG.md`.

Last layer: it describes what T-1723 built.

## Done when

The markdown lint and `flai check --strict` pass. Each document names the ADR.

## Notes
