---
id: T-0292
type: task
nature: feature
title: The login page says so when the token was right and the session did not stick
status: done
parent: S-0083
owner: alex
created: 2026-09-20T14:00:07Z
updated: 2026-09-20T14:03:25Z
transitions:
  - to: ready
    at: 2026-09-20T14:03:24Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T14:03:24Z
    by: system-flow
  - to: done
    at: 2026-09-20T14:03:25Z
    by: system-flow
stream: S-0083
tags: []
---
# T-0292 The login page says so when the token was right and the session did not stick

## Work
After a login the page asks something that needs the session; when that is refused although the token was accepted, it says the browser did not keep the cookie and what usually causes it, instead of going round to the login page again.

## Done when
- A component test for the accepted, the refused, and the accepted-but-not-kept cases

## Notes
