---
id: T-0892
type: task
nature: feature
title: An ADR, strategic-agents.md, flai-cli.md, the manifest design, and the operator and user guides describe the orchestrator
status: done
parent: S-0218
owner: alex
created: 2026-10-05T04:46:48Z
updated: 2026-10-05T08:15:57Z
transitions:
  - to: ready
    at: 2026-10-05T07:59:00Z
    by: agent-S-0218
  - to: in-progress
    at: 2026-10-05T07:59:00Z
    by: agent-S-0218
  - to: review
    at: 2026-10-05T08:15:57Z
    by: agent-S-0218
  - to: done
    at: 2026-10-05T08:15:57Z
    by: agent-S-0218
stream: S-0218
tags: [flai]
touches: [design/adrs, design/system/strategic-agents.md, design/system/flai-cli.md, design/system/project-manifest.md, docs/operators/settings.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-0881, T-0882, T-0883, T-0887]
usage:
  source: log
  seconds: 1017
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 158
      output: 64963
      cache_read: 9260076
      cache_write: 235320
      cost: 4.4688
---
# T-0892 An ADR, strategic-agents.md, flai-cli.md, the manifest design, and the operator and user guides describe the orchestrator

## Work

Write down what T-0881, T-0882, T-0883, and T-0887 built:

- A new ADR in `design/adrs`, refining ADR-0082 and ADR-0060: `flai serve` runs one orchestrator per project behind the `orchestrate` host action, restarting and stopping it with the action; `orchestration.agent` picks its agent; `orchestration.permissions`, each off by default, are enforced by `flai guard` and each refusal names the permission and is logged in `orchestrator.md`; its prompt is a loop of decisions logged with `activity_log` and waits on `wait_for_events`.
- `design/system/strategic-agents.md`: the table's Host action and Started by rows for the orchestrator, and a section "The orchestrator" shaped like the planner's: starting and stopping it, its agent, how it runs, what it is told, the guard's table of permissions, and the run and how it ends. Trim "The orchestrator and the analyzer" to the analyzer.
- `design/system/flai-cli.md`: the `orchestrate` host action under `flai serve`, and the orchestrator's rules under `flai guard`.
- `design/system/project-manifest.md` and `docs/operators/settings.md`: `orchestration.permissions` and `orchestration.agent`, each value and default, beside T-0809's `policy` and `release`.
- `docs/users/flai.md` and `docs/users/flai-reference.md`: how to turn the orchestrator on and off, what it may do under each permission, and where its decisions are logged.

This task waits for the four code tasks, so it describes what they built. It runs with T-0889, whose paths it does not share.

## Done when

- the ADR is accepted in `design/adrs` and linked from `strategic-agents.md`
- each document above describes the action, the permissions, the guard's refusals, and the decision log as built
- `flai check --strict` and the markdown lint report nothing on them

## Notes
