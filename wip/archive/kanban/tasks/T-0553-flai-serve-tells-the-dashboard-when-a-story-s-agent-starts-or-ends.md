---
id: T-0553
type: task
nature: feature
title: flai serve tells the dashboard when a story's agent starts or ends
status: done
parent: S-0154
owner: alex
created: 2026-09-29T07:08:18Z
updated: 2026-09-29T07:12:38Z
transitions:
  - to: ready
    at: 2026-09-29T07:08:31Z
    by: agent-S-0154
  - to: in-progress
    at: 2026-09-29T07:08:32Z
    by: agent-S-0154
  - to: done
    at: 2026-09-29T07:12:38Z
    by: agent-S-0154
stream: S-0154
tags: []
usage:
  source: log
  seconds: 246
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 47
      output: 13247
      cache_read: 3379812
      cache_write: 44147
      cost: 1.2943
---

# T-0553 flai serve tells the dashboard when a story's agent starts or ends

## Work

An agent's run starting or ending changes no file of the project, so the dashboard hears of it only by polling every 15 s. The launcher calls back when it records a run started, failed to start, or ended, and flai serve sends the channel notification `agent` with the project and the story.

## Done when

A behaviour test shows the launcher calls back with the story on start and on end, `flai serve` sends `agent` on the channel, and `design/system/dashboard-host-channel.md` lists the notification.

## Notes
