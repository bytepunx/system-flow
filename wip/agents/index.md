---
title: Active streams
updated: 2026-11-02T18:29:20Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0509](S-0509.md) | The planner, the orchestrator, the analyzer, and the curator run in a container when their harness is container, with the main checkout mounted as an ADR settles, pilot passing their role, and their activity logged as on the host | in-progress | agent-S-0509 | 2026-11-02T18:20:57Z |
| [S-0510](S-0510.md) | A project's manifest names an image layer for its test tiers' tools, flai serve runs it, flai check warns when it is not built, and this repository's layer runs the close-out inside the container | in-progress | agent-S-0510 | 2026-11-02T18:18:44Z |
| [S-0517](S-0517.md) | A story's worktree mounted alone into a container is not a git repository, since its .git file names the main checkout's .git/worktrees folder | in-progress | agent-S-0517 | 2026-11-02T18:28:11Z |
| [S-0518](S-0518.md) | flai guard does not run as the hook inside the agent container, and the hook lets every call through when the guard cannot start | in-progress | agent-S-0518 | 2026-11-02T04:24:48Z |
| [S-0553](S-0553.md) | flai generates the planner's, the orchestrator's, the analyzer's, and the curator's harness definitions from their role files alone, flai check holds a .claude/agents copy to them, and a session started by hand reads the same definition | review | agent-S-0553 | 2026-11-02T02:36:19Z |
| [S-0566](S-0566.md) | Operators should not need to intervene during review acceptance | in-progress | agent-S-0566 | 2026-11-02T18:26:18Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 122 | 956.9301 USD | 130434 | 2026-11-02T18:13:32Z |
| [orchestrator](orchestrator.md) | 1732 | 1961.7717 USD | 1505820 | 2026-11-02T18:20:07Z |
| [analyzer](analyzer.md) | 1 | 5.0040 USD | 454 | 2026-10-15T22:49:08Z |
| [curator](curator.md) | 50 | 87.5112 USD | 8461 | 2026-11-02T18:15:32Z |
