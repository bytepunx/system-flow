---
id: T-3607
type: task
nature: improvement
title: The strategic agents, flai CLI, manifest, workflow, and publishing designs, and the user, operator, and contributor guides, describe a publish as a release_publish call
status: done
parent: S-0837
owner: alex
created: 2026-11-17T12:38:28Z
updated: 2026-11-28T00:40:51Z
transitions:
  - to: ready
    at: 2026-11-28T00:31:25Z
    by: agent-S-0837
  - to: in-progress
    at: 2026-11-28T00:31:25Z
    by: agent-S-0837
  - to: done
    at: 2026-11-28T00:40:51Z
    by: agent-S-0837
stream: S-0837
tags: [docs, flai]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, design/system/workflow.md, design/system/pushing-from-the-board.md, docs/users/flai.md, docs/operators/index.md, docs/operators/settings.md, docs/contributors/template.md, design/system/template.md]
after: [T-3601, T-3602, T-3603, T-3605]
usage:
  source: log
  seconds: 531
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 115
      output: 673
      cache_read: 6479201
      cache_write: 177798
      cost: 2.5614
---
# T-3607 The strategic agents, flai CLI, manifest, workflow, and publishing designs, and the user, operator, and contributor guides, describe a publish as a release_publish call

## Work

This task waits for T-3601, T-3602, T-3603, and T-3605, because it describes what they built.

Rewrite each passage that describes a publish as a story or names the publish formula or `scripts/fetch.sh`. A publish the orchestrator starts is one `release_publish` call: it fetches, evaluates again, publishes, or answers `covered_by`, and makes no story. The explorer found these places:

| File | Where |
|------|-------|
| `design/system/strategic-agents.md` | § Formulas, § What starts a formula, § Publishing a release (S-0222), § Its permissions and the guard |
| `design/system/flai-cli.md` | the `flai formula`, `flai release`, and `release_publish` entries under § Commands |
| `design/system/project-manifest.md` | `orchestration.permissions.publish` |
| `design/system/workflow.md` | the formula's instance in the transitions |
| `design/system/pushing-from-the-board.md` | the opening current-state paragraphs |
| `docs/users/flai.md` | the publish section, the formula section's publish paragraph, and the `stop` exits it names |
| `docs/operators/index.md` | the orchestrator's `publish` permission and the formula host action |
| `docs/operators/settings.md` | the `publish` permission's row |
| `docs/contributors/index.md`, `docs/contributors/template.md` | the formulas the template ships |

- Search the design and docs folders for `publish formula`, `formula run publish`, and `publish.md`. Change the current-state text. Leave a finding's history, such as `design/system/dashboard-host-channel.md` and `design/system/server-performance.md`, as it was written.
- Bump each changed document's `updated`.

## Done when

- No design or docs passage describes a publish the orchestrator starts as a story or names the retired formula, but as history.
- `flai check --strict`, the markdown lint, and the `settings-index` tier pass.

## Notes

- Drafted by planner-S-0837.
- `docs/users/flai-reference.md` is generated from the help, and T-3602 regenerates it.
