---
id: T-0789
type: task
nature: feature
title: "Design, ADR, and user and operator docs for the planner: strategic-agents.md, flai-cli.md, project-manifest.md, the user guide, and settings"
status: done
parent: S-0208
owner: alex
created: 2026-10-04T00:44:14Z
updated: 2026-10-04T03:42:18Z
transitions:
  - to: ready
    at: 2026-10-04T00:44:32Z
    by: agent-S-0208
  - to: in-progress
    at: 2026-10-04T03:28:58Z
    by: agent-S-0208
  - to: done
    at: 2026-10-04T03:42:18Z
    by: agent-S-0208
stream: S-0208
tags: []
touches: [design/system, design/adrs, docs/users, docs/operators, template/root/docs, template/CHANGELOG.md, template/template.yaml]
after: [T-0784, T-0785, T-0786, T-0787, T-0788]
usage:
  source: log
  seconds: 800
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 320
      output: 92931
      cache_read: 23263547
      cache_write: 403575
      cost: 8.9592
---
# T-0789 Design, ADR, and user and operator docs for the planner: strategic-agents.md, flai-cli.md, project-manifest.md, the user guide, and settings

## Work

- An ADR for the planner as flai serve runs it: a claude-code session in the main checkout on `planning.agent` over `agent`, run as `.claude/agents/planner.md`, behind the `plan` host action, one run per item, recorded as a story's run, and held to planning by `flai guard` through `FLAI_ROLE=plan`.
- A new `design/system/strategic-agents.md` (topics planning, orchestration, analysis) describing the strategic agents as built so far: the planner's start, run, record, refusal, prompt, guard, and activity log, linking ADR-0075, ADR-0079, and the new ADR.
- `design/system/flai-cli.md`: `flai plan`, `plan` among the host actions, `plan.run`, the MCP tool `plan`, the guard's planner rules, and `planning.agent`.
- `design/system/project-manifest.md`: `planning.agent`.
- `docs/users/flai.md` and the user guide: how to turn on `plan` and run the planner; `docs/users/flai-reference.md` and `docs/operators/settings.md` regenerated with `make flai-reference`, and settings' rows for `planning.agent` and `FLAI_ROLE`/`FLAI_ITEM` where its tests ask.

Waits for every other task: it describes what they built.

## Done when

- [x] The ADR is accepted and indexed, and the design links it
- [x] The docs say how to run the planner, and the reference and settings tests pass

## Notes
