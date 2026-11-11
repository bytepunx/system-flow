---
title: Active streams
updated: 2026-11-11T01:39:51Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0667](S-0667.md) | A docker run from a story's agent's shell mounts the main checkout as an empty folder, so a hand call in the agent image cannot read the project | review | agent-S-0667 | 2026-11-10T19:29:32Z |
| [S-0675](S-0675.md) | The merge queue builds one-story batches and runs three verifications per story landed | review | agent-S-0675 | 2026-11-10T20:20:12Z |
| [S-0677](S-0677.md) | A whole flai check takes 3.3 to 4.9 s on a loaded host after S-0601, not under 2 s: the item list, workItems alone, and the digest rule are its critical path | review | system-flow | 2026-11-10T20:15:07Z |
| [S-0679](S-0679.md) | flai check finds `item.heading` outside the story at close-out | review | agent-S-0679 | 2026-11-10T20:01:14Z |
| [S-0688](S-0688.md) | A strategic agent's container can write the main checkout's .git/config, which can name commands git runs on the host | review | agent-S-0688 | 2026-11-10T19:50:12Z |
| [S-0708](S-0708.md) | TestGatewayStreamIsJudged and TestGatewayMissingNamesWhatIsMissing run in no tier, because go-test and integration skip ^TestGateway and the gateway tier runs only the two paid checks by name | review | agent-S-0708 | 2026-11-10T19:36:00Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 230 | 1264.0888 USD | 149634 | 2026-11-10T11:08:32Z |
| [orchestrator](orchestrator.md) | 2873 | 2102.1177 USD | 1545864 | 2026-11-10T20:21:43Z |
| [analyzer](analyzer.md) | 2 | 20.4997 USD | 1247 | 2026-11-07T08:18:10Z |
| [curator](curator.md) | 133 | 210.0617 USD | 24818 | 2026-11-10T19:28:49Z |
