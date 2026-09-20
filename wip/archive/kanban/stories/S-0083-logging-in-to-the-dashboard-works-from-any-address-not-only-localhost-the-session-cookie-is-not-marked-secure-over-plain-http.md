---
id: S-0083
type: story
nature: remediation
title: "Logging in to the dashboard works from any address, not only localhost: the session cookie is not marked Secure over plain HTTP"
status: done
parent: E-0003
owner: alex
created: 2026-09-20T13:56:28Z
updated: 2026-09-20T14:22:43Z
transitions:
  - to: ready
    at: 2026-09-20T13:57:05Z
    by: alex
  - to: in-progress
    at: 2026-09-20T13:59:18Z
    by: system-flow
  - to: review
    at: 2026-09-20T14:07:40Z
    by: system-flow
  - to: done
    at: 2026-09-20T14:22:43Z
    by: alex
tags: [dashboard]
touches: [flaiover/src, flaiover/server.js, flaiover/Dockerfile, docs/operators]
---
# S-0083 Logging in to the dashboard works from any address, not only localhost: the session cookie is not marked Secure over plain HTTP

## Goal
Logging in works from any address the dashboard is reached at, not only `localhost`. Today a correct token from another machine, or from this one by its LAN address, lands back on the login page.

## Acceptance criteria
- [x] Over plain HTTP the session cookie is not marked `Secure`; over HTTPS, and behind a proxy that says so with `X-Forwarded-Proto: https`, it is
- [x] The server knows the scheme it is actually served over: the image no longer lets adapter-node assume `https` when nothing says so, and how an operator behind a TLS-terminating proxy or tunnel tells it (`ORIGIN`, or the protocol and host headers) is decided, tested, and written in the operators' documentation
- [x] Tried in a browser against a scratch dashboard by a non-loopback address over HTTP: the login link and a pasted token both log in, a reload stays logged in, a write from the board works (the CSRF header path), and logout clears the cookie; and by `localhost` as before
- [x] A test holds the cookie's attributes for HTTP, HTTPS, and a forwarded HTTPS request, so this cannot come back unnoticed
- [x] The login page says something useful when the token was right and the session still did not stick, instead of silently asking again

## Tasks
- T-0291 The server knows the scheme it is served over, and the cookie follows it
- T-0292 The login page says so when the token was right and the session did not stick
- T-0293 Operators are told how the scheme is known behind a proxy, and the design records the cause
- T-0294 Tried in a browser by a non-loopback address over HTTP, and by localhost

## Notes
Reported by the operator on 2026-09-20. Reproduced the same day against the running dashboard (flaiover 0.22.3): `POST /api/login` answers 200 with `Set-Cookie: ...; SameSite=Lax; Secure` both by `127.0.0.1` and by the host's LAN address, over plain HTTP. Browsers treat `http://localhost` as a secure context and keep a `Secure` cookie there; on any other plain-HTTP origin they drop it, so the next request has no cookie and is sent to `/login`. `isSecure()` in `flaiover/src/lib/server/auth.ts` trusts `url.protocol`, and adapter-node reports `https:` whenever neither `ORIGIN` nor a protocol header is configured; `flaiover/server.js` and the Dockerfile set neither. It is not known since when: possibly since authentication landed (S-0036), noticed now because the dashboard is published on every interface by default. Bearer tokens (agents, curl) are not affected.
