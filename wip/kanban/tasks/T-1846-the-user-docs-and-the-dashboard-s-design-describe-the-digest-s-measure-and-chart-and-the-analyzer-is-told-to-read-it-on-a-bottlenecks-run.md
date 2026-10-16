---
id: T-1846
type: task
nature: feature
title: The user docs and the dashboard's design describe the digest's measure and chart, and the analyzer is told to read it on a bottlenecks run
status: backlog
parent: S-0434
owner: alex
created: 2026-10-16T06:08:45Z
updated: 2026-10-16T06:08:45Z
transitions: []
stream: S-0434
tags: [cli, dashboard, template]
touches: [docs/users/flai.md, docs/users/flaiover.md, design/system/flaiover-dashboard.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-1844, T-1845]
---
# T-1846 The user docs and the dashboard's design describe the digest's measure and chart, and the analyzer is told to read it on a bottlenecks run

## Work

- `docs/users/flai.md`: the `digest` section of `flai stats`.
- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md`: the digest chart on the charts page.
- `strategic-agents.md` (template and here), "As the analyzer": on a bottlenecks run, read `digest.weeks[]` and say whether the threads to the operator per completed story, and the repeats, fell.
- `template/CHANGELOG.md` and the version in `template/template.yaml`.

Waits for the CLI and chart tasks, so the docs describe what was built.

## Done when

The documents say what the code does, and `flai check --strict` and the markdown lint report nothing on them.

## Notes

The analyzer's prompt (`analyzePrompt`) already reads every figure from `flai stats --json`; the convention is enough.
