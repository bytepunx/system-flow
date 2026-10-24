---
id: T-1830
type: task
nature: feature
title: The strategic-agents convention gains As the curator, and the design and the user and operator docs describe the role, its action, and digest.agent
status: done
parent: S-0431
owner: alex
created: 2026-10-16T06:06:28Z
updated: 2026-10-24T01:06:47Z
transitions:
  - to: ready
    at: 2026-10-24T00:54:11Z
    by: agent-S-0431
  - to: in-progress
    at: 2026-10-24T00:54:12Z
    by: agent-S-0431
  - to: done
    at: 2026-10-24T01:06:46Z
    by: agent-S-0431
stream: S-0431
tags: [cli, template]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/conventions/README.md, template/root/design/conventions/README.md, design/system/strategic-agents.md, design/system/agent-narrative.md, design/system/agent-adapters.md, design/system/project-manifest.md, design/system/flai-cli.md, docs/operators/index.md, docs/operators/settings.md, docs/users/flai.md, wip/agents/README.md, template/root/wip/agents/README.md, template/CHANGELOG.md, template/template.yaml, design/system/metrics.md, flai/internal/serve/agents.go, flai/internal/serve/curate_test.go]
after: [T-1826, T-1829]
usage:
  source: log
  seconds: 754
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 210
      output: 76293
      cache_read: 14954576
      cache_write: 493024
      cost: 7.4893
---
# T-1830 The strategic-agents convention gains As the curator, and the design and the user and operator docs describe the role, its action, and digest.agent

## Work

- `strategic-agents.md` (template and here): the curator joins the strategic agents, with "As the curator": keep the digest through `flai digest` alone, one fact per entry citing its source, merge rather than duplicate, record `asked_again`, never a secret's value, edit nothing else, author no stories; and its activity document under "Logging your activity".
- `design/system/strategic-agents.md`: a section "The curator" on the analyzer's pattern: starting it, its agent, how it runs, what it is told, its guard, its run and how it ended, what it costs.
- `design/system/project-manifest.md` and `docs/operators/settings.md`: `digest.agent`. `docs/operators/index.md`: the `curate` action. `docs/users/flai.md` and `design/system/flai-cli.md`: `flai curate`.
- `template/CHANGELOG.md` and the version in `template/template.yaml`.

Waits for the definitions and entry points tasks, so the docs describe what was built.

## Done when

The documents say what the code does, the settings test finds a row for `digest.agent`, and `flai check --strict` and the markdown lint report nothing on them.

## Notes

The convention's frontmatter `roles` gains `curate`, so the curator's pack holds it.
