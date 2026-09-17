---
id: ADR-0018
title: flaiover authenticates every request with a per-project token
status: accepted
date: 2026-09-17
supersedes: []
superseded_by: []
refines: [ADR-0016]
---

# ADR-0018 flaiover authenticates every request with a per-project token

## Context

ADR-0007 and ADR-0016 kept the dashboard unauthenticated because it was a local, read-mostly tool bound to loopback. S-0035 published it on every interface on the operator's instruction, and E-0006 makes it the place where the designer edits design intent and accepts work. A writer reachable over the network needs a credential, and later a multi-project hub reached over a tunnel will present the same credential from outside.

## Decision

Every request except `/_health` and `/_ready` requires a per-project token. `flai dashboard` generates it on first run (32 random bytes, base64url), stores it at `.flai-cache/dashboard.token` with mode 0600, mounts the file read-only into the container, and points flaiover at it with `FLAIOVER_TOKEN_FILE`. The token is never an environment variable, a command argument, a log field, or a committed file; `flai check` refuses a repository where `.flai-cache/` is not git-ignored.

Browsers receive a login link with the token in the URL fragment, which never leaves the browser; the login page sets an HttpOnly, SameSite cookie (Secure over HTTPS) and rewrites history. Agents and the hub send `Authorization: Bearer`. The bearer check is primary and the cookie is a convenience over it. Comparison is constant time; writes from the browser also require a custom header. `flai dashboard token` prints the token and link, `--rotate` replaces it and restarts the container, invalidating every session. `/metrics` requires the token unless the operator opts it public.

flaiover speaks plain HTTP. Beyond a trusted LAN the operator terminates TLS in the tunnel or a proxy; the dashboard does not grow a TLS stack.

## Consequences

- One token per project, so the hub holds a keyring and rotation is local to a project.
- Anyone with Docker socket access on the host can read the mounted file; that is host ownership and out of scope.
- The security posture text changes from "keep it on localhost" to "keep the host private, use the tunnel or TLS beyond the LAN".
- Sessions are the token itself, so there is no session store to persist or leak.

## Alternatives considered

- Passwords and users: the dashboard has one designer per project today; the hub is where identity belongs.
- Token in `~/.flai/config.json`: shared across projects and copied more freely; rejected.
- Token as an environment variable: visible in `docker inspect`; rejected.
- Built-in TLS with a self-signed certificate: adds a certificate lifecycle for little gain; the tunnel already terminates TLS.
