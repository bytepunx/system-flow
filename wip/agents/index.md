---
title: Active streams
updated: 2026-11-08T02:58:20Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0408](S-0408.md) | An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles | cancelled | system-flow | 2026-11-07T22:14:13Z |
| [S-0409](S-0409.md) | The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is | cancelled | system-flow | 2026-11-07T22:04:57Z |
| [S-0410](S-0410.md) | A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's | cancelled | agent-S-0410 | 2026-11-07T21:59:02Z |
| [S-0591](S-0591.md) | TestGatewayLiteLLM fails its setup check: the template's settings run flai guard 6 times, and it wants 4 | review | system-flow | 2026-11-08T00:00:12Z |
| [S-0598](S-0598.md) | A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s | review | agent-S-0598 | 2026-11-08T01:11:49Z |
| [S-0599](S-0599.md) | flai serve's replanner writes every forecast it moves in one pass with one check, so an acceptance's replan lands in seconds, not 80 s | review | agent-S-0599 | 2026-11-08T02:06:30Z |
| [S-0601](S-0601.md) | flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s | review | agent-S-0601 | 2026-11-08T02:12:06Z |
| [S-0603](S-0603.md) | agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes under 200 ms, not 0.8 to 2.4 s | in-progress | agent-S-0603 | 2026-11-08T02:38:37Z |
| [S-0604](S-0604.md) | board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing | in-progress | agent-S-0604 | 2026-11-08T02:43:58Z |
| [S-0606](S-0606.md) | item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms | in-progress | agent-S-0606 | 2026-11-08T02:12:40Z |
| [S-0609](S-0609.md) | flai accept and flai task done lint markdown once per set with the installed markdownlint, not five times through npx, saving about 4 s an acceptance | in-progress | agent-S-0609 | 2026-11-08T02:56:48Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 165 | 1088.1764 USD | 141678 | 2026-11-08T00:11:17Z |
| [orchestrator](orchestrator.md) | 2294 | 2031.4583 USD | 1524452 | 2026-11-08T02:57:58Z |
| [analyzer](analyzer.md) | 2 | 20.4997 USD | 1247 | 2026-11-07T08:18:10Z |
| [curator](curator.md) | 98 | 157.6368 USD | 15682 | 2026-11-08T01:51:17Z |
