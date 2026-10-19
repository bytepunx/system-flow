---
title: Active streams
updated: 2026-10-19T23:41:20Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0401](S-0401.md) | A batch settles a rebase that stops inside the manifest's shared paths by the regenerate-and-merge steps and never sends a story back for it, and tells the story's agent of a conflict outside them by message | review | agent-S-0401 | 2026-10-19T22:35:28Z |
| [S-0403](S-0403.md) | With accept_reviews on and the queue running, the orchestrator enqueues each story in review it has verified instead of accepting it, never lands or bisects by hand, and opens no gate thread on a story the queue will verify | in-progress | agent-S-0403 | 2026-10-19T23:41:28Z |
| [S-0404](S-0404.md) | metrics.md defines the queue wait per story, the batch size, and the bisections per week by an ADR, and flai stats and the dashboard's charts report them | in-progress | agent-S-0404 | 2026-10-19T23:41:29Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 79 | 751.0673 USD | 31390 | 2026-10-16T21:51:32Z |
| [orchestrator](orchestrator.md) | 655 | 943.9724 USD | 914349 | 2026-10-19T22:36:49Z |
| [analyzer](analyzer.md) | 1 | 5.0040 USD | 454 | 2026-10-15T22:49:08Z |
