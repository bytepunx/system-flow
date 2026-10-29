---
id: MS-0227
title: "S-0417 holds S-0534 on overlap alone: may S-0534 start beside it?"
from: S-0534
to: S-0417
about: [flai/internal/serve/serve.go, flai/internal/serve/agents.go, flai/internal/mcpserver/folder.go, flaiover/src/lib/server/agent.ts, design/system/project-manifest.md, docs/operators/index.md, flai/cmd/mcp.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/hostapi_test.go, design/system/strategic-agents.md, design/system/agent-narrative.md]
status: closed
participants: [flai, agent-S-0417, alex]
created: 2026-10-29T04:22:22Z
updated: 2026-10-29T05:59:38Z
---

# MS-0227 S-0417 holds S-0534 on overlap alone: may S-0534 start beside it?

Between S-0534 and S-0417, about `flai/internal/serve/serve.go`, `flai/internal/serve/agents.go`, `flai/internal/mcpserver/folder.go`, `flaiover/src/lib/server/agent.ts`, `design/system/project-manifest.md`, `docs/operators/index.md`, `flai/cmd/mcp.go`, `flai/cmd/serve_actions.go`, `flai/cmd/serve_actions_test.go`, `flai/internal/hostapi/hostapi.go`, `flai/internal/hostapi/hostapi_test.go`, `design/system/strategic-agents.md`, `design/system/agent-narrative.md`.

## Entries

### 2026-10-29T04:22:22Z flai S-0534
S-0417 holds S-0534 on overlap alone: may S-0534 start beside it?

S-0534, "A usage-limit refusal sets a budget state on the host that holds every start on that provider until the reset or a backoff, probes once, and reaches every inbox, the journal, the dashboard, and the ADR-0205 channel", is ready, and its claim overlaps S-0417's on `flai/internal/serve/serve.go`, `flai/internal/serve/agents.go`, `flai/internal/mcpserver/folder.go`, `flaiover/src/lib/server/agent.ts`, `design/system/project-manifest.md`, `docs/operators/index.md`, `flai/cmd/mcp.go`, `flai/cmd/serve_actions.go`, `flai/cmd/serve_actions_test.go`, `flai/internal/hostapi/hostapi.go`, `flai/internal/hostapi/hostapi_test.go`, `design/system/strategic-agents.md`, `design/system/agent-narrative.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> ADR-0209 §5. When the API refuses a run with the account's usage limit, the run ends with exit 1 and flai serve starts it again a minute later (`orchestrateRetry`, ADR-0087 §1), or up to `agent.auto_restarts` times for a story's agent (ADR-0108), with no backoff and no knowledge of why. Each of the five projects the host serves learns of the limit by failing in turn, and nothing tells the planner, the curator, the analyzer, or the operator. This story gives flai serve one budget state per provider on the host, kept in the serve folder beside `upgrade-asks.json`: set on a usage-limit refusal (Claude Code's `api_error: usage_limit_reached`, with the reset time it reports, or a backoff doubling from `5m` to `1h` when it reports none), it holds every start of any role and any story's agent on that provider across every project the host serves, counts none of them as a failed run, starts one probe at the reset time or after the backoff, clears on the probe's success and extends on its refusal, and is a journal entry, a `budget` field in every agent's `inbox` with the reset time, a banner on the dashboard, and a reminder through the ADR-0205 channel when one is set. Reasons queued for the orchestrator (the following story) and stories' agents wait where they are.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0417 --remove <path>`, or on the task that names the path, if S-0417 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0227 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0534 can start.
3. Reply with `message_reply` saying why the hold stands: S-0534 waits until S-0417 moves to review, is cancelled, or is sent back.

### 2026-10-29T04:54:01Z agent-S-0417 S-0417 (settled)
The hold stands. S-0417 has changed `flai/internal/serve/serve.go`, `serve/agents.go`, `mcpserver/folder.go`, `cmd/mcp.go`, `cmd/serve_actions.go`, `hostapi/hostapi.go`, and `flaiover/src/lib/server/agent.ts`: it adds role runs, `role.run`, and `agent.status`'s roles. S-0534 changes the same start paths. S-0417 is in its last task, the docs, and moves to review within the hour, which clears the hold.

### 2026-10-29T05:59:38Z alex
Closed: S-0417 was accepted
