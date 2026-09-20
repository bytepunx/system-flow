---
id: T-0255
type: task
nature: research
title: "Research: the candidate channels under the same headings, and what comparable tools do"
status: done
parent: S-0071
owner: alex
created: 2026-09-20T07:00:45Z
updated: 2026-09-20T07:12:31Z
transitions:
  - to: ready
    at: 2026-09-20T07:07:14Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:07:14Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:12:31Z
    by: system-flow
stream: S-0071
tags: []
---
# T-0255 Research: the candidate channels under the same headings, and what comparable tools do

## Work
Read flaiover's server (SSE, the MCP bridge, the flai wrapper) and flai dashboard for what exists. For each candidate (WebSocket from flai, SSE plus POST, long polling, MCP in either direction over the existing endpoint, a Unix socket or named pipe mounted into the container, a request queue of files under .flai-cache) fill the story's headings. Cite how CI runners, editor remote agents, and tunnel clients solve the same problem.

## Done when
- Each candidate has every heading filled or marked unknown
- Comparable tools are cited with the date read

## Notes
