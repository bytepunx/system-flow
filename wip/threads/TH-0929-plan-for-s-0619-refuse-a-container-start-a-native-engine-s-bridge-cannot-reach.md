---
id: TH-0929
title: "Plan for S-0619: refuse a container start a native engine's bridge cannot reach"
anchor:
  path: wip/kanban/stories/S-0619-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md
  item: S-0619
status: resolved
participants: [planner-S-0619, orchestrator]
created: 2026-11-09T17:04:39Z
updated: 2026-11-09T17:06:19Z
---

# TH-0929 Plan for S-0619: refuse a container start a native engine's bridge cannot reach

On wip/kanban/stories/S-0619-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md.

## Entries

### 2026-11-09T17:04:39Z planner-S-0619
S-0619 is planned: 2 tasks in 2 layers, 12 file touches, a 40m forecast, and a cost of delay of 37.50 USD a week.

**Remedy.** I-0254 names no remedy, so the plan takes the one its S-0508 instance proposes. Before `flai serve` starts a container on the bridge, it refuses the start when three things hold:

- the MCP URL's host is loopback (`mcp.bind` unset);
- the network is not `host`;
- `docker info` says the engine is not Docker Desktop.

The refusal names `flai serve mcp bind <bridge gateway>`, and `0.0.0.0` with its warning as the other choice. The default bind stays `127.0.0.1`, because S-0506 measured that Docker Desktop under WSL2 cannot reach a bind to WSL's `docker0` address.

**Tasks.**

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-2735 the preflight in `flai/internal/serve/container.go`, called from `startOf` in `agents.go`, with tests in `container_test.go` that reproduce I-0254 | none | 3 Go files |
| 2 | T-2736 the operators guide, `settings.md`, the install and agent-container runbooks, `docs/users/flai.md`, `design/system/flai-cli.md`, `design/tech/docker.md`; then `flai issue close I-0254` | T-2735, whose refusal the documents quote | 7 documents and the issue files |

**Assumptions.**

- A `docker info` that cannot answer does not refuse a start. It logs a warning, so the check never refuses on a guess.
- The engine's kind is read once per `flai serve` process, not once per start.
- When strategic agents' container starts go through a path other than `startOf`, T-2735 gives them the same check and widens its touches.
- `scripts/agent-container.sh` keeps its own rewrite and gets no check: it is a script the operator runs by hand.
- Nothing can be measured on a native engine here. The test uses a fake docker, and the issue's close reason will say that.
- No `after` on another story, and no ADR planned. The story's agent may add an ADR if it judges the refusal a new decision.

To redirect: if you would rather change the default bind than refuse, or probe from a real container instead of reading `docker info`, say so here and I will replan T-2735.

### 2026-11-09T17:06:01Z orchestrator
**Approved as planned.** I am finalizing S-0619 now.

- **Remedy fits.** I-0254 names no remedy, and its only instance (S-0508) proposes this one: `flai serve` refuses a container start when the bind cannot be reached from the bridge. Keeping the `127.0.0.1` default is right, because S-0506 measured that Docker Desktop under WSL2 cannot reach a `docker0` bind.
- **Tasks fit.** T-2735 builds the preflight and the test that reproduces I-0254 (criterion 1). T-2736 updates the documents and closes I-0254 (criterion 2). T-2736 waits for T-2735 because the documents quote its refusal.
- **Touches fit.** All 12 are files, and every one exists. `containerAllowed` in `container.go` is already called from `agents.go`, so the check has a natural home.
- **Figures are flai's.** The 40m forecast raises flai's 22m and says why. 37.50 USD a week comes from the inputs flai set from I-0254.
- **Assumptions accepted.** When `docker info` cannot answer, flai warns and does not refuse. The engine kind is read once per process. `scripts/agent-container.sh` gets no check.

**One note for the story's agent.** A start `flai serve` refuses is a tool-enforced rule. Under `design/conventions/decisions.md` § Rules, that needs an ADR, so write one with `flai adr new --commit`. `design/adrs` is in `claims.shared`, so adding it to the touches holds nothing.

Source: design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md

Source: design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md

### 2026-11-09T17:06:19Z orchestrator
Resolved: Plan approved under plan_backlog_stories: the remedy is the one I-0254's instance proposes, the 2 tasks and 12 file touches fit both criteria, and the figures are flai's.
