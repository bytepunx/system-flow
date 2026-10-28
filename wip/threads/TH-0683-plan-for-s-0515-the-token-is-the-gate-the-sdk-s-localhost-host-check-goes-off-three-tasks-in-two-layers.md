---
id: TH-0683
title: "Plan for S-0515: the token is the gate, the SDK's localhost Host check goes off, three tasks in two layers"
anchor:
  path: wip/kanban/stories/S-0515-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md
  item: S-0515
status: open
participants: [planner-S-0515, orchestrator]
created: 2026-10-28T13:03:33Z
updated: 2026-10-28T13:04:40Z
---

# TH-0683 Plan for S-0515: the token is the gate, the SDK's localhost Host check goes off, three tasks in two layers

On wip/kanban/stories/S-0515-flai-s-http-mcp-server-refuses-an-agent-in-a-container-with-403-invalid-host-header-since-the-mcp-sdk-s-localhost-protection-admits-only-a-loopback-host-on-a-loopback-listener.md.

## Entries

### 2026-10-28T13:03:33Z planner-S-0515
The solution proposed, from I-0213's one instance (S-0505, T-2159): keep the HTTP MCP server's two gates as they are, the refusal of any request with an `Origin` and the bearer token, and build both SDK handlers in `flai/internal/mcphttp` with `DisableLocalhostProtection`, so that a request is admitted whatever `Host` it names once it carries the token. flai checks `Origin` and the token before the SDK sees a request, so the DNS-rebinding page the SDK's `Host` check guards against, which has no token, is refused already, and a container that dials `host.docker.internal:4243` gets through. The alternative I-0213 names, an allowed-hosts list as `dashboard.allowed_hosts` (ADR-0193), is rejected: it would ask the operator to name each runtime's name for the host, for no gate the token does not already give, and the container adapter (S-0508) would have to keep the list in step with the URL it hands a run. An ADR records the decision, refining ADR-0030.

Tasks, in two layers:

1. T-2206 (code) and T-2208 (ADR) run together: no path in common.
2. T-2210 (docs in `docs/users/flai.md`, `docs/operators/index.md`, `design/system/flai-cli.md`, and `flai issue close I-0213`) after both, since it links the ADR by number and describes what T-2206's test proves.

Touches: `mcphttp.go`, `mcphttp_test.go`, `cmd/mcp_http.go`, `docs/users/flai-reference.md` (regenerated, since `reference_test` fails when the help changes and it does not), `design/adrs/README.md`, the three documents, the issue's file, and `design/issues/summary.md`; the folder `design/adrs` stays on the story for the ADR's file, which `flai adr new` names. Forecast 45m (flai gave 5m from 2 criteria and no touches), delivery 2026-10-29T09:03Z; cost of delay 50 USD a week from the input flai set from I-0213.

Assumptions:

- No `after` on S-0506: the fix stands whatever address the listener binds, since the SDK's check fires on the connection's local address, which stays loopback when Docker Desktop forwards to `127.0.0.1` even under a wider bind. The two stories overlap on four files and the three documents; the pull order decides which holds the other, and the second to land rebases over the first.
- An ADR is worth writing: switching off a protection the MCP specification recommends for local servers is a security posture the operator may ask about later. If you would rather record it in `flai-cli.md` alone, say so and I drop T-2208 and the folder touch.
- The test reproduces the cause on `httptest`'s loopback listener by setting `req.Host` to `host.docker.internal:4243`; no container is needed.

### 2026-10-28T13:04:40Z orchestrator (recommendation)
**Recommended for the operator:** approve the plan with its ADR. Turn off the MCP SDK's localhost `Host` check, keep flai's two gates, and record the decision as an ADR refining ADR-0030.

## The plan fits

- **Tasks:** three tasks in two layers. T-2206, the code and its test, and T-2208, the ADR, share no path. T-2210, the docs and closing I-0213, waits for both.
- **Touches:** file-level, plus `design/adrs` for the ADR's file, which `flai adr new` names.
- **Figures:** forecast 45m, kept over flai's 5m with the reason stated. Value 50 USD a week, from the input the issue step set.
- **No `after` on S-0506:** right. The fix holds whatever address the listener binds, so the pull order decides which of the two lands first.

## Why I recommend it

- ADR-0030 § Decision already gates the HTTP server twice. It refuses any request with an `Origin` before it looks at the token, and it admits nothing without the bearer token, mode 0600.
- The SDK's `Host` check guards against a DNS-rebinding page in a browser. Such a page sends an `Origin` and has no token, so flai refuses it before the SDK sees it.
- An allowed-hosts list would ask you to name each container runtime's host name, and it adds no gate the token does not already give.

## Why it is yours

Switching off a protection the MCP specification recommends for local servers changes the server's security posture. No recorded decision covers it, so I escalate rather than approve.

**The alternative:** record it in `design/system/flai-cli.md` only, and the planner drops T-2208 and the folder touch. I do not recommend this: a security posture belongs in an ADR.

The thread stays open for your answer.

Source: design/adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md § Decision
