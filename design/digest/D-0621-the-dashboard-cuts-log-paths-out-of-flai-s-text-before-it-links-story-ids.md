---
id: D-0621
title: The dashboard cuts log paths out of flai's text before it links story IDs, because an agent log's name holds its story ID
kind: decision
sources: [S-0785, T-3356]
topics: [flaiover, dashboard]
paths: [flaiover/src/lib/activity.ts, flaiover/src/lib/logs.ts, flaiover/src/lib/components/LogText.svelte]
date: 2026-11-14
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0621 The dashboard cuts log paths out of flai's text before it links story IDs, because an agent log's name holds its story ID

S-0785 made every log path the dashboard shows a link to the log page (T-3356). `reasonParts` in `flaiover/src/lib/activity.ts` already cut a reason or a refusal into text, story IDs, and `/settings/<page>` paths. An agent run's log is named `<key>-<subject>-<stamp>.log`, with the story's ID as subject ([ADR-0300](../adrs/0300-the-dashboard-reads-a-log-flai-wrote-through-log-read-confined-to-the-project-s.md)), so linking story IDs first would break the path in two. The story's narrative, `wip/archive/agents/S-0785.md`, records under Decisions:

> T-3356 widened past its touches to activity.ts (reasonParts cuts log paths out before story IDs, since an agent log's name holds its story ID), DismissibleNotice.svelte (links log paths in every notice once), …

`reasonParts` runs `splitLogPaths` (`flaiover/src/lib/logs.ts`) first and links IDs and settings pages only in the text between the paths. A later change that adds another kind of link to flai's text keeps the log paths cut out before it.
