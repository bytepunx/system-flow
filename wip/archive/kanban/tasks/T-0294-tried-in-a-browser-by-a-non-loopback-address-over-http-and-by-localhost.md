---
id: T-0294
type: task
nature: feature
title: Tried in a browser by a non-loopback address over HTTP, and by localhost
status: done
parent: S-0083
owner: alex
created: 2026-09-20T14:00:07Z
updated: 2026-09-20T14:07:32Z
transitions:
  - to: ready
    at: 2026-09-20T14:07:31Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T14:07:32Z
    by: system-flow
  - to: done
    at: 2026-09-20T14:07:32Z
    by: system-flow
stream: S-0083
tags: []
---
# T-0294 Tried in a browser by a non-loopback address over HTTP, and by localhost

## Work
Scratch dashboard from the branch image under another name and port: by the host's LAN address, the login link and a pasted token, a reload, a write from the board, logout; the same by localhost; a request that claims forwarded HTTPS gets a Secure cookie.

## Done when
- What was tried and what was not is in the narrative; every unexpected answer explained before review

## Notes
