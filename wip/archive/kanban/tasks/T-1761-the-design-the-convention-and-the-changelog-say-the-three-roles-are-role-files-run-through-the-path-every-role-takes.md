---
id: T-1761
type: task
nature: improvement
title: The design, the convention, and the changelog say the three roles are role files run through the path every role takes
status: done
parent: S-0419
owner: alex
created: 2026-10-16T05:46:30Z
updated: 2026-10-30T13:25:15Z
transitions:
  - to: ready
    at: 2026-10-30T13:04:07Z
    by: agent-S-0419
  - to: in-progress
    at: 2026-10-30T13:04:08Z
    by: agent-S-0419
  - to: done
    at: 2026-10-30T13:08:39Z
    by: agent-S-0419
stream: S-0419
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md, design/system/flai-cli.md, design/system/repository-layout.md, design/workflows/roles/README.md, docs/operators/index.md, docs/users/flai-reference.md, docs/users/flai.md, flai/cmd/check.go, template/root/design/workflows/roles/README.md, design/system/workflow.md, flai/cmd/check_stats_test.go]
after: [T-1754, T-1757, T-1759]
usage:
  source: log
  seconds: 271
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 67
      output: 18884
      cache_read: 6780155
      cache_write: 176399
      cost: 2.9656
---
# T-1761 The design, the convention, and the changelog say the three roles are role files run through the path every role takes

## Work

- In `design/system/strategic-agents.md`, under each role's How it runs and What it is told, say that its prompt is its role file's body and that it starts through the defined-role path. Say what the embedded baseline does for a project without role files.
- In the convention, template and here, say that a project changes a built-in role below its file's marker.
- Name the change and `roles.baseline` in `template/CHANGELOG.md`.

Last layer.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes
