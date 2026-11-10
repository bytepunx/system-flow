---
id: T-2932
type: task
nature: remediation
title: An ADR records that the orchestrator accepts and verifies through MCP calls that run on the host, and that its shell's accept and verify fail closed in a container
status: done
parent: S-0687
owner: alex
created: 2026-11-10T08:49:27Z
updated: 2026-11-10T08:55:15Z
transitions:
  - to: ready
    at: 2026-11-10T08:51:02Z
    by: agent-S-0687
  - to: in-progress
    at: 2026-11-10T08:51:02Z
    by: agent-S-0687
  - to: done
    at: 2026-11-10T08:55:15Z
    by: agent-S-0687
stream: S-0687
tags: [flai, adr]
touches: [design/adrs/README.md, design/adrs/0278-the-orchestrator-accepts-through-an-mcp-tool-that-runs-flai-accept-on-the-host.md]
usage:
  source: log
  seconds: 253
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 54
      output: 12160
      cache_read: 9259761
      cache_write: 204229
      cost: 3.233
---
# T-2932 An ADR records that the orchestrator accepts and verifies through MCP calls that run on the host, and that its shell's accept and verify fail closed in a container

## Work

Propose the remedy for I-0256 from its instance, as the story's goal asks, and record it with `flai adr new` before any code. The proposal the plan assumes: a new orchestrator-only MCP tool, `accept`, whose handler runs `flai accept --by orchestrator` on the host, where the MCP server runs with the host's configuration, as the host API's `accept.run` does (`flai/internal/hostapi/writes_items.go`), with a dry run that answers the blockers; the orchestrator verifies through the MCP tool `verify`, the close-out job (ADR-0229), and reads the result with `flai verify --last`; and, run as the orchestrator with a `FLAI_CONFIG` inside a role's cache folder (`.flai-cache/agents/<role>`), the shell's `flai accept` and tier-running `flai verify` refuse, naming those tools. Record the alternative weighed and left: a copy of the host actions in the role's cache, which a shell write could change unseen by the guard (ADR-0220). Refine ADR-0247 and ADR-0093. If reading the code shows the proposal does not fit, write the ADR for what does and say so on the plan's thread before the next layer starts. The ADR's file name in this task's touches is predicted; correct it to the one `flai adr new` writes. This task waits for nothing: the layers after it build what it decides.

## Done when

- [ ] An accepted ADR names the MCP route for the orchestrator's acceptance and verification, the fail-closed refusal of the shell's, and the alternative left, with its reasons
- [ ] `design/adrs/README.md` lists it, as flai generates it

## Notes
