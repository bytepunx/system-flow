---
id: S-0036
type: story
nature: feature
title: Dashboard authentication with a per-project token
status: review
parent: E-0006
owner: alex
created: 2026-09-17T19:46:28Z
updated: 2026-09-17T20:03:22Z
transitions:
  - to: ready
    at: 2026-09-17T19:52:11Z
    by: alex
  - to: in-progress
    at: 2026-09-17T19:52:11Z
    by: alex
  - to: review
    at: 2026-09-17T20:03:22Z
    by: alex
blocked:
  - from: 2026-09-17T20:02:21Z
    until: 2026-09-17T20:02:21Z
    reason: auth smoke
tags: [dashboard, cli]
---

# S-0036 Dashboard authentication with a per-project token

## Goal
The dashboard requires a per-project token for every request except liveness and readiness, so it can be published on every interface and later take on document editing. Browsers use a login link and an HttpOnly cookie; agents and the future hub send a bearer header.

## Acceptance criteria
- [x] `flai dashboard` creates `.flai-cache/dashboard.token` (32 random bytes, base64url, mode 0600) on first run, mounts it read-only into the container, and passes `FLAIOVER_TOKEN_FILE`; the token is never an environment variable, argument, or log field
- [x] `flai dashboard token` prints the token and a login link with the token in the URL fragment; `--rotate` replaces it and restarts the container; the start-up message prints the login link
- [x] flaiover accepts `Authorization: Bearer <token>` and a session cookie set by `/login` from the fragment (HttpOnly, SameSite=Lax, Secure when served over HTTPS); the login page rewrites history so the token is not retained; comparison is constant time; writes from the browser also require `X-Requested-With`
- [x] `/_health` and `/_ready` stay open; `/metrics` requires the token unless `FLAIOVER_METRICS_PUBLIC=true`; `/api/*`, `/`, and every page require it; 401 answers carry no detail
- [x] `flai check` errors when `.flai-cache/` is not git-ignored; the template's `.gitignore` covers it
- [x] Request logs never carry header or cookie values; the logging convention gains a project addition saying so; unit tests cover bearer, cookie, missing, wrong, and rotated tokens
- [x] ADR-0018 accepted; docs/operators describes setting, storing, retrieving, rotating, and the exposure table; the security posture says plain HTTP beyond a trusted LAN needs the tunnel or a TLS proxy

## Tasks
- T-0116 flai: token file creation and rotation, read-only mount and FLAIOVER_TOKEN_FILE, flai dashboard token, login link at start-up, check rule for .flai-cache in .gitignore
- T-0117 flaiover: auth hook with bearer and cookie, /login from the fragment, constant-time compare, X-Requested-With on writes, open health and ready, metrics gate, 401 without detail
- T-0118 Tests on both sides; dev script and compose carry the token file; served-build verification
- T-0119 Docs: operators token section and posture, users, logging convention addition, template gitignore; ADR-0018 already accepted

## Notes
First story of E-0006 on purpose: nothing else in the epic ships writes to a dashboard without it. Design in ADR-0018 and design/system/flaiover-dashboard.md (Authentication).
