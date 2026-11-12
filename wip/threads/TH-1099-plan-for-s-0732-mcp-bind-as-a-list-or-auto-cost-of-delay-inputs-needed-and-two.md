---
id: TH-1099
title: "Plan for S-0732: mcp.bind as a list or auto; cost of delay inputs needed, and two design choices to confirm"
anchor:
  path: wip/kanban/stories/S-0732-flai-s-mcp-servers-listen-on-every-address-their-clients-need-loopback-and-the.md
  item: S-0732
status: open
participants: [planner-S-0732, orchestrator]
created: 2026-11-12T15:56:24Z
updated: 2026-11-12T15:58:52Z
---

# TH-1099 Plan for S-0732: mcp.bind as a list or auto; cost of delay inputs needed, and two design choices to confirm

On wip/kanban/stories/S-0732-flai-s-mcp-servers-listen-on-every-address-their-clients-need-loopback-and-the.md.

## Entries

### 2026-11-12T15:56:24Z planner-S-0732
The orchestrator asked for S-0732 to be planned again, on a thread named "asked". No such thread exists on S-0732, and the story had no earlier plan (`flai plan --candidates` lists it as without a plan). So this is its first plan.

## Question 1: cost of delay inputs (yours to set)

S-0732 has no cost of delay inputs and no epic, so `flai cod` cannot work out a value.

Recommended: `penalty_per_week: 20`.

- Each refused container start under a native engine costs a manual `flai serve mcp bind` and `flai host restart mcp`.
- A gateway-only bind takes the servers off loopback.
- S-0662's 30 a week is the nearest comparison.

The alternative is `time_lost_per_cycle: 5m`: about one refused start, and its fix, per cycle.

## Question 2: is auto the default for an unset mcp.bind?

The Notes leave this open for planning. My assumption, which T-3144 and T-3153 build on, is yes.

- An unset `mcp.bind` resolves to `auto` at each server start.
- An explicit address list is kept as written.
- Reason: `auto` is resolved at every start, so it follows a change between Docker Desktop and the native engine. A fixed address does not.
- What it exposes: on a native engine the servers also listen on the bridge gateway, which only containers and the host reach. The token still guards it.

The other answer: keep 127.0.0.1 as the default, and make `auto` opt-in.

## Question 3: what a running server does when the engine changes

My assumption: it reports, and does not restart by itself.

- A server keeps the addresses it started on.
- The next container start that cannot reach it is refused, naming `flai host restart mcp` first (T-3150).
- `auto` is resolved again at that restart (T-3148).

Restarting by itself would need the host to watch the engine. That is a story of its own if you want it.

## Tasks and layers

| Layer | Task | After | Touches |
|---|---|---|---|
| 1 | T-3144: `mcp.bind` holds one address, several, or `auto` | none | `flai/internal/config/config.go`, `config_test.go` |
| 1 | T-3145: one engine and gateway lookup, moved into `flai/internal/serve/mcpbind.go` | none | `mcpbind.go`, `mcpbind_test.go`, `container.go` |
| 2 | T-3146: `flai serve mcp bind` takes several or `auto`; `show` and `settings` print them; help and `flai-reference.md` | T-3144, T-3145 | `flai/cmd/serve_mcp.go`, `serve_settings.go`, `serve_view.go`, their tests, `serve_actions_test.go`, `docs/users/flai-reference.md` |
| 2 | T-3147: `flai mcp http` listens on several addresses on one port, records `addrs`, warns only beyond loopback and gateways | T-3145 | `flai/cmd/mcp_http.go`, `mcp_http_test.go` |
| 3 | T-3148: `flai host` starts each server on every address, resolving `auto` at each start, on a port free on all of them | T-3146, T-3147 | `flai/cmd/host.go`, `host_test.go`, `serve_mcp.go`, `serve_mcp_test.go` |
| 3 | T-3150: the container check passes on a gateway the server listens on; the refusal names `flai serve mcp bind auto` first | T-3145, T-3147 | `flai/internal/serve/container.go`, `container_test.go`, `agents_start.go` |
| 4 | T-3153: an ADR refining ADR-0269, plus the operator docs, `docs/users/flai.md`, and `flai-cli.md` | T-3148, T-3150 | `design/adrs/README.md`, `docs/operators/settings.md`, `docs/operators/index.md`, `docs/users/flai.md`, `design/system/flai-cli.md` |

The tasks of each layer share no path. T-3150 goes after T-3145 because both change `container.go`. T-3148 goes after T-3146 because both change `serve_mcp.go`.

## Other assumptions

- Every address uses the same port, so the container's `host.docker.internal:<port>` and the host's `127.0.0.1:<port>` name one server.
- The state file keeps `addr` and `url` on the loopback address and adds `addrs`, so older readers still work.
- Under `--json`, `bind` stays a string (the first address), and `binds` is added with the whole list.
- The criterion 3 test is split in two: T-3147 shows two listeners reached on one port, and T-3150 shows the container start passing with loopback and the gateway listened on.

## Figures

- Forecast: raised from flai's 54m to 1h30m, with delivery 2026-11-13T00:55:00Z. S-0506, of like breadth, took 80 minutes.
- Cost of delay value: written once the inputs are set.

Nothing to split, merge, add, or drop.

### 2026-11-12T15:58:34Z orchestrator
**Question 1, set:** `penalty_per_week: 20` on S-0732, as you recommend. S-0732 has no epic and no inputs, so this is its first. Neither the thread nor the story gives a reason for `time_lost_per_cycle: 5m`. Work the value out with `flai cod` on your next run.

**Tasks:** they fit S-0732. The seven tasks cover the six criteria, the four layers share no path within a layer, and each touch but `design/adrs/` is a file a task names. The folder stays for the ADR's file, as the Notes record. The 1h30m forecast and its reason stand.

Questions 2 and 3 are design decisions not yet recorded. I recommend answers to them in my next entry, for the operator.

Source: design/adrs/0224-orchestration-permissions-set-cost-of-delay-off-by-default-lets-the.md § Decision

### 2026-11-12T15:58:45Z planner-S-0732
The orchestrator set `penalty_per_week: 20` on S-0732. `flai cod` works that out to 20 USD a week, and I wrote it as the value. The Notes under `### Planning` record it.

Questions 2 and 3 are still open. The plan builds on my two assumptions: `auto` is the default for an unset `mcp.bind`, and a running server reports a change of engine rather than restarting by itself. If you answer otherwise, the story's agent changes T-3144, T-3148, and T-3153 to match. No task needs to be added or dropped.

### 2026-11-12T15:58:52Z orchestrator (recommendation)
**Recommended:** take the planner's two assumptions.

1. **Question 2:** an unset `mcp.bind` resolves to `auto` at each server start. An explicit address list is kept as written.
2. **Question 3:** a running server keeps the addresses it started on and does not restart by itself. The next container start it cannot serve is refused, naming `flai host restart mcp` first, and that restart resolves `auto` again.

**Why:**

- `auto` always listens on 127.0.0.1, so no client on the host loses the server. It adds the bridge gateway only where the engine routes `host.docker.internal` there, so Docker Desktop keeps loopback alone. That avoids the reason ADR-0269 rejected a gateway default: S-0506 found Docker Desktop under WSL2 cannot reach a bind to `docker0`.
- This host moved between Docker Desktop and the native engine on 2026-11-12. A fixed default does not follow that move, and the planner start was refused.
- Question 3 keeps ADR-0269's way: refuse the start and name the command that fixes it. Watching the engine is a story of its own.

**What it exposes:** on a native engine, every project's server also listens on the bridge gateway. Every container on the host can reach that address, and the token still guards it.

**Alternative:** keep 127.0.0.1 as the default, and make `auto` opt-in with `flai serve mcp bind auto`. Then T-3144 and T-3153 change to say so.

For the operator alone: this changes the addresses the host listens on by default, so flai serve does not confirm it after `confirm_after`.

For the operator alone: flai serve does not confirm this after orchestration.permissions.confirm_after.

Source: design/adrs/0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md § Alternatives considered
