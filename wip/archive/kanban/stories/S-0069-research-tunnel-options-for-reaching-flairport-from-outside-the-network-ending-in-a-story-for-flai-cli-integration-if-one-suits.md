---
id: S-0069
type: story
nature: research
title: Research tunnel options for reaching flairport from outside the network, ending in a story for flai CLI integration if one suits
status: cancelled
parent: E-0007
owner: alex
created: 2026-09-20T03:35:30Z
updated: 2026-09-20T06:11:55Z
transitions:
  - to: ready
    at: 2026-09-20T03:38:59Z
    by: alex
  - to: in-progress
    at: 2026-09-20T04:07:41Z
    by: system-flow
  - to: cancelled
    at: 2026-09-20T06:11:55Z
    by: alex
tags: [flairport, cli]
touches: [design/system, design/adrs]
---
# S-0069 Research tunnel options for reaching flairport from outside the network, ending in a story for flai CLI integration if one suits

## Goal
Find out which tunnels can put flairport on the internet the way E-0007 describes: reachable from outside the network over HTTPS that ends at the provider's edge, with a simple passkey or identity check at that edge, and nothing opened inbound on the host. The deliverable is a written finding with a recommendation, the operator's decision recorded, and, if one or more options are suitable, a story for integrating it into the flai CLI. No integration is built here, and "none of them is suitable" is an acceptable outcome.

## Acceptance criteria
- [ ] The finding compares the candidates under the same headings: how the tunnel is made (outbound only, no inbound port), where TLS ends and who holds the certificate, what identity check the edge offers (passkeys, one-time codes, an identity provider, access policies by email) and whether it is part of the free or the paid offer, custom domain versus a provider-assigned name, what an account and a domain cost, limits that matter here (request and body size, idle and streaming timeouts, WebSocket and server-sent events support), how it is installed and run as a service on Linux, WSL2, macOS, and Windows, how it is configured (file, CLI, API, so that flai could drive it), its licence and whether it can be self-hosted, and what the provider can see of the traffic
- [ ] The candidates include at least cloudflared with Cloudflare Access, Tailscale Funnel and Tailscale Serve, ngrok, and one self-hosted option (such as frp, rathole, or a WireGuard VPS with a reverse proxy), plus any better one found on the way; each claim cites the provider's documentation with the date it was read
- [ ] What flairport and flaiover need from a tunnel is written down first and each candidate is checked against it: long-lived server-sent events (the board and the ADRs page), streamed responses (acceptance is newline-delimited JSON), MCP over Streamable HTTP with a bearer token that must pass the edge check or bypass it deliberately (ADR-0024), one hostname or a wildcard (S-0068 decides between a path and a subdomain per project), and the identity the edge established being available to the origin as a header flairport can trust or ignore
- [ ] Anything the finding says works was tried against a scratch service on this host, never the operator's dashboards or accounts: at least the options that need no paid account or domain, with what was run and what was seen; anything not tried is marked as not tried, and anything that needs the operator's account or domain is listed as a question for them rather than attempted
- [ ] The security of the arrangement is stated plainly for each option: flairport has no authentication of its own and holds every project's token (S-0067), so the edge check is the only lock; what happens if the tunnel is up and the edge policy is missing or misconfigured, whether the origin can refuse requests that did not come through the edge (a signed header, a token, binding to loopback), and how access is revoked
- [ ] The operator's answers are recorded in the finding, from a conversation in the session or a thread: which providers they already have accounts or domains with, whether a paid plan is acceptable, who besides them will pass the edge check, and whether agents on other machines must reach MCP through the same tunnel
- [ ] The finding ends with one recommendation and why each other option lost, in a document under `design/system`; the operator's decision is recorded as an ADR when it changes what is exposed or how (refining ADR-0018 and ADR-0024), and as a note in the living design otherwise
- [ ] If at least one option is suitable, a follow-up story exists under E-0007 for integrating it with the flai CLI, with goal and acceptance criteria drawn from the finding: what flai would do (check that the tunnel client is installed, write its configuration, start, stop, and report it beside the dashboards, refuse to expose flairport when no edge policy can be confirmed), what stays the operator's to do by hand (accounts, DNS, the identity policy), and which option or options it covers. If none is suitable, the story's notes say why and no follow-up is created

## Tasks
- T-0245 Research: what flairport needs from a tunnel, and each candidate against the same headings
- T-0246 Try the options that need no account against a scratch service: SSE, streamed responses, bearer tokens
- T-0247 Conversation with the operator: what they use today, accounts and domains, paid plans, who passes the edge, MCP from other machines
- T-0248 Write the finding with the security of each option and one recommendation, get the decision, record it, and queue the flai integration story

## Notes
Raised by the operator on 2026-09-20, after the first four E-0007 stories were written: "add a story to research tunnel options for E-0007. this story's output should result in a story about the possibility of adding integration for this to the flai CLI assuming we can find one or more suitable options."

It covers the epic's second outcome, which had no bullet: "Expose this surface via something like cloudflared or a similar tunnel so that it can be accessed from outside the network over HTTPS (https terminated at the tunnel provider's network edge) with a simple passkey/identity check at the edge." And the epic's note: authentication is not flairport's job; it comes from "the tunnel vendor or ... authstar". Whether authstar can sit behind a plain tunnel as the identity check, instead of the vendor's, is worth a paragraph.

What exists. `docs/operators/index.md` has a section on the tunnel expectation (S-0043): flaiover speaks plain HTTP, and beyond a trusted network the operator ends TLS in a tunnel or proxy. ADR-0018 says the same. ADR-0024 sketched a hub that projects dial out to; flairport (S-0067) is the single-host form of it. The operator told S-0052 on 2026-09-19 that the dashboard is reached by them alone, over a public tunnel, so something is already in use: ask what, first.

Research stories are accepted like any other and cut no release (ADR-0025); the finding and any ADR land on main. S-0052 is the model for the shape: research tasks, then a conversation with the operator, then the finding, the decision, and the follow-up stories.

It does not depend on the other E-0007 stories and can be done first; S-0068's choice between a path and a subdomain per project depends partly on what tunnels make easy (wildcard hostnames and certificates), so doing this before S-0068 helps.
- 2026-09-20T06:11:55Z: moved to cancelled: I am thinking of taking a very different route
