---
title: Active streams
updated: 2026-11-10T09:39:38Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0408](S-0408.md) | An OpenCode agent's protected write is asked of the story's owner through flai's hold-and-ask from the plugin's permission.asked, or refused by deny rules under --auto, as an ADR settles | cancelled | system-flow | 2026-11-07T22:14:13Z |
| [S-0409](S-0409.md) | The explorer, the verifier, and the strategic agents run on OpenCode as agents flai generates from .claude/agents, each on the model the story's roles name and read-only where the definition is | cancelled | system-flow | 2026-11-07T22:04:57Z |
| [S-0410](S-0410.md) | A story on OpenCode has its usage read from OpenCode's JSON events, priced by the harness or its gateway's spend log, and reaches its provider through an entry derived from the host's | cancelled | agent-S-0410 | 2026-11-07T21:59:02Z |
| [S-0668](S-0668.md) | A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps | in-progress | agent-S-0668 | 2026-11-10T09:36:53Z |
| [S-0670](S-0670.md) | scripts/agent-container.sh always runs flai-agent-claude-code:local and ignores the manifest's agent.image.layer, so a close-out run by hand in it has none of the project's test tools | review | agent-S-0670 | 2026-11-10T09:35:34Z |
| [S-0671](S-0671.md) | scripts/close-out.sh takes no --fresh, so at a commit the verify record holds a pass for it runs no tier, and a second close-out in another place cannot check the same commit again | in-progress | agent-S-0671 | 2026-11-10T09:39:31Z |
| [S-0682](S-0682.md) | scripts/agent-container.sh can run only a story's own agent, under the story's own name, so a trial run of the image acts on the live story | in-progress | agent-S-0682 | 2026-11-10T09:26:22Z |
| [S-0684](S-0684.md) | flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest | in-progress | agent-S-0684 | 2026-11-10T09:32:25Z |
| [S-0687](S-0687.md) | A strategic agent's container runs flai with no host configuration, so the orchestrator's flai accept and flai verify inside it see no host actions and lack the test tiers' tools | in-progress | agent-S-0687 | 2026-11-10T09:38:30Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 210 | 1165.3067 USD | 146972 | 2026-11-10T09:38:04Z |
| [orchestrator](orchestrator.md) | 2751 | 2088.0096 USD | 1541730 | 2026-11-10T09:38:54Z |
| [analyzer](analyzer.md) | 2 | 20.4997 USD | 1247 | 2026-11-07T08:18:10Z |
| [curator](curator.md) | 123 | 196.8129 USD | 22780 | 2026-11-10T09:13:17Z |
