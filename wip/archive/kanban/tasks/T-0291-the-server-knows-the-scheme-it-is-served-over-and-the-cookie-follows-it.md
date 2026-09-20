---
id: T-0291
type: task
nature: feature
title: The server knows the scheme it is served over, and the cookie follows it
status: done
parent: S-0083
owner: alex
created: 2026-09-20T14:00:06Z
updated: 2026-09-20T14:03:24Z
transitions:
  - to: ready
    at: 2026-09-20T14:00:08Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T14:00:08Z
    by: system-flow
  - to: done
    at: 2026-09-20T14:03:24Z
    by: system-flow
stream: S-0083
tags: []
---
# T-0291 The server knows the scheme it is served over, and the cookie follows it

## Work
adapter-node reports https when no protocol header is configured, so the session cookie was marked Secure over plain HTTP and browsers dropped it everywhere but localhost. The image sets the protocol header to X-Forwarded-Proto, and the server entry fills that header in from the connection when no proxy sent it, in a small module of its own so it can be tested. isSecure follows the request's scheme. How an operator behind a proxy tells it is decided and recorded.

## Done when
- Tests: the scheme for a plain connection, a TLS one, and a forwarded one; the cookie's attributes for HTTP, HTTPS, and forwarded HTTPS through the login route
- The flaiover tests and build pass

## Notes
