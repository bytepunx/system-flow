---
id: T-1846
type: task
nature: feature
title: The user docs and the dashboard's design describe the digest's measure and chart, and the analyzer is told to read it on a bottlenecks run
status: done
parent: S-0434
owner: alex
created: 2026-10-16T06:08:45Z
updated: 2026-10-24T22:24:03Z
transitions:
  - to: ready
    at: 2026-10-24T22:20:43Z
    by: system-flow
  - to: in-progress
    at: 2026-10-24T22:20:44Z
    by: system-flow
  - to: done
    at: 2026-10-24T22:24:03Z
    by: agent-S-0434
stream: S-0434
tags: [cli, dashboard, template]
touches: [docs/users/flai.md, docs/users/flaiover.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, template/template.yaml]
after: [T-1844, T-1845]
usage:
  source: log
  seconds: 199
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 37
      output: 11895
      cache_read: 1992095
      cache_write: 98984
      cost: 1.1972
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
