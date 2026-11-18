---
title: Active streams
updated: 2026-11-18T13:04:38Z
---

# Active streams

| Stream | Title | Status | Last agent | Updated |
|--------|-------|--------|------------|---------|
| [S-0809](S-0809.md) | The merge queue's batch runs the suites verify_at.batch lists over the merge of its stories | in-progress | agent-S-0809 | 2026-11-18T09:19:01Z |
| [S-0810](S-0810.md) | flai serve runs a state's suites and the cluster's deploy and checks as a job when a story enters it, skips the cluster's steps where none is configured, and sends the story back when a step fails | review | agent-S-0810 | 2026-11-18T11:23:23Z |
| [S-0811](S-0811.md) | The release's publish runs the suites verify_at.publish lists at the main branch's head before it tags, and publishes nothing when one fails | review | agent-S-0811 | 2026-11-18T09:07:44Z |
| [S-0812](S-0812.md) | flai test --at ci runs the suites verify_at.ci lists, and the template's and this repository's CI call it | review | agent-S-0812 | 2026-11-18T09:15:38Z |
| [S-0834](S-0834.md) | flai kluster exec and log read .flai-cache/kluster, which a story's agent in a container does not mount | review | agent-S-0834 | 2026-11-18T08:30:21Z |
| [S-0839](S-0839.md) | No flai command writes agent.image.args, so a project's layer build arguments can only be written by hand | review | agent-S-0839 | 2026-11-18T10:01:45Z |
| [S-0843](S-0843.md) | Callers of the workflow-aware guard, forecast replay, and share hold still pass the baseline, so a project's own workflow file will not reach them | review | agent-S-0843 | 2026-11-18T11:06:28Z |

## Strategic agents

| Agent | Activities | Cost | Seconds | Last run |
|-------|------------|------|---------|----------|
| [planner](planner.md) | 308 | 1469.4340 USD | 182962 | 2026-11-17T17:23:48Z |
| [orchestrator](orchestrator.md) | 4046 | 2267.9559 USD | 1579855 | 2026-11-18T11:25:21Z |
| [analyzer](analyzer.md) | 3 | 31.8290 USD | 2136 | 2026-11-16T06:17:09Z |
| [curator](curator.md) | 261 | 371.9171 USD | 40884 | 2026-11-18T08:36:42Z |
