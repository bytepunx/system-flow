---
id: T-3048
type: task
nature: improvement
title: Record in an ADR that the operator's outside folders reach a story's agent and its sub-agents as --add-dir and container mounts, and that Monitor joins the default tools
status: backlog
parent: S-0701
owner: alex
created: 2026-11-10T10:34:59Z
updated: 2026-11-10T10:34:59Z
transitions: []
stream: S-0701
tags: [adr, design]
touches: [design/adrs/README.md, design/system/agent-adapters.md]
---
# T-3048 Record in an ADR that the operator's outside folders reach a story's agent and its sub-agents as --add-dir and container mounts, and that Monitor joins the default tools

## Work

Write the ADR for remedy A as TH-1025 settles it, with `flai adr new`. The operator names outside folders per harness on the host (`flai serve agent harness <name> --outside <dir>...`). flai serve puts them in the start request. The claude-code adapter passes each as `--add-dir`, which the session's sub-agents inherit. The container harness mounts each at its own path and refuses home and any folder that holds the project. `Monitor` joins the default `--allowedTools`. The start prompt names the granted folders. `permission_prompt` stays as ADR-0086 and ADR-0106 leave it. Say why ADR-0038 holds: a story still names no folder, and only the operator widens what an agent reaches. List remedies B and C from TH-1025 as alternatives considered. Cite D-0480.

Update `design/system/agent-adapters.md` § Today › Permissions and § Sub-agents and strategic agents to say what the agent and its sub-agents may reach outside the project, linking the ADR.

This task waits for nothing. Every other task builds what this one decides.

## Done when

- The ADR is in `design/adrs/` and listed in `design/adrs/README.md`, and `flai check --strict` is clean.
- `design/system/agent-adapters.md` says which outside folders a story's agent and its sub-agents reach, how, and where the operator sets them.

## Notes
