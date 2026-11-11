---
id: T-3048
type: task
nature: improvement
title: Record in an ADR that the operator's outside folders reach a story's agent and its sub-agents as --add-dir and container mounts, and that Monitor joins the default tools
status: done
parent: S-0701
owner: alex
created: 2026-11-10T10:34:59Z
updated: 2026-11-11T23:09:37Z
transitions:
  - to: ready
    at: 2026-11-11T23:06:31Z
    by: agent-S-0701
  - to: in-progress
    at: 2026-11-11T23:06:31Z
    by: agent-S-0701
  - to: done
    at: 2026-11-11T23:09:37Z
    by: agent-S-0701
stream: S-0701
tags: [adr, design]
touches: [design/adrs/README.md, design/system/agent-adapters.md, design/adrs/0289-flai-serve-gives-a-story-s-agent-and-its-sub-agents-the-outside-folders-the.md]
usage:
  source: log
  seconds: 186
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 44
      output: 258
      cache_read: 1616339
      cache_write: 84809
      cost: 0.6319
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
