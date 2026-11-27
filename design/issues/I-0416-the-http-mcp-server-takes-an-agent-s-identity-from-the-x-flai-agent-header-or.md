---
id: I-0416
title: The HTTP MCP server takes an agent's identity from the X-Flai-Agent header or the client name, so any holder of the project's one token is served as any agent
class: impression
status: open
count: 1
cost: 5m
first_reported: 2026-11-27T22:46:58Z
last_reported: 2026-11-27T22:46:58Z
updated: 2026-11-27T22:46:58Z
usage:
  strategic:
    - kind: analyzer
      seconds: 85
      estimated: true
      models:
        - model: claude-fable-5-1
          input: 55
          output: 5532
          cache_read: 391176
          cache_write: 29513
          cost: 0.9652
        - model: claude-opus-5-5
          input: 29
          output: 4664
          cache_read: 1264948
          cache_write: 28114
          cost: 0.487
---

# I-0416 The HTTP MCP server takes an agent's identity from the X-Flai-Agent header or the client name, so any holder of the project's one token is served as any agent

## Description
The HTTP MCP server takes an agent's identity from the X-Flai-Agent header or the client name, so any holder of the project's one token is served as any agent

## Instances

### 2026-11-27T22:46:58Z
Report: design/analysis/2026-11-27-all.md.
Read in the code today; no occurrence observed. A story's agent that sends X-Flai-Agent: orchestrator, or another story's agent name, with the project token is served as that agent, with its tools and guard rules. A prompt injection in a thread or document is the way such a header gets sent.

## Impact
flai/internal/mcphttp/mcphttp.go:174-207 handler.ServeHTTP: after the Origin refusal and the constant-time bearer check (ADR-0216), agent := SafeAgent(r.Header.Get(AgentHeader)), else the client name from the body, then r.WithContext(context.WithValue(r.Context(), ctxKey{}, agent)) picks the per-agent server. One token per project (ADR-0033) is held by every story's container and every strategic agent's. ADR-0143 holds a served agent to its own name on the CLI side through flai guard; over HTTP only the header says who calls. No instance observed.

## Remediation

Found by the analysis in [2026-11-27-all.md](../analysis/2026-11-27-all.md).
