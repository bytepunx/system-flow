---
id: S-0177
type: story
nature: remediation
title: A story begun on another host can be given an agent on this one
status: in-progress
owner: alex
created: 2026-10-01T07:30:22Z
updated: 2026-10-01T10:08:19Z
transitions:
  - to: ready
    at: 2026-10-01T07:39:53Z
    by: alex
  - to: in-progress
    at: 2026-10-01T10:00:28Z
    by: agent-S-0177
tags: [flai, dashboard]
touches: [flai/internal/serve, flai/internal/hostapi, flai/cmd, flai/internal/mcpserver, flai/internal/harness, flai/internal/storygit, flai/internal/workitem, flaiover/src, design/adrs/0064-a-story-in-ready-or-in-progress-with-no-agent-run-on-this-host-is-started-here.md, design/adrs/README.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/agent-narrative.md, docs/users/flai.md, docs/users/flai-reference.md, docs/users/flaiover.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 678
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 208
      output: 1156
      cache_read: 9555876
      cache_write: 381313
      cost: 3.9678
---
# S-0177 A story begun on another host can be given an agent on this one

## Goal

A story an agent began on one host cannot be picked up on another. Seen 2026-10-01 with S-0173: an agent on another machine moved it to in-progress, opened its narrative, and wrote its tasks, which reached this clone through git; its story branch and worktree never did. On this host:

- The story page shows no agent block, because this host's `flai serve` has no run for it in `serve/agents.json`, and nothing says the story was started elsewhere.
- `flai serve agent restart` (the Retry button) refuses it, because flai serve here "has started no agent for it"; `flai serve agent start` (Start agent) refuses it, because it is not in `ready`. The launcher starts only stories in `ready`. The only way to get an agent is to move the story back to `ready`, which writes a false transition into its history and metrics.
- The designer answered the agent's question (TH-0041) here, and nothing resumed anything: the run waiting on it is on the other host.

A story in progress or in review with no agent run on this host should be startable here, and the agent should be told what state the story was left in.

## Acceptance criteria
- [ ] `flai serve agent restart` (and `agent.restart`) starts an agent for a story in ready or in progress that this host has no run for, instead of refusing it; it still refuses while an agent of this host's runs or waits for it
- [ ] The prompt for such an agent says the story was begun elsewhere (from its transitions and narrative: the agent that moved it and when) and to reconcile: use `story/S-nnnn` if it exists locally or on the remote, create it from the main branch otherwise, read the narrative and tasks, and go on from what is committed rather than starting over. If the story has threads answered since it was begun, the agent is told to read them
- [ ] `flai stream open` (or what the agent runs) recreates the worktree for a story whose narrative exists but whose branch and worktree do not, fetching `story/S-nnnn` from the remote when it is there
- [ ] The story page and the board card show a story in progress with no agent on this host as such ("begun by agent-S-nnnn on another host; no agent here"), with Start agent offered, rather than showing no agent block
- [ ] The MCP server can start or restart a story's agent, gated by the same agent host action as the dashboard and CLI, so an operator's own agent can do it
- [ ] Tests cover a story in progress with no run on this host, with and without its branch on the remote, and the refusal while this host's agent runs
- [ ] The design (`design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`) and the user guide describe picking up a story on another host

## Tasks
- T-0662 An ADR records how a story begun on another host gets an agent on this one
- T-0663 flai stream open reopens a stream whose worktree is gone, from the local or remote story branch, and the narrative records its host
- T-0664 flai serve agent restart starts a story with no run on this host, and its agent is told the story was begun elsewhere
- T-0665 agent.status names a story in progress with no agent here, and the story page and board card show it with Start agent
- T-0666 The MCP server starts or restarts a story's agent, gated by the agent host action

## Notes

- Related: an agent waiting on a thread on one host is not resumed when the thread is answered on another. Starting an agent here covers the case where the first host is gone; whether flai should say where a story's agent runs (a host name in the narrative's front matter, say) is worth deciding here.
- Pushing story branches as an agent works would make the work itself portable; out of scope unless it is the cheap way to meet the reconcile criterion.
