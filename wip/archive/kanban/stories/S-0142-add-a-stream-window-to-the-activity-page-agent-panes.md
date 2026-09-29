---
id: S-0142
type: story
nature: feature
title: Add a stream window to the activity page agent panes
status: done
parent: E-0011
owner: alex
created: 2026-09-28T23:26:54Z
updated: 2026-09-29T06:02:46Z
transitions:
  - to: ready
    at: 2026-09-29T03:17:59Z
    by: alex
  - to: in-progress
    at: 2026-09-29T05:43:11Z
    by: agent-S-0142
  - to: review
    at: 2026-09-29T06:00:49Z
    by: agent-S-0142
  - to: done
    at: 2026-09-29T06:02:46Z
    by: alex
tags: [dashboard, cli]
touches: [flaiover/src, flai/cmd, flai/internal/serve, flai/internal/hostapi, docs/users, docs/operators, design/system/flai-cli.md, design/system/flaiover-dashboard.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
---
# S-0142 Add a stream window to the activity page agent panes

## Goal

flai serve should provide active agent streams via command for the dashboard. the dashboard should provide live feeds from the agent under the activity page so that agent activity is more observable.

## Acceptance criteria
- [x] Active agent streams are available
- [x] Active agent streams are visible on the agent activity cards

## Tasks
- T-0536 flai reads an agent's stream from its log, from an offset and bounded
- T-0537 flai serve offers an agent's stream to the dashboard and flai serve agent stream prints it
- T-0538 The activity page's agent panes show a live stream window

## Notes

- flai serve offers each story's agent stream as `agent.stream` on the channel and as `flai serve agent stream <story> [--from] [-f] [--json]`, read from the log it already writes (`serve.Stream`). Verified by tests and by running the command against this session's own live log.
- The activity page's cards show it in a stream window (`AgentStream.svelte`), open and followed while the agent runs. Verified by component and route tests in jsdom. Not seen in a browser: a dev server has data only when a `flai serve` dials in to it, and starting a second one on this host is not done. The operator sees it once the host's flai and the dashboard image carry S-0142.
