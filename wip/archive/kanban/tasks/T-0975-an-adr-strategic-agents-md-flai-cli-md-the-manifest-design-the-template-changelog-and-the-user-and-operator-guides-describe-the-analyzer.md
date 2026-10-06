---
id: T-0975
type: task
nature: feature
title: An ADR, strategic-agents.md, flai-cli.md, the manifest design, the template changelog, and the user and operator guides describe the analyzer
status: done
parent: S-0223
owner: alex
created: 2026-10-05T05:48:06Z
updated: 2026-10-06T20:48:13Z
transitions:
  - to: ready
    at: 2026-10-06T20:35:49Z
    by: agent-S-0223
  - to: in-progress
    at: 2026-10-06T20:35:49Z
    by: agent-S-0223
  - to: review
    at: 2026-10-06T20:48:12Z
    by: agent-S-0223
  - to: done
    at: 2026-10-06T20:48:13Z
    by: agent-S-0223
stream: S-0223
tags: [flai]
touches: [design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md, template/CHANGELOG.md, template/template.yaml]
after: [T-0949, T-0957, T-0960, T-0968, T-0971, T-0972, T-0973]
usage:
  source: log
  seconds: 743
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 293
      output: 117147
      cache_read: 18725624
      cache_write: 456791
      cost: 8.612
---
# T-0975 An ADR, strategic-agents.md, flai-cli.md, the manifest design, the template changelog, and the user and operator guides describe the analyzer

## Work

Record the decision in an ADR (`flai adr new`): the analyzer runs behind the `analyze` host action, on demand (`flai analyze`, `analyze.run`, the MCP tool) or on `analysis.schedule`, one run per project at a time; it writes one report under `design/analysis/` and edits nothing else, which `flai guard` enforces; and its activity entry names the report.

Rewrite `design/system/strategic-agents.md`'s "The orchestrator and the analyzer" so that the analyzer has its own section as the planner has: starting it, its agent, its prompt, the report and its index, the guard, the schedule, and the run and how it ended. Add `flai analyze` and the MCP tool to `design/system/flai-cli.md`, and `analysis.agent` and `analysis.schedule` to `design/system/project-manifest.md`. Describe them for users in `docs/users/flai.md` ("The planner, the orchestrator, and the analyzer") and `docs/users/flai-reference.md`, and the `analyze` host action and the `analysis` block for operators in `docs/operators/settings.md`.

Add the template's entry for the analyzer's definition, the settings hook, and the `design/analysis` folder to `template/CHANGELOG.md`, and bump `template/template.yaml`.

It waits for every code task, T-0949, T-0957, T-0960, T-0968, T-0971, T-0972, and T-0973, so that it describes what they built. The dashboard task writes its own documents and runs with this one.

## Done when

- the ADR is accepted in `design/adrs/` and listed in its README
- each document named above describes the analyzer as built, with no behaviour left undescribed
- the template changelog has the entry and `template.yaml` its new version
- `flai check --strict` and the markdown lint pass

## Notes
