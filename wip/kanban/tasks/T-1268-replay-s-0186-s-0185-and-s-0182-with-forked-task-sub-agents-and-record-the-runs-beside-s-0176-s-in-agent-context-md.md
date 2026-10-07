---
id: T-1268
type: task
nature: experiment
title: Replay S-0186, S-0185, and S-0182 with forked task sub-agents and record the runs beside S-0176's in agent-context.md
status: backlog
parent: S-0241
owner: alex
created: 2026-10-07T23:12:53Z
updated: 2026-10-07T23:12:53Z
transitions: []
stream: S-0241
tags: [template]
touches: [design/system/agent-context.md]
after: [T-1267]
---
# T-1268 Replay S-0186, S-0185, and S-0182 with forked task sub-agents and record the runs beside S-0176's in agent-context.md

## Work

Repeat the measurement in `design/system/agent-context.md` § Tasks in parallel › Measured against runs without the plan, as TH-0059 settled it, with one change: each task sub-agent is started with `subagent_type: fork` instead of `general-purpose`. Replay S-0186, S-0185, and S-0182 in a scratch clone with no remote, each from the `main` it began at (`e33b80d`, `723efff`, `b8caca4`) with the story put back to `ready`. Lay over it the current conventions, `.claude/agents/`, and a flai built from the story branch for `scripts/flai.sh`, the MCP server, and the guard. Start each with the command `flai serve` builds, and tell the story's agent to fork its task sub-agents. Answer S-0182's question with the original answer's words.

Measure each run as § Measured does: minutes from `in-progress` to `review` less designer waits, cost from the newest `result`, and model calls and cache reads split between the agent and its sub-agents, deduplicated by message ID. Have one fresh verifier review each replay's change blind against the earlier replay's, labelled X and Y.

Add the three forked runs to the section's table, or a table beside it, and say what they show against the earlier replays: above all, whether the sub-agents' cache reads and the runs' cost fall.

Waits for T-1267: it cannot run until a headless session is shown to offer forks, and both write the same section of `design/system/agent-context.md`.

## Done when

- The section records minutes, cost, model calls, cache reads, and review defects for each of the three forked replays, beside the earlier replays of the same stories.
- It says whether the forked sub-agents' cache reads and each run's cost fell against the general-purpose replays, with the figures.

## Notes

The replays run in a scratch clone outside the repository, so this task changes no file but `design/system/agent-context.md`. If a fork needs a flag or setting in the `claude-code` adapter to be offered, that is a change under `flai/internal/harness`: widen this task's touches and say so in the narrative's `## Decisions`.
