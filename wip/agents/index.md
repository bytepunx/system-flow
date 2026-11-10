---
title: Active streams
updated: 2026-11-10T10:55:00Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0665](S-0665.md) | flai serve records no image, container name, or mounts for a container harness run, so a run cannot be shown afterwards to have used the project's layer image | in-progress | agent-S-0665 | 2026-11-10T10:41:47Z |
| [S-0668](S-0668.md) | A story's agent in a container cannot write its narrative's Decisions, Context, or Open questions, since wip is read-only there and flai stream state writes only Current state and Next steps | review | agent-S-0668 | 2026-11-10T10:50:44Z |
| [S-0684](S-0684.md) | flai/internal/mcpserver's tests still take 136 to 162 s under -race, since the server reads the session's story, role, session, reasons, and started-by from the environment and 95 tests must run serially or in env groups before the rest | review | agent-S-0684 | 2026-11-10T10:52:51Z |
| [S-0708](S-0708.md) | TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name | in-progress | agent-S-0708 | 2026-11-10T10:52:56Z |
| [S-0713](S-0713.md) | make agent-images fails at the first base image when Docker Desktop's credential helper desktop.exe cannot start from WSL | in-progress | agent-S-0713 | 2026-11-10T10:54:16Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 226 | 1256.2868 USD | 148461 | 2026-11-10T10:39:27Z |
| [orchestrator](orchestrator.md) | 2789 | 2092.7585 USD | 1543488 | 2026-11-10T10:31:01Z |
| [analyzer](analyzer.md) | 2 | 20.4997 USD | 1247 | 2026-11-07T08:18:10Z |
| [curator](curator.md) | 129 | 204.3154 USD | 23943 | 2026-11-10T10:24:03Z |
