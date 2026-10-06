---
id: S-0222
type: story
nature: feature
title: "The orchestrator publishes by the release policy: judgement, theme, or cost of delay threshold"
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:17Z
updated: 2026-10-06T02:57:19Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:24Z
    by: alex
tags: [flai]
topics: [orchestration, release]
touches: [flai/internal/harness, flai/internal/hostapi, flai/internal/release, design/system/strategic-agents.md, flai/internal/guard, flai/internal/manifest, design/system/project-manifest.md, docs/operators/settings.md, flai/internal/mcpserver, flai/cmd/release.go, ".claude/agents/orchestrator.md", template/root/.claude/agents/orchestrator.md, design/adrs, design/system/flai-cli.md, docs/users/flai.md, docs/operators/index.md]
after: [S-0218, S-0174]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 76.53
  by: planner-S-0222
  at: 2026-10-05T04:48:06Z
forecast:
  duration: 1h15m
  delivery: 2026-10-06T06:07:00Z
  basis: "Its own forecast of 1h15m; 4th in the pull order with an in-progress limit of 3, behind S-0219, S-0220 and S-0221."
  by: flai
  at: 2026-10-06T02:57:19Z
---
# S-0222 The orchestrator publishes by the release policy: judgement, theme, or cost of delay threshold

## Goal

With `publish` on, the orchestrator cuts releases when the project's release policy says so: a cost of delay or count threshold, a theme (an epic or tag whose stories are all accepted), or its own judgement when the policy is `judgement`.

## Acceptance criteria
- [ ] After each acceptance it runs `flai release --evaluate`; when the policy is met it publishes through `publish.run` (which requires the `push` host action too) and logs the release with the figures and the items bundled
- [ ] Under `judgement`, it publishes when it judges the unreleased work coherent and complete, and logs the reasoning; it never publishes a batch with a story whose epic is not in review or done when `orchestration.release.whole_epics` is set
- [ ] When publishing is refused (S-0174's tag check, a moved remote) it logs the refusal and opens a thread to the operator
- [ ] Tests cover each policy met and not met, and the refusal

## Tasks
- T-0885 orchestration.release.whole_epics holds back a batch with a story whose epic is not in review or done, and flai release --evaluate names it
- T-0886 flai guard lets the orchestrator publish only through release_publish, and only while the publish permission is on
- T-0891 The release_publish MCP tool publishes through publish.run when the release policy is met, and returns the release or the refusal
- T-0895 The orchestrator's prompt evaluates the release policy after each acceptance, publishes when it is met, and logs the release or the refusal with a thread
- T-0897 An ADR, the design, and the guides describe the orchestrator publishing by the release policy

## Notes

### Planning

Planned by planner-S-0222 on 2026-10-05. The plan is summarised on the story's plan thread.

Touches:

- Declared, kept: `flai/internal/harness`, `flai/internal/hostapi`, `flai/internal/release`, `flai/internal/guard`, `flai/internal/manifest`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/settings.md`.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (34%) and `docs/users/flai.md` (32%), where the new MCP tool is described. Also `docs/operators/index.md` (13%), whose push host action section gains the `publish` permission's need for it.
- Layout: `flai/internal/mcpserver`, for the `release_publish` tool beside S-0217's `release_evaluate` in `orchestrate.go`, registered in `folder.go`. Also `flai/cmd/release.go`, where `flai release --evaluate` prints the stories `whole_epics` holds back.
- Design: `.claude/agents/orchestrator.md` and `template/root/.claude/agents/orchestrator.md`, the orchestrator's definition S-0218 adds, kept the same as the planner's two copies are. Also `design/adrs`, for an ADR that refines ADR-0067 (agents publish only when the operator asks).
- Left out from co-change: `flaiover/src/lib/server/agent.ts`, `flai/cmd/serve_actions.go`, and `flai/internal/serve/agents.go`. Publishing adds no dashboard view and no host action; the orchestrator's run is S-0218's.

Topics `orchestration` and `release` were added: the tasks reach both.

Forecast: 1h15m, delivery 2026-10-05T15:16Z. `flai forecast` gave 45m: 132 s per unit of size over 14 done feature stories on claude-opus-5-5 in the large band, times size 20. It was raised by 30m for two reasons. The five tasks build on code that S-0217 and S-0218 have yet to write, which the agent must read first. And `release_publish` needs git fixtures with a bare remote for both refusals. S-0221, of like shape, is forecast at 1h30m. The delivery is flai's 14:31Z moved by the 30m added at its cycle factor of 1.5. It is 8th in the pull order and starts after S-0218.

Cost of delay: 76.53 USD a week, from `flai cod`. This is S-0222's share of E-0016's 1500 USD a week (the operator's inputs on the epic), 1h15m of the 24h30m forecast over the epic's 17 open stories without inputs. It replaces the 60.98 that planner-E-0016 set on 2026-10-04, because the epic's open forecasts have since shrunk. The story has no inputs of its own.
