---
id: S-0071
type: story
nature: research
title: Research how the dashboard and a flai process on the host can talk, so that actions in the dashboard trigger actions on the host
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:00:03Z
updated: 2026-09-20T07:28:36Z
transitions:
  - to: ready
    at: 2026-09-20T07:00:22Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T07:00:44Z
    by: system-flow
  - to: review
    at: 2026-09-20T07:27:16Z
    by: system-flow
  - to: done
    at: 2026-09-20T07:28:36Z
    by: alex
tags: [dashboard, cli]
touches: [design/system, design/adrs]
---
# S-0071 Research how the dashboard and a flai process on the host can talk, so that actions in the dashboard trigger actions on the host

## Goal
Find out how the dashboard and a flai process on the host can talk to each other, so that once flai has connected to the dashboard (the connection is made from the CLI to the dashboard, never the other way), something the user does in the dashboard can make that flai process act on the host. Today the dashboard can only change files in the clone: it cannot push unless given a key (ADR-0026), cannot start an agent when a story becomes ready, and cannot run anything on the host, by design (ADR-0027). The operator widened it on 2026-09-20, when asked what the channel is for: the dashboard should reach the project **only** through that flai process, reads included, and the mount of the clone into the container goes away. The channel is then not an occasional trigger but the dashboard's whole access to a project, and one flai on the host serves every project. The deliverable is a written finding that compares the ways to build this channel, says what was tried, states what each exposes, and ends with one recommendation; then the operator's decision, recorded, and the stories that would build it. Nothing is built here.

## Acceptance criteria
- [x] The finding starts from what the channel is for, with the operator's answers recorded: which actions in the dashboard should reach the host (candidates found so far: push an acceptance and publish the template, start an agent session when a story becomes ready or a thread is opened, restart or upgrade the dashboard, run the project's tests for a story in review), which of those matter first, and whether one flai serves one project or several
- [x] The finding lists everything the dashboard does with the mounted clone today (direct reads, file watching, the `flai` and `git` it runs inside the container, the MCP bridge, the token file) and what the channel must carry to replace each: request and response, streams (events, acceptance progress), large answers (diffs, search), and concurrency; and says what the container still needs once the mount is gone, and whether it still needs `flai` and `git` inside at all
- [x] It weighs one dashboard per project against one dashboard for every project that the single host flai serves, and what that means for the project token (ADR-0018), project identity (ADR-0024), and the multi-project view E-0007 was after; and notes what a dashboard that needs no mount makes possible (running somewhere other than the host) without designing it
- [x] The candidates are compared under the same headings: who opens the connection and to what address, transport and framing, how a request from the dashboard reaches flai and how the result comes back, what happens when no flai is connected (refused, queued, for how long), reconnection and delivery (at most once, at least once, acknowledgements, a request that outlives a restart of either side), how flai authenticates to the dashboard and how the dashboard knows which flai it is talking to, what it costs to build on each side, and how it behaves on Linux, WSL2, macOS, and Windows with Docker Desktop
- [x] The candidates include at least: a WebSocket that flai opens to the dashboard; server-sent events from the dashboard with results posted back over HTTP; long polling in the manner of MCP `wait_for_events`; flai as an MCP client of the dashboard's existing `/mcp` endpoint, or the reverse over the same connection (ADR-0024); a Unix domain socket or named pipe mounted into the container; a queue of request files under `.flai-cache` that flai watches; and any better one found on the way, with what comparable tools do (CI runners, editor remote agents, tunnel clients) cited
- [x] Security is stated plainly for each: what a holder of the dashboard token can make the host do, what a compromised container can, and what neither can. The finding proposes the boundary: named actions with typed arguments that flai implements and the operator enables, never a command line from the dashboard; which actions need a confirmation on the host or a second credential; how this sits with ADR-0018, ADR-0026, and ADR-0027, and which of them it would refine or supersede
- [x] How the host side runs is covered: a foreground `flai` command, something `flai dashboard` starts and stops with the container, or a user service; without root; what the operator sees when it is connected, idle, or gone (in the terminal, in `flai dashboard status`, and in the dashboard), and how it is stopped
- [x] Anything the finding says works was tried with scratch projects, scratch containers, and throwaway tokens on this host, never the operator's dashboard, credentials, or this repository's remote: at least the two or three strongest candidates, with a dashboard-side stub and a flai-side stub exchanging a request and a result through the container boundary, a reconnection after either side restarts, and a request made while nothing is connected. What was not tried is marked so
- [x] The finding ends with one recommendation and why each other option lost, in a document under `design/system`, linked from its index; the operator's decision is recorded as an ADR
- [x] Follow-up stories exist for what the decision needs, in the order they should be built, each with a goal and acceptance criteria drawn from the finding: the channel itself, the first action or two over it, and what the dashboard shows. If the decision is to build nothing, the notes say why

## Tasks
- T-0254 Conversation with the operator: which dashboard actions should reach the host, which first, and one flai per project or one for several
- T-0255 Research: the candidate channels under the same headings, and what comparable tools do
- T-0256 Try the strongest candidates between a scratch container and a host stub: request and result, reconnection, nothing connected
- T-0257 Write the finding with the security boundary and one recommendation, get the decision, record it as an ADR, and queue the follow-up stories

## Notes
The operator's answers, 2026-09-20, in the session. Which actions: all of push and publish, start an agent, manage the dashboard, run checks for review, and "make changes to the files via the tunnel through flai itself instead of side-loading the file system into the container, the dashboard is only able to view and interact with the project folder through the flai process itself". Which first: "I would like to move all file access back to flai and out of the container doing away with the mount". May actions run with nobody at the host: yes, once enabled. How the host side runs: one for all projects.

Raised by the operator on 2026-09-20: "I would like to queue up a story to research options for how the dashboard and an instance of the flai process can inter-communicate such that once the connection is established between them (from the cli to the dashboard) actions taken by the user in the dashboard can trigger actions in the flai process itself on the host machine."

Where it comes from. The same day the operator found that moving a story to ready starts no work: nothing on the host is told, and an agent only learns of it while it holds `wait_for_events` or at its next `inbox`. S-0052 found the same gap for pushing: an acceptance made from the board stays local unless the container is given a key (ADR-0026), and an agent on the host has pushed by hand after each one. A flai on the host that the dashboard can ask would close both without putting credentials or the ability to execute into the container.

What exists. The dashboard runs `flai` inside the container for every write (`flaiover/src/lib/server/flai.ts`). `/api/events` is server-sent events, one event per changed file. `/mcp` is MCP over Streamable HTTP with the project token as bearer (ADR-0024), bridged to `flai mcp` in the container. `flai dashboard` starts, stops, and reports the container and mounts `.git/hooks`, `.git/info`, `.git/config`, and the token read-only so the container cannot leave anything the host executes (ADR-0027); whatever this story recommends must not reopen that route by another name. The direction the operator gave, CLI to dashboard, fits: the container publishes one port and the host dials it, on every platform.

The model for the shape is S-0052: research tasks, a conversation with the operator, the finding, the decision, the follow-up stories. Research is accepted without a release (ADR-0025). E-0006 and E-0007 are cancelled; this sits under E-0003, the dashboard.
