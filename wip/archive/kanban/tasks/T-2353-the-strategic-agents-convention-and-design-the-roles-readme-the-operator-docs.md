---
id: T-2353
type: task
nature: improvement
title: The strategic-agents convention and design, the roles README, the operator docs, and the changelog say a strategic agent's definition is generated from its role file and changed there alone
status: done
parent: S-0553
owner: alex
created: 2026-10-30T20:47:44Z
updated: 2026-11-02T00:22:14Z
transitions:
  - to: ready
    at: 2026-11-02T00:06:49Z
    by: agent-S-0553
  - to: in-progress
    at: 2026-11-02T00:06:50Z
    by: agent-S-0553
  - to: done
    at: 2026-11-02T00:22:14Z
    by: agent-S-0553
stream: S-0553
tags: [template]
touches: [design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, design/system/strategic-agents.md, design/workflows/roles/README.md, template/root/design/workflows/roles/README.md, docs/operators/index.md, template/CHANGELOG.md]
after: [T-2348]
usage:
  source: log
  seconds: 924
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 45
      output: 13232
      cache_read: 3988000
      cache_write: 68189
      cost: 1.5115
---
# T-2353 The strategic-agents convention and design, the roles README, the operator docs, and the changelog say a strategic agent's definition is generated from its role file and changed there alone

## Work

Say in the strategic-agents convention, here and in the template, that a strategic agent's prompt and definition are changed in its role file alone, ending D-0411's double edit. In `design/system/strategic-agents.md`, rewrite the bullets that name `.claude/agents/<name>.md` as the definition, and the orchestrator's refusal naming `flai upgrade`, to the generation. Say in the roles README pair what `flai role definition` does. In `docs/operators/index.md`, replace the lines that say the orchestrator needs `.claude/agents/orchestrator.md` with what the operator does once, as the ADR says. Add the change to `template/CHANGELOG.md`.

Third layer: it waits for the ADR and for the command it describes. It shares no path with the `.claude/agents` pair task.

## Done when

No document here or in the template says a strategic agent's definition is read from `.claude/agents` or that a prompt change is made twice. `flai check --strict` and the markdown lint report nothing.

## Notes

Drafted by the planner of E-0024 on 2026-10-30.
