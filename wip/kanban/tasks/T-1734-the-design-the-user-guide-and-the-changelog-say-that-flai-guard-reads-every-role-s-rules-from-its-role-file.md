---
id: T-1734
type: task
nature: feature
title: The design, the user guide, and the changelog say that flai guard reads every role's rules from its role file
status: backlog
parent: S-0416
owner: alex
created: 2026-10-16T05:44:33Z
updated: 2026-10-28T20:33:06Z
transitions: []
stream: S-0416
tags: [cli, template]
touches: [design/system/strategic-agents.md, design/system/flai-cli.md, docs/users/flai.md, template/CHANGELOG.md, design/system/agent-adapters.md, flai/cmd/guard.go, docs/users/flai-reference.md, design/system/project-manifest.md]
after: [T-1731, T-1732, T-1733]
---
# T-1734 The design, the user guide, and the changelog say that flai guard reads every role's rules from its role file

## Work

- In `design/system/strategic-agents.md`, under § The guard, § Its permissions and the guard, and § The analyzer's guard, say that each table is what the role's file states and that the named conditions stay in flai.
- Describe the guard's role path and its fail-closed refusal in `design/system/flai-cli.md` under `flai guard`, and in `docs/users/flai.md`.
- In `design/system/agent-adapters.md`, name a defined role among the agents ADR-0143's own-name rule holds, with the new ADR. Give the hook's condition in its hooks table as "`FLAI_ROLE` is set, or `FLAI_STORY` is set", and say in its contract that `FLAI_ROLE` may name a defined role.
- In `flai guard`'s help in `flai/cmd/guard.go`, say that each role's rules are read from its role file and that the hook runs the guard for any `FLAI_ROLE`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Name the settings change and the role path in `template/CHANGELOG.md`.

Last layer: it describes the whole story.

## Done when

The markdown lint and `flai check --strict` pass.

## Notes

- Planner (S-0416, 2026-10-21): the orchestrator approved, on TH-0569, adding `agent-adapters.md` and `flai guard`'s help to the Work.
