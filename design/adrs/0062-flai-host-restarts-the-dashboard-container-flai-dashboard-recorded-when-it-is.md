---
id: ADR-0062
title: "flai host restarts the dashboard container flai dashboard recorded when it is gone or not answering, unless dashboard.no_restart is set"
status: accepted
date: 2026-10-01
supersedes: []
superseded_by: []
refines: [ADR-0040]
---

# ADR-0062 flai host restarts the dashboard container flai dashboard recorded when it is gone or not answering, unless dashboard.no_restart is set

## Context

The dashboard container runs `docker run --detach --rm` with no restart policy. A pattern kill meant for a scratch server once stopped its process, because a container's processes match `pkill` on the host (I-0025). S-0170 made the server exit cleanly on SIGTERM and made a stop signal only the process flai started ([ADR-0058](0058-the-agent-host-action-lets-the-dashboard-stop-a-story-s-agent-and-a-stop.md)). A stray signal still took the dashboard down until the operator ran `flai dashboard` again. `flai dashboard status` said only whether docker ran the container, not whether it answered.

flai host ([ADR-0040](0040-one-flai-host-per-machine-runs-flai-serve-and-each-project-s-mcp-server-as-its.md)) already keeps `flai serve` and the MCP servers running. The operator asked in S-0184 for the dashboard to recover the same way, with a back-off, a log line that says why, and a way to turn it off.

## Decision

flai host restarts the dashboard container flai dashboard recorded when it is gone or not answering, unless `dashboard.no_restart` is set.

1. **The record.** Each time `flai dashboard`, `flai dashboard restart`, or `flai dashboard upgrade` starts the shared container, it writes the container's name, image reference, and publish to `dashboard-container.json` beside the dashboard token. `flai dashboard` also records a container an older flai started, as it runs. `flai dashboard stop` forgets the record when it leaves no project registered, so a dashboard the operator stopped is not brought back.
2. **The look.** Every 15 seconds the host asks docker whether the recorded container runs, and its published port whether `/_health` answers. The answer is running, not answering, or gone. Without a record, or with `dashboard.no_restart` set, the dashboard is not watched. The configuration is read at every look.
3. **The restart.** After two looks in a row that do not find it running, the host removes what is left of the container and starts it again from the recorded image reference. It never pulls, so a floating tag cannot change what runs. A swap that `flai dashboard restart` or `upgrade` is making is shorter than one look.
4. **The back-off.** After a restart, the next waits 30 seconds, doubled after each restart up to five minutes. It is reset once the dashboard has answered for a minute after a restart.
5. **Said.** Each restart is logged at warn with `reason` (gone or not-answering). A restart that fails is logged at error with the error. `flai host status` shows the watch: its state, how many restarts, the last one's time and reason, and its error.
6. **Health.** The image has a `HEALTHCHECK` on `/_health`. `flai dashboard status` reports running, not answering, or gone from its own probe, with docker's health verdict beside it.

## Consequences

- A stray signal costs the dashboard 15 to 30 seconds, not an outage that lasts until someone notices.
- A `docker stop flaiover` by hand is undone within half a minute. To keep it stopped, run `flai dashboard stop` or set `dashboard.no_restart`.
- A dashboard that crashes at start is restarted at most every five minutes, and each failure is in the host's log.
- Only a dashboard on this host's docker is watched, and only while flai host runs. `flai dashboard --no-serve` starts no host, so nothing watches it.

## Alternatives considered

- **Docker's `--restart unless-stopped`.** Docker refuses it together with `--rm`, so the container would outlive its stops and need cleaning up. Docker also restarts on exit only, never for a container that runs and does not answer.
- **flai serve watches it.** flai serve is the process the dashboard talks to. When flai serve goes, the host restarts it, and the host is what already keeps processes running.
- **Restart from the configuration's image and tag.** A dashboard built with `--build`, or started with `--tag`, would come back as something else. A floating tag would pull an upgrade nobody asked for.
