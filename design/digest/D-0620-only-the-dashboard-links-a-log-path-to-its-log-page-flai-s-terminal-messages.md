---
id: D-0620
title: "Only the dashboard links a log path to its log page; flai's terminal messages keep printing the path, since a CLI user has a shell on the host"
kind: decision
sources: [S-0785, TH-1193]
topics: [flaiover, dashboard, cli, flai]
paths: [flaiover/src/lib/components/LogText.svelte, flaiover/src/lib/logs.ts, flai/cmd]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0620 Only the dashboard links a log path to its log page; flai's terminal messages keep printing the path, since a CLI user has a shell on the host

S-0785 asked that operators be given a way to view a log rather than be told to open a log file. Its planner assumed on TH-1193 that criterion 1 covers what the dashboard shows:

> Criterion 1 covers what the dashboard shows. The CLI's terminal output keeps printing paths, since a terminal user has the shell.

The orchestrator approved the plan on TH-1193: "The story's goal is an operator without shell access, so terminal output stays as it is." The story's narrative, `wip/archive/agents/S-0785.md`, records under Decisions:

> Scope of criterion 1 is the dashboard: an operator at the CLI has a shell on the host and can open the path, so flai's terminal messages keep their paths; what the dashboard renders from them, statuses, whys, and errors, links each log path through LogText.

So the story dropped its `flai/cmd` folder touch, and no flai command's output changed. A later story that touches how flai reports a log need not make the CLI print links or open a viewer: the dashboard turns each log path in flai's text into a link to `/ops/log` (`LogText`, `splitLogPaths` in `flaiover/src/lib/logs.ts`).
