---
title: Active streams
updated: 2026-11-09T08:04:55Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0408](S-0408.md) | An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles | cancelled | system-flow | 2026-11-07T22:14:13Z |
| [S-0409](S-0409.md) | The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is | cancelled | system-flow | 2026-11-07T22:04:57Z |
| [S-0410](S-0410.md) | A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's | cancelled | agent-S-0410 | 2026-11-07T21:59:02Z |
| [S-0598](S-0598.md) | A checked write runs only the rules an item write can affect, so creating or editing an item takes about 2 s, not 7 to 18 s | review | agent-S-0598 | 2026-11-09T08:03:53Z |
| [S-0601](S-0601.md) | flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s | review | agent-S-0601 | 2026-11-08T02:12:06Z |
| [S-0603](S-0603.md) | agent.status keeps agents.json parsed and reads the conversations once a call, so it stops queueing on agentsMu and takes about 250 ms, not 0.8 to 2.4 s | review | agent-S-0603 | 2026-11-09T07:39:00Z |
| [S-0604](S-0604.md) | board.get takes one release lock per repository, runs show-ref only when a ref moved, and reads no threads with nothing in review, so 13 boards at once stop queueing | review | agent-S-0604 | 2026-11-09T08:01:25Z |
| [S-0606](S-0606.md) | item.get answers from a kept item list and the file watcher stops statting every archived file every 300 ms | in-progress | agent-S-0606 | 2026-11-09T08:01:57Z |
| [S-0610](S-0610.md) | The close-out runs independent test tiers at once and smoke skips what the close-out already ran, saving about 40 to 150 s a close-out | in-progress | agent-S-0610 | 2026-11-09T07:52:06Z |
| [S-0612](S-0612.md) | flai guard refuses a story's task sub-agent flai hostapi, so a measurement task cannot use the documented hostapi --timing | in-progress | agent-S-0612 | 2026-11-09T08:03:52Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 166 | 1090.8430 USD | 142062 | 2026-11-09T08:00:09Z |
| [orchestrator](orchestrator.md) | 2342 | 2037.5400 USD | 1526486 | 2026-11-09T08:03:51Z |
| [analyzer](analyzer.md) | 2 | 20.4997 USD | 1247 | 2026-11-07T08:18:10Z |
| [curator](curator.md) | 102 | 165.2491 USD | 16869 | 2026-11-09T08:01:36Z |
