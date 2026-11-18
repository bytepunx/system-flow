---
title: Active streams
updated: 2026-11-18T09:20:10Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0807](S-0807.md) | A story's move needs the suites its transition's verify lists to have passed at its head, and the close-out runs the suites its next transitions need | in-progress | agent-S-0807 | 2026-11-18T09:09:15Z |
| [S-0809](S-0809.md) | The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories | review | agent-S-0809 | 2026-11-18T09:19:01Z |
| [S-0810](S-0810.md) | flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails | in-progress | agent-S-0810 | 2026-11-18T09:16:11Z |
| [S-0811](S-0811.md) | The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails | review | agent-S-0811 | 2026-11-18T09:07:44Z |
| [S-0812](S-0812.md) | flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it | review | agent-S-0812 | 2026-11-18T09:15:38Z |
| [S-0834](S-0834.md) | flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount | review | agent-S-0834 | 2026-11-18T08:30:21Z |
| [S-0839](S-0839.md) | No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand | in-progress | agent-S-0839 | 2026-11-18T09:20:08Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 308 | 1469.4340 USD | 182962 | 2026-11-17T17:23:48Z |
| [orchestrator](orchestrator.md) | 4033 | 2265.2108 USD | 1579332 | 2026-11-18T09:18:16Z |
| [analyzer](analyzer.md) | 3 | 31.8290 USD | 2136 | 2026-11-16T06:17:09Z |
| [curator](curator.md) | 261 | 371.9171 USD | 40884 | 2026-11-18T08:36:42Z |
