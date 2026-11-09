---
id: S-0619
type: story
nature: improvement
title: An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal
status: done
owner: alex
created: 2026-11-07T23:39:38Z
updated: 2026-11-09T17:28:12Z
transitions:
  - to: ready
    at: 2026-11-09T17:06:54Z
    by: orchestrator
  - to: in-progress
    at: 2026-11-09T17:07:02Z
    by: agent-S-0619
  - to: review
    at: 2026-11-09T17:26:27Z
    by: agent-S-0619
  - to: done
    at: 2026-11-09T17:28:12Z
    by: orchestrator
tags: []
topics: [mcp, docker]
touches: [flai/internal/serve/container.go, flai/internal/serve/container_test.go, flai/internal/serve/agents.go, docs/operators/index.md, docs/operators/settings.md, docs/operators/runbooks/install.md, docs/operators/runbooks/agent-container.md, docs/users/flai.md, design/system/flai-cli.md, design/tech/docker.md, design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md, design/issues/summary.md, design/adrs/0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md, design/adrs/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1175
  priced_by: harness
  turns:
    - day: 2026-11-09
      ceremony: 1
      hand_edits: 1
      work: 50
  models:
    - model: claude-opus-5-5
      input: 132
      output: 28353
      cache_read: 16807807
      cache_write: 329888
      cost: 6.5031
  strategic:
    - kind: orchestrator
      seconds: 219
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 61
          output: 530
          cache_read: 2261685
          cache_write: 252845
          cost: 0.6636
        - model: claude-sonnet-5-5
          input: 6
          output: 51
          cache_read: 40925
          cache_write: 37362
          cost: 0.0581
cost_of_delay:
  inputs:
    time_lost_per_cycle: 15m
    by: flai
    at: 2026-11-07T23:39:38Z
  value: 37.5
  by: planner-S-0619
  at: 2026-11-09T17:04:19Z
forecast:
  duration: 40m
  delivery: 2026-11-10T12:00:00Z
  basis: "flai forecast's 22m (100 s per unit of size over 138 done large-band improvement stories, size 13) raised to 40m because the story first settles a remedy I-0254 does not name and its tests need a fake docker info; delivered a day out since it is a draft until the operator finalizes it."
  by: planner-S-0619
  at: 2026-11-09T17:04:19Z
finalized:
  by: orchestrator
  at: 2026-11-09T17:06:25Z
---
# S-0619 An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal

## Goal

This story remediates [I-0254](../../../design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md), "An agent container on a native Linux Docker engine cannot reach flai's HTTP MCP server on its default 127.0.0.1 bind through host.docker.internal". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0254 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0254 is closed with `flai issue close I-0254 --reason` saying what fixed it

## Tasks
- T-2735 flai serve refuses a container start whose loopback MCP bind a native engine's bridge cannot reach, naming flai serve mcp bind and the bridge's gateway
- T-2736 The operators guide, the runbooks, the user guide, and the CLI and Docker design describe the refusal, and I-0254 is closed

## Notes

Cost of delay inputs set by flai from I-0254. time_lost_per_cycle 15m: 15m per occurrence × 1 occurrence ÷ 1 cycle of 168h (first reported 2026-11-02T07:30:21Z, 5.7 days before this story; under one cycle counts as one).

### Planning

The remedy planned is the one I-0254's instance proposes. Before it starts a container on the bridge, `flai serve` refuses the start when three things hold: the MCP URL's host is loopback (`mcp.bind` unset), the network is not `host`, and `docker info` says the engine is not Docker Desktop. The refusal names `flai serve mcp bind <gateway>`. The default bind stays `127.0.0.1`. Binding the bridge's gateway on its own was not chosen: S-0506 measured that Docker Desktop under WSL2 cannot reach a bind to WSL's `docker0` address, so that default would break the one engine measured here.

Touches, all files, no folder:

| Touch | Source |
|-------|--------|
| `flai/internal/serve/container.go` | layout: `containerAllowed`, `containerMCP`, and `containerHost` live here (explorer) |
| `flai/internal/serve/container_test.go` | layout: `TestAContainerStartIsRefusedWhatItLacks` and the fake docker |
| `flai/internal/serve/agents.go` | layout: `startOf` calls `containerAllowed`, then `containerMCP` |
| `docs/operators/index.md` | design: § "MCP over HTTP" says when to set `mcp.bind` |
| `docs/operators/settings.md` | design: the `mcp.bind` row |
| `docs/operators/runbooks/install.md` | design: tells a native engine's operator to set the bind |
| `docs/operators/runbooks/agent-container.md` | design: says a Docker Engine host was not tried |
| `docs/users/flai.md` | co-change: 27% of the commits that change the other touches; it shows `flai serve mcp bind 172.17.0.1` |
| `design/system/flai-cli.md` | design: says a native Linux engine needs `mcp.bind` |
| `design/tech/docker.md` | design: the "Agent harness" row lists the docker calls |
| `design/issues/I-0254-…` and `design/issues/summary.md` | criterion 2: `flai issue close` writes both |

`flai touches suggest` listed nothing at first, since the story declared no touches. Run again from the touches above, it lists `docs/users/flai-reference.md` (13%) and wider design documents, such as `design/system/strategic-agents.md`. They are left out because the story changes no command or flag, and the reference is generated from those. The story has no `after`: no other story writes anything it reads.

Figures:

- Forecast 40m, delivery 2026-11-10T12:00:00Z. `flai forecast` gives 22m: 100 s per unit of size, over 138 done improvement stories in the large band, times size 13. It is raised because the story first settles a remedy the issue does not name, and its tests need a fake `docker info` and a fake `docker network inspect`; the size count sees neither. Delivery is a day out because the story is a draft until the operator finalizes it.
- Cost of delay 37.50 USD a week, as `flai cod` gives it: 15m lost per 168h cycle at 150 USD an hour. It stands, since the inputs are the operator's, set from I-0254. The issue's one occurrence was not measured: no native engine runs on this host.

### Accepted by the orchestrator

- Verified: 799ab41ec311c0ba2bfa36263ccc858ef72405b6
- At: 2026-11-09T17:28:12Z

Verdict: pass
- 1: flai/internal/serve/container.go, flai/internal/serve/agents.go, flai/internal/serve/container_test.go, design/adrs/0269-flai-serve-refuses-a-container-start-whose-http-mcp-server-listens-on-loopback.md
- 2: design/issues/I-0254-an-agent-container-on-a-native-linux-docker-engine-cannot-reach-flai-s-http-mcp.md, design/issues/summary.md
