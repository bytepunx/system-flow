---
id: ADR-0029
title: The dashboard reaches a project only through flai on the host
status: accepted
date: 2026-09-20
supersedes: []
superseded_by: []
refines: [ADR-0007, ADR-0016, ADR-0018, ADR-0024]
---

# ADR-0029 The dashboard reaches a project only through flai on the host

## Context

`flai dashboard` mounts the whole clone into the flaiover container, read-write and at its host path ([ADR-0022](0022-repository-mounted-at-its-host-path.md)), with the host's git identity, and since S-0064 with read-only mounts over the parts of `.git` the host executes ([ADR-0027](0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md)); optionally with a push key ([ADR-0026](0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md)). The container runs its own `flai` and `git` for every write ([ADR-0016](0016-dashboard-delegates-to-flai.md)) and reads Markdown directly. It can do nothing on the host: an acceptance made from the board is not pushed without the key, and a story moved to ready starts no work, which the operator found on 2026-09-20.

S-0071 researched how the dashboard and a flai process on the host could talk, with flai opening the connection, and tried the candidates between a container with no project mount and a process on the host: [A channel between the dashboard and flai on the host](../system/dashboard-host-channel.md). Asked what the channel is for, the operator widened it: "I would like to move all file access back to flai and out of the container doing away with the mount", with "the dashboard ... only able to view and interact with the project folder through the flai process itself". The inventory found that the dashboard reads nothing but Markdown and the manifest and already writes only through flai, so the mount gives the container far more than it uses: all source, ignored files, and a writable `.git`.

## Decision

The dashboard reaches a project only through a flai process on the host, over a connection that flai opens.

- **The channel.** flai dials the dashboard's agent endpoint with a WebSocket and the two speak JSON-RPC 2.0 over it. The dashboard never connects to the host, and the container is given no socket, pipe, or folder of the host's.
- **Named methods.** Every request is a named method with typed arguments, implemented in flai and validated as data there: an item ID, a state, a repository-relative Markdown path. There is no method that takes a command line and none that reads or writes an arbitrary file. Writes are the commands the dashboard runs today. Reads are structured answers from flai (the project, the board, items, threads, documents and their tree, the ADR list, narratives and activity, the designer's inbox, search); the dashboard stops parsing project files. File changes reach the dashboard as notifications from a watcher in flai.
- **Delivery.** At most once. With no flai connected the dashboard shows that and queues nothing. A request in flight when the connection drops fails visibly. A write carries a request ID, and flai answers a repeat of one it has completed with the recorded result. Both sides ping; flai reconnects with backoff.
- **Who is who.** flai authenticates with a credential of its own, generated on the host and handed to the container as a secret, accepted only on the agent endpoint; both sides prove they hold it in the first exchange. Every request names the project by its key ([ADR-0024](0024-mcp-over-http-and-project-identity.md)) and flai refuses a key it does not serve.
- **One flai per user.** One process on the host serves every project registered with it, started without root by `flai dashboard` when it is not running. One dashboard for every project is the end state and comes after the mount is gone; until then `flai dashboard` keeps starting a container per project.
- **Host actions.** Things the dashboard may ask the host to do beyond the project's files (push and publish after an acceptance, start an agent when a story becomes ready, manage the dashboard, run checks for a story in review) are off until the operator enables each by name in the host's flai configuration. Once enabled they run with nobody at the host, by the operator's choice. The command that starts an agent is written in that configuration by the operator; the dashboard supplies only the story's ID. Every host action is journalled on the host.
- **The mount goes.** When reads and writes are on the channel, `flai dashboard` stops mounting the clone, and with it the guard mounts, the git identity, the excludes, and the push key; the image stops carrying `flai`, `git`, and `ssh`. The story that does it supersedes ADR-0022, ADR-0026, and ADR-0027.
- **MCP lives in flai.** Agents work with flai on the host, not through the dashboard: the dashboard's `/mcp` endpoint goes, and flai serves MCP itself, over stdio as now and over HTTP as a process of its own. The story that does it supersedes the MCP-over-HTTP part of ADR-0024; project identity stays.

## Consequences

- flaiover needs a custom server entry around adapter-node's handler, because SvelteKit has no WebSocket support; the development server needs the same hook. This refines [ADR-0007](0007-sveltekit-spa-with-node-adapter-in-docker.md). Server-sent events down with results posted back carry the same methods and were tried; they are the fallback if the custom server proves fragile.
- A compromised container can call the named methods and nothing else: it cannot read source or ignored files, write outside flai's rules, touch `.git`, hold a credential for the remote, or execute on the host. A holder of the dashboard token can do what the API allows plus the host actions the operator enabled, so the token is worth more once pushing is enabled; [ADR-0018](0018-dashboard-token.md) is refined by the agent credential now and by a per-user login token when one dashboard serves every project.
- Without a flai on the host the dashboard shows nothing. `flai dashboard` starting both keeps that from being a step the operator has to remember.
- Choosing structured reads over restricted file reads means more to build before the mount can go, and removes the duplicate readers in TypeScript (`board.ts` mirrors `flai board` on purpose) instead of re-sourcing them.
- A dashboard that needs no mount can run somewhere other than the host, and because flai dials out the host then needs no inbound port and no tunnel client. That needs TLS and a stronger credential and is not designed here.
- ADR-0016 is refined: flai does the reads too.

## Alternatives considered

- Server-sent events down and POST up: works without a custom server; two channels to pair, chunk order to manage, and a long-lived event stream is what proxies buffer.
- A long-poll job lease in the manner of CI runners: right when work must wait for an absent agent; here nothing can be shown without flai, so nothing should queue.
- MCP with the dashboard asking flai: the protocol is removing server-initiated requests, and they never meant "perform this".
- A Unix socket or request files shared with the container: a mount by another name, in the direction ruled out, and sockets do not cross Docker Desktop's file sharing.
- A reverse TCP tunnel: it hands the container a way into the host's network.
- flai serving the pages itself with no container: the smallest system if the dashboard only ever runs beside flai; it gives up a dashboard that can live elsewhere and rewrites flaiover's server half in Go.
- Restricted file reads first (`files.list`, `file.read` for Markdown under the manifest's folders): recommended by the finding to remove the mount sooner; the operator chose structured methods from the start.
