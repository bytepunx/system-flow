---
id: S-0115
type: story
nature: feature
title: "flai serve agent start <story> starts the story's agent on demand, and the story page has a button for it"
status: backlog
owner: alex
created: 2026-09-24T08:30:36Z
updated: 2026-09-24T08:30:36Z
transitions: []
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

- [ ] `flai serve agent start <story>` starts the agent the story names (its harness, model, and options, or the host's command), the way flai serve does when a story enters ready, and records the run in `serve/agents.json` so that the serving flai tracks it from then on: the dot on the card, the outcome when it ends, the restart on an answer
- [ ] it starts whether or not the story was ready before serve started, and whether or not anyone is attending; it still refuses while the `agent` host action is off for the project, when the story is not in ready, when the story already has an agent running, or when nothing can start it (no harness and no command), and each refusal says why
- [ ] the story's page has a **Start agent** button under the same conditions, through a hostapi write gated by the `agent` action, journalled like every host action
- [ ] tried live: a story in ready, stale to the running serve, gets its agent from the command and from the button

## Tasks

## Notes

Asked by the operator on 2026-09-24 after S-0112's research: a story that was in ready when flai serve started is never started by that serve, and a story cannot be sent back to ready (ready allows in-progress and cancelled only), so there is no way to make flai pick it up short of cancelling and recreating it. S-0112 fixes the stale rule; this story gives the operator a direct way in either case.

The launcher runs inside flai serve, and a command in another process cannot reach it. Two ways in: the command runs the launcher's own start path in its own process against the shared `serve/agents.json`, and serve's next look settles the run when it ends (`settleOrphans` already handles a run no launcher waits for); or the command asks flai host, which asks serve. The first is smaller and needs no new API. `agent.status` should show the run either way.
