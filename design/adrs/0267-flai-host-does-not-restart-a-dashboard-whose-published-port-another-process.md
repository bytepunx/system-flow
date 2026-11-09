---
id: ADR-0267
title: "flai host does not restart a dashboard whose published port another process holds, and names the docker-proxy that holds it"
status: accepted
date: 2026-11-09
supersedes: []
superseded_by: []
refines: [ADR-0062]
topics: [cli, dashboard, server-side]
---

# ADR-0267 flai host does not restart a dashboard whose published port another process holds, and names the docker-proxy that holds it

## Context

flai host restarts the dashboard container when it is gone or not answering ([ADR-0062](0062-flai-host-restarts-the-dashboard-container-flai-dashboard-recorded-when-it-is.md)). It judges "not answering" from `docker port <name> 3000/tcp` and a probe of `/_health` at the address that gives. When `docker port` gives nothing, the container counts as not answering.

[I-0243](../issues/I-0243-flai-host-cannot-see-a-dashboard-container-docker-lost-track-of-which-keeps-the.md) shows where that is wrong. A flaiover container started on 2026-10-25 dropped out of Docker's list, even under `docker ps --all`. Its processes kept running: its `containerd-shim`, `tini`, `node server.js`, and the `docker-proxy` that held `0.0.0.0:4242` for `172.17.0.2:3000`. From then on, every container flai host started had an empty port map, because the port was taken. Each one got the same bridge address, `172.17.0.2`, because Docker's address management no longer knew the orphan. The host found each new container not answering and restarted it every five minutes, 1,697 times in six days. Browsers and flai serve reached the orphan the whole time.

A restart cannot win the port back while another process holds it. flai host can tell that this is the case, and it can say what holds the port.

## Decision

flai host does not restart a dashboard whose published port another process holds. It says what holds the port and how to free it.

1. **The look.** The watch has a fourth answer beside running, not answering, and gone: `port-held`. It is `port-held` when all three of these are true: docker runs the recorded container, `docker port <name> 3000/tcp` gives no address for it, and the host port in the record's `publish` accepts a TCP connection all the same. A wildcard bind is probed at `127.0.0.1`. When the port does not accept a connection, the answer is not answering, as before.
2. **No restart.** A `port-held` look is not a miss. It sets the misses back to zero and leaves the back-off as it stands.
3. **Said once.** The watch logs `port-held` at error once each time it enters the state, not at every look, naming the holder. `flai host status` shows the state and the holder.
4. **The holder.** On Linux, flai reads `/proc/<pid>/cmdline` of each process and takes the `docker-proxy` whose arguments carry `-host-port <port>`. It names its PID and its `-container-ip` and `-container-port`. A process's command line is readable by any user. `ss -p` and `/proc/<pid>/fd` show who owns a socket only to root when the owner runs as root, as `docker-proxy` does.
5. **No holder found.** Where no `docker-proxy` names the port, or there is no `/proc`, flai says only that a process other than the container holds the port. It names `ss -ltnp` as the way to find it.
6. **What to do.** `flai dashboard status` and `flai host status` say to stop that process, or restart the Docker daemon. Once the holder is gone, `docker port` publishes nothing or the port stops answering. Two looks then find the dashboard not answering or gone, and the watch restarts it on its own.

## Consequences

- An orphaned container no longer causes a restart every five minutes. The operator learns what is wrong from one log line and the status commands.
- The dashboard browsers reach meanwhile is the orphan, which may run an older release with older settings. flai says so, but it does not fix it.
- Freeing the port stays the operator's. It usually needs root.
- A container that publishes nothing for another reason, while something unrelated listens on the same port, is reported as port held, not restarted. A restart could not publish that port either.
- The holder is named only on Linux. Elsewhere the status says a process holds the port, without naming it.

## Alternatives considered

- **flai kills the `docker-proxy`.** It runs as root, so flai cannot without root, and killing it leaves the orphan's other processes running. Choosing what to stop on the host is the operator's.
- **flai restarts the Docker daemon.** That stops every container on the host and needs root.
- **Publish the dashboard on another port.** Browsers, bookmarks, and `dashboard.port` name the recorded port. Moving it silently would leave the orphan answering there.
- **Keep restarting, with a longer back-off.** Every restart fails the same way, so it only spaces out the noise.
