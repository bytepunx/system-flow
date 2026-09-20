---
id: T-0256
type: task
nature: research
title: "Try the strongest candidates between a scratch container and a host stub: request and result, reconnection, nothing connected"
status: done
parent: S-0071
owner: alex
created: 2026-09-20T07:00:45Z
updated: 2026-09-20T07:12:32Z
transitions:
  - to: ready
    at: 2026-09-20T07:12:32Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:12:32Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:12:32Z
    by: system-flow
stream: S-0071
tags: []
---
# T-0256 Try the strongest candidates between a scratch container and a host stub: request and result, reconnection, nothing connected

## Work
A dashboard-side stub in a container run the way flai dashboard runs flaiover (published port, host user, project mounted) and a flai-side stub on the host. For the two or three strongest candidates: a request from the dashboard side reaches the host stub and the result comes back; either side restarts and the channel recovers; a request made while nothing is connected is refused or queued as the candidate says. Scratch projects and throwaway tokens only; stop every process by PID or container name (I-0025).

## Done when
- Each trial is written up with what was run and what was seen
- Everything is removed afterwards

## Notes
