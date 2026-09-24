---
id: S-0115
type: story
nature: feature
title: "flai serve agent start <story> starts the story's agent on demand, and the story page has a button for it"
status: done
owner: alex
created: 2026-09-24T08:30:36Z
updated: 2026-09-24T20:31:35Z
transitions:
  - to: ready
    at: 2026-09-24T08:32:38Z
    by: alex
  - to: in-progress
    at: 2026-09-24T09:06:12Z
    by: agent-S-0115
  - to: review
    at: 2026-09-24T09:32:21Z
    by: agent-S-0115
  - to: done
    at: 2026-09-24T20:31:35Z
    by: alex
tags: [cli, dashboard]
touches: [flai/internal/serve, flai/cmd, flaiover/src]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0115 flai serve agent start <story> starts the story's agent on demand, and the story page has a button for it

## Goal

The operator can have flai start a story's agent now, from a shell on the host or from the story's page, without waiting for the launcher's own rules to fire.

## Acceptance criteria

- [x] `flai serve agent start <story>` starts the agent the story names (its harness, model, and options, or the host's command), the way flai serve does when a story enters ready, and records the run in `serve/agents.json` so that the serving flai tracks it from then on: the dot on the card, the outcome when it ends, the restart on an answer
- [x] it starts whether or not the story was ready before serve started, and whether or not anyone is attending; it still refuses while the `agent` host action is off for the project, when the story is not in ready, when the story already has an agent running, or when nothing can start it (no harness and no command), and each refusal says why
- [x] the story's page has a **Start agent** button under the same conditions, through a hostapi write gated by the `agent` action, journalled like every host action
- [x] tried live: a story in ready, stale to the running serve, gets its agent from the command and from the button

## Tasks
- T-0420 flai serve agent start <story> starts a ready story's agent now, on the start path S-0116 builds for restart
- T-0421 The story page has a Start agent button, through a hostapi write gated by the agent action
- T-0422 Record start in the design and the operator and user docs
- T-0423 Tried live: a ready story stale to the running serve gets its agent from the command and from the button

## Notes

Asked by the operator on 2026-09-24 after S-0112's research: a story that was in ready when flai serve started is never started by that serve, and a story cannot be sent back to ready (ready allows in-progress and cancelled only), so there is no way to make flai pick it up short of cancelling and recreating it. S-0112 fixes the stale rule; this story gives the operator a direct way in either case.

The launcher runs inside flai serve, and a command in another process cannot reach it. Two ways in: the command runs the launcher's own start path in its own process against the shared `serve/agents.json`, and serve's next look settles the run when it ends (`settleOrphans` already handles a run no launcher waits for); or the command asks flai host, which asks serve. The first is smaller and needs no new API. `agent.status` should show the run either way.

Built on S-0116's start path (TH-0010): `serve.Start` in `flai/internal/serve/start.go`, over S-0116's `serve.StartNow`, which now hands the agent to the serving flai instead of waiting for it. The command and the hostapi write `agent.start` share restart's code. Start refuses while the action is off, when the story is not in ready, while its agent runs or waits for an answer, and when nothing can start it. A full in-progress limit does not refuse it: it starts past the limit with a warning, as `flai move` does, because criterion 2 does not list the limit. The Start agent button shows only for a ready story with no run; one whose run failed has S-0116's Restart agent.

Tried live on 2026-09-24 (criterion 4), with a scratch project, config, and serve folder under `/tmp/s0115-live`. It used a scratch `flai serve` and a flaiover dev server on port 5199, both built from story/S-0115, and a harmless agent command that prints and sleeps 45 s. The in-progress limit was 1, with S-0003 in progress. S-0001 and S-0002 were moved to ready before serve started, so the launcher held both ("the in-progress limit leaves no room").
- Refusals: `flai serve agent start S-0003` (in progress) said `rule: S-0003 is in in-progress; only a story in ready is started; flai serve agent restart starts ...`. A second start of S-0001 while its agent ran said `rule: S-0001's agent is running (pid ..., started ...)`.
- `flai serve agent start S-1` started the command for S-0001 and warned `agent started past the in-progress limit`. The run was in `serve/agents.json` and the journal. When the process ended, serve settled it at a look as `failed` (story left in ready). The launcher then held S-0001 as having had its agent, and its page showed the red dot with Restart agent and no Start agent.
- Headless Chromium logged in to the dev dashboard and opened S-0002. The panel read "flai serve has started no agent for this story: the in-progress limit leaves no room for S-0002", with Start agent. Pressing it gave "agent working (command) agent-S-0002". The journal had `agent.start` for alex and `serve.agent` for the start. Serve settled the run at a look once the process ended (`failed`, story in ready).
- A handed-over run that ends asking is started again on the answer. That was not tried live; `TestAReadyStorysAgentIsStartedOnTheOperatorsWord` pins it.
- The first scratch serve inherited `FLAI_HOST_URL` and `FLAI_HOST_TOKEN` from this agent's environment. It told the operator's flai host which MCP servers to keep, which stopped the operator's six for about two minutes. It was stopped by PID, the six came back, and the trial went on with those variables unset. Recorded as I-0044.
