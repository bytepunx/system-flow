---
id: ADR-0040
title: "One flai host per machine runs flai serve and each project's MCP server as its children"
status: accepted
date: 2026-09-24
supersedes: [ADR-0034]
superseded_by: []
refines: [ADR-0029]
---

# ADR-0040 One flai host per machine runs flai serve and each project's MCP server as its children

## Context

[ADR-0034](0034-flai-serve-keeps-each-served-project-s-http-mcp-server-running.md) made `flai serve` start and supervise each served project's `flai mcp http`. `flai dashboard` started `flai serve` itself ([ADR-0029](0029-the-dashboard-reaches-a-project-only-through-flai-on-the-host.md)). Nothing supervised `flai serve`. A crash left the dashboards without flai until someone ran `flai dashboard` or `flai serve start` again. Nothing could restart or upgrade it on the dashboard's behalf, since the process that answers the dashboard cannot replace itself. Two `flai serve` processes with different configurations could also run on one machine and fight over one dashboard (I-0029).

The operator asked (S-0106) for one process that loads and manages the others. Its jobs: run `flai serve` and the MCP servers, keep them running, stop them all when it stops, and take commands from `flai serve` (status, version, upgrade, restart) so that the dashboard's host panel (S-0107) can manage them. They also asked for one such process per machine.

## Decision

**`flai host`, one per machine, runs `flai serve` and each served project's `flai mcp http` as its own children. `flai serve` asks it for MCP servers and starts none.**

- **One per machine.** The host listens on a fixed loopback address, `127.0.0.1:4241`, which `FLAI_HOST_ADDR` moves. Only one process can hold it, and a crash frees it. A second host is refused and told which one holds the address (`/_health`: its pid, version, and config file).
- **Control API.** The host answers HTTP on that address. Every route but `/_health` needs a bearer token, which the host makes at each start. It hands the token to its children in `FLAI_HOST_TOKEN` (with `FLAI_HOST_URL`) and writes it with mode 0600 to `host/token` beside the config for the CLI. A request with an `Origin` header is refused, so no web page can reach it. Routes: status; start, stop, or restart `serve`, `mcp`, or `all`; the set of projects whose MCP server to keep; check for a newer flai; upgrade.
- **Children.** The host starts `flai serve` and restarts it when it ends unasked, after a wait that doubles from one second to thirty and resets after a minute of running. `flai serve` tells the host every project it serves, at once when the set changes and every 15 seconds otherwise. The host keeps a `flai mcp http` for each project, choosing ports as ADR-0034 did, and stops the ones no longer named. When `flai serve` stops, it tells the host nothing, so a restart of serve does not bounce the MCP servers. A process the host did not start is used and left alone (`external`): a `flai serve` started by hand, or an MCP server from `flai mcp start`.
- **Leaving with the host.** When the host returns, it stops every child: SIGTERM, then SIGKILL after ten seconds. Every child is started with `--exit-with <host pid>` and ends within seconds of a host that was killed outright. A Linux parent-death signal was not used: Go may fork from a thread that later exits, which would kill the child.
- **Upgrade.** The host runs `flai self-upgrade`. When that installs a release, the host answers, stops its children, and replaces itself with the new binary: exec on Unix, a new process on Windows. The new host then starts new children.
- **Who starts it.** `flai dashboard` registers its projects and starts the host, not `flai serve`. `flai serve start` and `stop` go through the host. Bare `flai serve` in a terminal still runs, keeps no MCP server, and says so.
- **The dashboard.** The dashboard reaches the host only through `flai serve`, with four hostapi methods, each running `flai host … --json`. `host.status` and `host.check` are reads. `host.process` and `host.upgrade` are gated on a new host action, `host`, and detached from the request, because restarting serve ends the connection they came on.

## Consequences

- A crashed `flai serve` comes back by itself, and the dashboard can restart it and upgrade flai without a shell.
- One host serves the machine, with one config file. `flai dashboard` under another config is refused with the host's config named, where two `flai serve` processes used to compete (I-0029).
- The machine has a new loopback listener on port 4241. A program already on that port stops the host from starting until `FLAI_HOST_ADDR` moves it.
- The `host` action gives a holder of the dashboard token the power to restart the operator's processes and to install a flai release with the operator's GitHub credentials. Like `push`, it is off until enabled in a shell.
- Operators who stopped `flai serve` with `flai serve stop` get the same effect, now kept until `flai serve start`. `flai host stop` ends everything.
- ADR-0034's rules for MCP ports, for a server started by hand, and for `--exit-with` stand. Only who starts the servers changes.

## Alternatives considered

- **A lock file for one per machine.** Needs clean-up after a crash, and a Windows lock behaves differently. The listening socket needs neither.
- **A Unix socket for the control API.** No port to collide, but the CLI and Windows would need a second path. Loopback HTTP with a token is what `flai mcp http` already does.
- **The host watches serve's registry itself** rather than being told. `flai serve` alone knows the projects it serves below a folder (ADR-0036), so it stays the one that decides.
- **Keeping MCP supervision in `flai serve`** and having the host supervise only serve. Then a restart of serve would bounce every agent's MCP server, and the story asks that serve no longer manage them.
