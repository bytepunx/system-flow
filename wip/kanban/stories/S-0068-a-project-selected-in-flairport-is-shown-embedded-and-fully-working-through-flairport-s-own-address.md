---
id: S-0068
type: story
nature: feature
title: A project selected in flairport is shown embedded and fully working, through flairport's own address
status: ready
parent: E-0007
owner: alex
created: 2026-09-20T01:18:59Z
updated: 2026-09-20T03:32:02Z
transitions:
  - to: ready
    at: 2026-09-20T03:32:02Z
    by: alex
tags: [flairport, dashboard]
touches: [flairport/src, flaiover/src]
---
# S-0068 A project selected in flairport is shown embedded and fully working, through flairport's own address

## Goal
Choosing a project in flairport shows that project's flaiover dashboard inside flairport's page, fully working, so the operator moves between projects without leaving the one address, and that one address is all a tunnel has to expose.

## Acceptance criteria
- [ ] How the dashboard is brought into the page is decided first, with what was tried, in an ADR: flairport reverse-proxies each project under its own origin (a path such as `/p/<key>/` or a subdomain per project) and frames that, rather than framing `http://localhost:<port>` directly. Direct framing is the thing to rule out or in with evidence: through a tunnel only flairport's address is reachable, and flaiover's session cookie is HttpOnly and SameSite, so a cross-origin frame would not be logged in
- [ ] flairport adds the project's bearer token to proxied requests on the server side, from the instance record (S-0065). The browser never receives a project token, and the operator is not asked to log in to each project
- [ ] Everything in the embedded dashboard works as it does on its own: navigation and deep links (the address bar reflects the project and the page inside it, and a reload comes back to the same place), the board's drag and drop, the editor, forms, server-sent events for live refresh, streamed acceptance, downloads, and the light and dark theme following flairport's
- [ ] If flaiover has to change to live under a path or behind a proxy (SvelteKit `paths.base`, API calls that assume `/api`, the login redirect, the CSRF header check, cookies, `X-Forwarded-*`), it changes here, behind a setting that leaves a dashboard run on its own exactly as it is today, with tests
- [ ] MCP over HTTP (`/mcp`, ADR-0024) either works through the proxy with the agent's own bearer token or is deliberately not proxied, and the story says which and why
- [ ] A project whose dashboard is stopped or unreachable shows that in the pane with the start action from S-0067, not a browser error; switching projects never leaks one project's page, token, or events into another's
- [ ] flaiover may only be framed by flairport: it sends `Content-Security-Policy: frame-ancestors` (none by default, flairport's origin when told of it) so that no other site can frame a dashboard; flairport's own pages cannot be framed at all
- [ ] Tried in a browser with at least two scratch projects running at once, including through a real tunnel or a stand-in that terminates HTTPS on another origin; `design/system/flairport.md`, `flaiover-dashboard.md`, and `docs/operators/index.md` (the tunnel expectation) updated

## Tasks

## Notes
From E-0007, the operator's fourth bullet: "On selection of a project, present it, embedded, in the browser screen as a 'integration at the glass' approach." And from the epic's outcome: the surface is exposed "via something like cloudflared or a similar tunnel ... with a simple passkey/identity check at the edge".

Why a proxy is likely. A tunnel publishes one local address. If flairport framed `http://localhost:4250`, a browser outside the network could not load it. So each project's dashboard has to be reachable through flairport's own address, which makes flairport a reverse proxy, which in turn lets it hold the tokens server side (ADR-0018 has one token per project and the hub "holds a keyring"; this is that keyring).

What is known about flaiover that bears on it. It authenticates by bearer token or by an HttpOnly, SameSite cookie set at `/login`; writes need an `x-requested-with` header (seen in S-0059: a plain fetch without it is 403); pages build links with `resolve()` from `$app/paths`, but the client's API calls are absolute (`/api/...`), and `EventSource('/api/events')` likewise; the login redirect is to `/login`. A path prefix therefore touches the `api()` helper, every `EventSource`, and the base path setting; a subdomain per project touches none of that and needs wildcard DNS and certificates at the tunnel. The ADR should weigh exactly this.

Depends on S-0065 and S-0067. Setting up the tunnel itself (cloudflared or another, and the identity check at its edge) is the epic's second outcome and is not a bullet yet: this story makes flairport correct behind one and documents what the tunnel must provide; a story for the tunnel's set-up and documentation is worth adding to the epic.
