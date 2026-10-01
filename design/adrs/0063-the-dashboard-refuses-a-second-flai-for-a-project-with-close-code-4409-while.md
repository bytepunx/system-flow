---
id: ADR-0063
title: "The dashboard refuses a second flai for a project with close code 4409 while the one it has answers, and flai waits a minute before it dials again"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
refines: [ADR-0029]
---

# ADR-0063 The dashboard refuses a second flai for a project with close code 4409 while the one it has answers, and flai waits a minute before it dials again

## Context

The dashboard keeps one flai connection per project ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). A newer proven connection replaced the older one at once (close 4000). A flai reconnects within a fraction of a second. So two flai serving the same project at one dashboard evicted each other every few seconds, and each request went to whichever held the project at that moment (I-0029). Those two could be a `flai serve` run by hand, or a host on another address. One host per machine ([ADR-0040](0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)) refuses a second host, but not those. T-0339's attempt at this was abandoned with S-0084. Its criterion stood: refuse the newcomer while the holder answers, replace only a holder that has stopped answering, and have a refused flai back off and say so.

## Decision

The dashboard refuses a second flai for a project with close code 4409 while the one it has answers, and flai waits a minute before it dials again.

1. **Judged on adoption.** When a flai proves itself for a project that already has a connection, the dashboard pings the one it has. A pong within two seconds keeps the project with it, and the newcomer is closed with 4409, "the project is held by a flai that answers". No pong, or a holder that has closed, means the holder is closed with 4000 and the newcomer adopted.
2. **One at a time.** Adoptions for a project are judged in turn, so two newcomers are not both compared with the same holder. What a newcomer sends while it is judged is kept and handled once it is adopted.
3. **flai backs off.** flai reads the close code. On 4409 it waits a minute, spread over its second half, before it dials that dashboard for that project again, instead of the 250 ms to 4 s it waits after any other end. It logs this once at warn. Its connection state's `last_error` says another flai serves the project there and when it tries again.
4. **Logged by the dashboard.** The dashboard logs each refusal at warn with both flai versions and the project's key.

## Consequences

- The first flai to connect keeps a project for as long as it answers. A second learns it is not wanted, and dials once a minute in case the first goes.
- A holder that freezes without closing loses the project on the next newcomer, within two seconds, rather than after the next missed ping.
- The refused flai has proven itself, so it marks itself connected for a moment and logs that it connected, before the refusal arrives.
- The protocol stays 1: a flai that does not know 4409 retries at once as before and is refused each time, which costs it a handshake every few seconds and no longer costs the holder anything.

## Alternatives considered

- **Refuse at hello, before the proof.** It would avoid the moment of being connected. But an unproven caller could then learn which projects are held, and make the dashboard ping the holder.
- **Newest wins, as before, with a back-off on 4000.** The project would still change hands each time the older flai came back.
- **A lock per project in flai's state.** Two flai with different configurations, or on different hosts, do not share their state. The dashboard is the one place both reach.
