---
id: T-0409
type: task
nature: remediation
title: The claude-code harness names its agent to flai's MCP server itself, so an environment that sets FLAI_AGENT cannot merge every started agent into one name
status: done
parent: S-0114
owner: alex
created: 2026-09-24T08:34:03Z
updated: 2026-09-24T08:39:45Z
transitions:
  - to: ready
    at: 2026-09-24T08:37:57Z
    by: agent-S-0114
  - to: in-progress
    at: 2026-09-24T08:37:58Z
    by: agent-S-0114
  - to: done
    at: 2026-09-24T08:39:45Z
    by: agent-S-0114
stream: S-0114
tags: []
touches: [flai/internal/harness, flai/cmd]
---
# T-0409 The claude-code harness names its agent to flai's MCP server itself, so an environment that sets FLAI_AGENT cannot merge every started agent into one name

## Work

On this host every agent flai serve started connected to MCP as `system-flow`: the project's `.claude/settings.local.json` sets `FLAI_AGENT=system-flow` in its `env`, and Claude Code puts that into the MCP server it spawns, over what flai serve set. So the shared cursor `.flai-cache/mcp/system-flow.json` was nobody's own sign, and each of its own agents ending held the next start back six minutes (the gaps the operator saw between 08:05 and 08:11, and 08:21 and 08:27 on 2026-09-24). It is also I-0037's root cause. Give `flai mcp` an `--agent` flag that names the agent over `FLAI_AGENT`, and have the `claude-code` adapter pass `["mcp", "--agent", <name>]` in its `--mcp-config`. Tests in `harness_test.go`; bump I-0037 with this instance and record the remediation.

## Done when

- [x] `flai mcp --agent <name>` serves as that agent whatever `FLAI_AGENT` says
- [x] The claude-code adapter's MCP config names the agent with `--agent`
- [x] I-0037 records this occurrence and its remediation
- [x] `make flai-test` passes

## Notes
