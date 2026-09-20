---
title: A channel between the dashboard and flai on the host
updated: 2026-09-20
status: active
---

# A channel between the dashboard and flai on the host

The finding of S-0071. The operator asked how the dashboard and a flai process on the host can talk, with flai opening the connection, so that what the user does in the dashboard makes flai act on the host; and then widened it: the dashboard should reach a project **only** through that flai, reads included, and the clone is no longer mounted into the container. This document says what the dashboard does with the mount today, how comparable tools solve "an agent dials out and is sent work", what was tried, what each way of building the channel exposes, and ends with one recommendation. Nothing here is built. The decision and the stories that follow are at the end.

Documentation was read on 2026-09-20. **Tried** means run on this host between a scratch container and scratch processes; everything else was read and not tried.

## What the operator said

Asked in the session on 2026-09-20:

- **What the channel is for.** Pushing and publishing after an acceptance, starting an agent when a story becomes ready, managing the dashboard, running checks for a story in review, and: "make changes to the files via the tunnel through flai itself instead of side-loading the file system into the container, the dashboard is only able to view and interact with the project folder through the flai process itself".
- **What comes first.** "I would like to move all file access back to flai and out of the container doing away with the mount."
- **May an action run with nobody at the host.** Yes, once enabled.
- **How the host side runs.** One flai for all projects.

It began with the operator finding that a story moved to ready starts no work: the move changes a file, and nothing on the host is told.

## What the dashboard does with the mounted clone today

`flai dashboard` mounts the whole clone read-write at its host path ([ADR-0022](../adrs/0022-repository-mounted-at-its-host-path.md)), `.git` included, runs the container as the host user, gives it the host's git identity and excludes, and lays read-only mounts over `.git/hooks`, `.git/info`, `.git/config`, and the token so that the container cannot leave anything the host executes ([ADR-0027](../adrs/0027-git-hooks-config-and-info-are-read-only-in-the-dashboard-container.md)); optionally a push key ([ADR-0026](../adrs/0026-the-dashboard-may-push-with-a-key-the-operator-gives-it.md)). The image carries its own `flai`, `git`, and `ssh`.

The server half of flaiover has no page server code; its API routes are the whole surface. They touch the project in four ways:

| How | What | Where |
|-----|------|-------|
| Direct reads | The manifest; every item file, archive included; `board.md`; threads; narratives (log and open questions); the document tree with each file's front matter; one document; the ADR list; a crawl of every Markdown file into the search index. Front matter is parsed in TypeScript, and the board, the inbox, and the activity view are composed there | `repo.ts`, `board.ts`, `inbox.ts`, `activity.ts`, `search.ts` |
| Watching | chokidar on the design, docs, and wip folders and the manifest. One `change` event per file feeds `/api/events` (SSE to the browser), the search index, the stats and inbox caches, and the notifier | `repo.ts`, `sse.ts`, `notify.ts` |
| Running `flai` in the container | Every write, and every git operation, already goes through flai ([ADR-0016](../adrs/0016-dashboard-delegates-to-flai.md), [ADR-0023](../adrs/0023-documents-are-saved-through-flai.md)): `move`, `order`, `block`, `unblock`, `epic`/`story new`, `accept` (progress streamed as NDJSON) and its dry run, `stream diff`, `stream log`, `thread new`/`reply`/`resolve`, `doc show`/`save`, `adr new`/`accept`, `stats`, `check`, `push --pending --dry-run`. About twenty commands with `--json`: this is already an API | `flai.ts` and its callers |
| A process per MCP session | `/mcp` spawns `flai mcp` in the project and relays JSON-RPC ([ADR-0024](../adrs/0024-mcp-over-http-and-project-identity.md)) | `mcpbridge.ts` |

Nothing writes a file directly, nothing runs git directly, and nothing but Markdown and the manifest is ever read. So the mount gives the container far more than it uses: all source, ignored files such as `.env`, and a writable `.git`.

What the channel must carry to replace it: request and response for reads and for the twenty commands; one stream down to the dashboard (file changes) and streamed progress for acceptance; answers up to a few megabytes (diffs, the document tree); several requests at once, because one page load makes several. With the mount gone the container needs no `flai`, no `git`, no `ssh`, no git identity, no excludes, no guard mounts, and no push key. It still needs its login token, which is a secret handed to it and not part of the clone.

## What comparable tools do

- **CI runners dial out and lease jobs over plain HTTPS.** GitHub's runner long-polls a broker, acknowledges, acquires, and renews a lock while the job runs ([runner source](https://github.com/actions/runner/blob/main/src/Sdk/WebApi/WebApi/BrokerHttpClient.cs)); GitLab's runner polls `POST /jobs/request`, optionally held open, and appends logs by byte offset so a resend is harmless ([long polling](https://docs.gitlab.com/ee/ci/runners/long_polling.html)). Jobs wait in a durable queue while no runner is there. Buildkite moved from polling to a persistent stream in 2026 for latency ([changelog](https://buildkite.com/resources/changelog/345-streaming-job-dispatch-is-now-generally-available/)).
- **The agent decides what may run, not the control plane.** Buildkite's `no-command-eval`, plugin allow-lists, and signed pipelines exist so that a compromised control plane cannot make an agent run arbitrary things ([securing](https://buildkite.com/docs/agent/v3/securing), [signed pipelines](https://buildkite.com/docs/agent/v3/signed-pipelines)). It is the closest mature match to "the container must not be able to make the host do what it likes".
- **Credentials are per agent and separate from users'.** A registration step yields a long-lived agent credential; each job gets a short-lived one ([GitHub runner auth](https://github.com/actions/runner/blob/main/docs/design/auth.md), [Buildkite tokens](https://buildkite.com/docs/agent/v3/tokens)).
- **Tunnel clients hold one outbound connection and multiplex over it**, with application heartbeats and backoff on reconnect (cloudflared, the ngrok agent, VS Code's tunnel). Rancher's remotedialer gives the server a dial into the client's network over a client-opened WebSocket, which is exactly the capability to withhold here.
- **MCP is moving away from servers asking clients.** The 2025-06-18 revision let a server send requests down the client's stream, but only sampling, elicitation, and roots, none of which means "do this"; the 2026-07-28 revision removes sessions, the GET stream, resumption, and server-initiated requests ([changelog](https://modelcontextprotocol.io/specification/latest/changelog)), and the Go SDK cannot receive a custom method from a server.
- **SvelteKit has no WebSocket support** and none planned before version 3 ([issue 1491](https://github.com/sveltejs/kit/issues/1491)); the documented way is a small custom server that wraps adapter-node's handler. Server-sent events from a route work today, and flaiover already uses them. In Go, `coder/websocket` is maintained (ISC, context-aware); an SSE reader is a few lines of the standard library.

## How it was tried

A stub dashboard in a container started the way `flai dashboard` starts flaiover (as the host user, one port published on loopback) with **no project mount**, and a stub on the host in Go that dials it with its own bearer token and serves named methods. The stub dashboard answers "browser" requests only by asking whichever host stub is connected. Two channels were built, a WebSocket carrying JSON-RPC 2.0 and server-sent events down with results posted back, and a third run carried a real `flai mcp` (1.4.4, scratch project) over the host-dialled socket. Scratch project, throwaway tokens; every process was stopped by its PID and the container by its name; the operator's dashboard was not touched. Host: WSL2, Docker Engine.

| What | WebSocket, JSON-RPC | SSE down, POST up |
|------|---------------------|-------------------|
| Request and result, small | about 0.9 ms | about 1.1 ms |
| A 5 MB answer | 53 ms | 53 ms |
| Progress streamed from the host and relayed to the browser as NDJSON | a line a second, in order | a line a second, in order (each chunk its own POST) |
| 20 requests at once while a stream runs | 64 ms for all | not tried |
| Nothing connected | 503 in 8 ms | the same |
| The agent endpoint without the agent token | 401 | 401 |
| The container replaced | the host stub was back 440 ms after it listened again | the same |
| The host stub killed during a streamed request | the browser's request ended at once with "host flai went away"; the next got 503 | not tried |
| The host stub frozen, not dead | marked gone after 9 s, by a ping every 5 s | not tried; needs a heartbeat of its own |
| A real `flai mcp` behind the socket | the stub dashboard called its `board` tool and got the scratch project's board | not tried |

Both work, and on one machine neither is measurably faster. The differences are in what they cost to build and how they fail.

## The candidates

| | Who connects, to what | How a request reaches flai and the answer returns | Nothing connected | Build cost | Verdict |
|-|-----------------------|---------------------------------------------------|-------------------|------------|---------|
| **WebSocket, JSON-RPC 2.0** | flai dials `ws://<dashboard>/agent` with its own token | One ordered pipe; ids pair answers with requests; notifications carry progress and file changes; either side can cancel | 503 at once; requests in flight fail when the socket drops (**tried**) | Go: small. Node: a custom server entry of some thirty lines around adapter-node's handler, and the same hook in the dev server | **Recommended** |
| **SSE down, POST up** | flai opens a GET stream and posts answers | Two channels paired by id; chunks posted separately need sequence numbers to stay in order under load | The same | No custom server; stock routes | Works (**tried**). The fallback if the custom server proves troublesome. A long-lived event stream is also the thing proxies buffer (S-0069's quick-tunnel trial), which matters if the dashboard ever runs elsewhere |
| **Long-poll job lease** | flai polls `POST /agent/lease`, held open | Each request is a leased job with acknowledge, renew, complete | Jobs wait in a durable queue | Most: a job table, leases, expiry | Right when work must wait for an absent agent. Here nothing can even be shown without flai, so nothing should queue; a lease per board read is ceremony |
| **MCP, the dashboard asking flai** | flai as MCP client of `/mcp` | Server-initiated requests | | | No: the protocol is removing them, and they never meant "perform this" |
| **MCP spoken over the reversed socket** | flai dials; the dashboard is the MCP client of flai's existing server | Tool calls | 503 | Small (**tried**) | Works, and would keep `/mcp` alive without a process in the container. Not the dashboard's own API: tool results are text for agents, and the protocol is churning. The operator decided MCP belongs to flai on the host, not behind the dashboard (see Decision) |
| **Unix socket or named pipe mounted into the container** | The container dials the host | | | Small on Linux | No: it is a mount, the direction is the one ruled out, and sockets do not cross Docker Desktop's file sharing on macOS and Windows |
| **Request files under `.flai-cache`** | | | | | No: it needs the mount |
| **A reverse TCP tunnel (remotedialer)** | flai dials; the dashboard may then dial into the host | | | | No: it hands the container a way into the host's network |
| **No channel: flai serves the pages itself** | The browser talks to flai | | | The server half of flaiover rewritten in Go, the pages built static and embedded | The smallest system if the dashboard will only ever run beside flai. It gives up a dashboard that can live somewhere else, and the container. Named so that it is rejected knowingly |

**Framing.** JSON-RPC 2.0, as LSP and MCP use: requests, answers, and notifications in both directions, ids for pairing, trivially logged. It promises nothing about delivery; that is decided below.

**What the methods are.** Not a command line and not MCP tools: named methods with typed arguments, one per thing the dashboard does today. For writes they are the twenty commands above (`item.move`, `item.order`, `story.new`, `accept.run` with progress, `doc.save` with its hash, and so on), built in flai from the arguments the way `_moveArgs` and its siblings build them in TypeScript today, so flai's rules stay the only rules. For reads there are two ways, and the recommendation takes them in order: first `files.list` and `file.read`, limited by flai to Markdown under the manifest's three folders and the manifest itself, which is exactly what the dashboard reads now, so `repo.ts` changes its source and nothing above it changes; then structured answers where flai already has them (`board`, threads, one item, `doc show`) and where TypeScript duplicates flai on purpose (`board.ts` mirrors `flai board`). File changes arrive as notifications from a watcher in flai and feed everything the chokidar events feed now.

**Delivery.** At most once, and honest about it. With no flai connected the dashboard cannot show a project at all, so it says so and queues nothing. A request in flight when the connection drops fails visibly (**tried**). A write carries a request id; flai remembers the ids it has completed for a few minutes and answers a repeat with the recorded result, so a retry after a reconnect cannot accept or push twice. Pings every few seconds in both directions; reconnection with backoff and jitter from a quarter of a second to a few seconds (**tried** without jitter). Cancellation is a notification; flai cancels the work and, for a child process, terminates and then kills it. Messages are capped where flaiover's `maxBuffer` is today (16 MiB); a larger diff is paged.

**Who is who.** flai authenticates with a credential of its own, not the browser's token: generated on the host, handed to the container as a secret when it starts, accepted only on the agent endpoint, one connection per credential. In its first message each side proves it holds the credential over a nonce, so a flai that dials the wrong port does not hand its credential to a stranger. Every request names the project by its key ([ADR-0024](../adrs/0024-mcp-over-http-and-project-identity.md)), and flai refuses a key it does not serve.

## One dashboard, or one per project

Without a mount, nothing in the container belongs to one project. The operator chose one flai for all projects; one dashboard for all of them follows naturally: flai dials once and announces the projects it serves, the dashboard routes by project key, and the multi-project view E-0007 wanted exists without a block of ports, a registry of tokens, or pages embedded through a proxy. It changes [ADR-0018](../adrs/0018-dashboard-token.md): the login token becomes the operator's, kept in flai's home and not in each project's `.flai-cache`. It can come second: the channel and the removal of the mount work with today's one container per project, and `flai dashboard` keeps starting it.

A dashboard that needs no mount can also run somewhere other than the host, and because flai dials out, the host then needs no inbound port and no tunnel client: the channel is the tunnel. That answers most of what S-0069 was asked. It needs TLS and a stronger credential, and it lets the machine that runs the dashboard read what passes. It is noted, not designed.

## Security

Today a holder of the dashboard token can do what the API allows, and a compromised container can read the whole clone, write the work tree and the refs, and use whatever key it was given. With the channel and no mount:

- **A holder of the dashboard token** can do what the API allows, as now, plus the host actions the operator has enabled. That is the point of "yes, once enabled", and it makes the token worth more: with pushing enabled, the token holder can publish any story an agent has put in review, which is what ADR-0026 already accepted for the push key.
- **A compromised container** can call the same named methods and nothing else. It cannot read source or ignored files, cannot write a file except through flai's rules, cannot touch `.git`, holds no credential for the remote, and cannot run anything on the host. It can lie to the browser and misuse the enabled actions while it lasts. The route ADR-0027 closed stops existing, and so does the reason for the push key.
- **Neither** can send a command line. Every method is implemented in flai, takes typed arguments that are validated as data (an item ID, a state, a repository-relative Markdown path), and never reaches a shell. Host actions are off until the operator enables them by name in the host's flai config. The command that starts an agent is written in that config by the operator; the dashboard supplies only the story's ID.
- **What flai must get right.** Path arguments confined to the manifest's folders; the project key checked on every request; the credential never logged; a journal of host actions (what, for which project, asked by whom, when, the outcome) so that an unattended action leaves a trace.

ADRs this would touch: 0016 refined (flai does the reads too); 0022, 0026, and 0027 superseded; 0018 and 0024 refined; 0007 refined (a custom server entry).

## How the host side runs

One process per user, `flai serve` for the sake of a name, serving every project registered with it (a list in flai's home, added to by `flai dashboard` in a project). It needs no root: `flai dashboard` starts it detached if it is not running, as it starts the container, and `flai dashboard stop` leaves it for the other projects; a systemd user unit or a login item can come later. `flai dashboard status` says whether it runs, since when, which projects it serves, and whether the dashboard has it connected. The dashboard shows the same in its header, and when no flai is connected every page says so, with the command to start it, and offers nothing else. A foreground `flai serve` in a terminal stays possible, for watching it work.

## Recommendation

**A WebSocket that flai opens to the dashboard, carrying JSON-RPC 2.0, with named, typed methods implemented in flai; one `flai serve` per user; reads through restricted file methods first and structured methods after; host actions off until enabled; then the mount removed.** It is the one candidate that is a single ordered pipe in both directions with cancellation and progress for free, it was as fast as anything tried, it fails visibly, and it is the transport proxies treat best if the dashboard ever leaves the host. Its one real cost is a custom server entry for flaiover, because SvelteKit has no WebSocket support of its own; if that proves fragile, SSE down with POST up carries the same methods and was tried.

Build order as recommended, each a story, the dashboard working after every one. The decision below changed two things, how reads move and where MCP lives, and lists the stories as they were created:

1. The channel and nothing else: `flai serve`, the agent endpoint, the credential, the handshake, pings, reconnection, status on both sides. The mount stays.
2. Reads through the channel: `files.list`, `file.read`, and change notifications from a watcher in flai; `repo.ts` reads from the channel.
3. Writes through the channel: every call in `flai.ts` becomes a method, acceptance streamed; the image drops `flai`, `git`, and `ssh`.
4. The mount removed: `flai dashboard` stops mounting the clone, the guard mounts, the git identity, the excludes, and the push key; the ADRs above are superseded; the documentation follows.
5. Push and publish after an acceptance, as an enabled host action; the board's "accepted, not pushed" state becomes "pushing".
6. Start an agent when a story becomes ready, with the operator's command.
7. One dashboard for every project.
8. Later: manage the dashboard from its own page; run checks for a story in review; `/mcp` relayed to `flai mcp` on the host.

## Decision

The operator decided on 2026-09-20, in the session, with this finding in front of them ([ADR-0029](../adrs/0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)):

- **The channel:** a WebSocket that flai opens to the dashboard, carrying JSON-RPC 2.0, as recommended.
- **Dashboards:** one dashboard for every project is the end state, built after the mount is gone, as recommended.
- **Reads:** structured methods only. The recommendation was restricted file methods first, to remove the mount sooner; the operator chose the cleaner end state from the start: flai answers the board, items, threads, documents, the ADR list, narratives and activity, the designer's inbox, and search as data, and the dashboard stops parsing project files. There is no `file.read`. The read stories below are larger for it, and the mount goes later than it would have.
- **MCP:** it moves down to flai. In the operator's words: "lets move the MCP integration down to flai as a process again (flai MCP can start a server process) since we want flai and agents working together on the host and not having agents interacting with a project through the dashboard." The dashboard's `/mcp` endpoint goes; flai serves MCP itself on the host, over stdio as now and over HTTP as a process of its own for agents that need it. The MCP-over-HTTP part of ADR-0024 is superseded when that story lands; its project identity stays.

The stories, in build order, under E-0003:

1. S-0072: the channel and nothing else; the mount stays.
2. S-0073: work items read through the channel (project, board, items, threads) and change notifications from a watcher in flai.
3. S-0074: documents, ADRs, narratives and activity, the inbox, and search read through the channel; the dashboard parses no project file.
4. S-0075: every write through the channel; the image drops `flai`, `git`, and `ssh`.
5. S-0076: MCP served by flai on the host, over HTTP as well as stdio; the dashboard's `/mcp` removed. Landed as [ADR-0030](../adrs/0030-mcp-is-served-by-flai-on-the-host-over-stdio-and-http-and-the-dashboard-s-api.md): a process of its own per project (`flai mcp start`), not part of `flai serve`, serving revisions with and without sessions at one address.
6. S-0077: the mount removed, with everything that existed to make it safe. Landed as [ADR-0031](../adrs/0031-the-dashboard-s-container-holds-nothing-of-the-project-a-port-and-two-secrets.md): a port and two secrets, the push key retired, the container still run as the host user so it can read them.
7. S-0078: push and publish after an acceptance, as a host action.
8. S-0079: an agent started when a story becomes ready.
9. S-0080: one dashboard for every project.
10. S-0081 and S-0082, later: the dashboard managed from its own page; checks run for a story in review.
