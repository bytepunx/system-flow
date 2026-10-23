---
id: D-0014
title: "Criterion 1: may someone attending still hold a newly ready story back: drop the hold for stories that have room, add a restart command and button to stories when the agent drops or fails so that a…"
kind: answer
sources: [TH-0010]
topics: [flai, cli, go, flaiover, dashboard, sveltekit, code, all]
paths: [flaiover/src, flai/cmd, flai/internal/serve, flai/internal/hostapi, flai/internal/harness, flai/internal/config, flai/internal/manifest, flai/internal/itemedit, design/adrs, design/system, design/issues, docs/operators, docs/users]
date: 2026-09-24
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0014 Criterion 1: may someone attending still hold a newly ready story back: drop the hold for stories that have room, add a restart command and button to stories when the agent drops or fails so that a…

The operator, alex, answered TH-0010 at 2026-09-24T08:58:56Z.

What was asked:

```text
flai serve already starts an agent for each ready story while the in-progress limit has room. The one exception is while someone is attending, meaning an MCP cursor or a narrative was written in the last 6 minutes. That exception is why S-0115 and S-0116 sat in ready today. Every agent connected to MCP as `system-flow`, so their activity looked like someone attending (I-0037). The serve log shows "an agent is attending the project (.flai-cache/mcp/system-flow.json ...)" at every look from 08:37 to 08:44.
```

The answer:

```text
drop the hold for stories that have room, add a restart command and button to stories when the agent drops or fails so that a new agent is assigned to the story without it having to re-enter the ready state. write the new ADR to supersede where necessary
```

Source: TH-0010.
