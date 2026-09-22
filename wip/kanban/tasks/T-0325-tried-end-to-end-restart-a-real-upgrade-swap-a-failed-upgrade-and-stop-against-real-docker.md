---
id: T-0325
type: task
nature: feature
title: "Tried end to end: restart, a real upgrade swap, a failed upgrade, and stop, against real Docker"
status: backlog
parent: S-0081
owner: alex
created: 2026-09-22T21:13:07Z
updated: 2026-09-22T21:22:45Z
transitions: []
stream: S-0081
tags: []
---
# T-0325 Tried end to end: restart, a real upgrade swap, a failed upgrade, and stop, against real Docker

## Work
The shared container's name is now fixed (`sharedContainerName = "flaiover"`, S-0080), not derived from a project, so it cannot be given a different name to avoid colliding with the operator's real dashboard container, which is also named `flaiover` and is running on this host's real Docker daemon right now (found while verifying T-0322). The scratch lab must therefore run against its own, isolated Docker daemon — a `docker:dind` container, or an equivalent separate daemon/context — never the host's real one, so `docker ps`, `docker stop`, and every other call this story's commands make can only ever see the scratch lab's own containers.

In that isolated daemon (own `flai serve`, a scratch project, throwaway token): register the project, enable the `dashboard` action, then from the `/host` page (Docker + Playwright):
1. Restart while the project list and board are open elsewhere; confirm the page shows the action running, then reconnected, and the other pages recover on their own.
2. Check for updates against a target that genuinely differs (build `flaiover:s81a` and `flaiover:s81b` from two different commits, or two tags of the same build with a file changed) and confirm it reports one available without changing the running container.
3. Upgrade against that same target: confirm the swap happens, the new version shows, and the previous image is gone from `docker ps` (a new container, not the old one restarted).
4. Force an upgrade failure (a target image that does not serve `/_health`, e.g. a bare `alpine sleep` tag under the same name) and confirm the previously running container is still running, unchanged, and the page says the upgrade failed and why.
5. Stop from the page; confirm it behaves exactly as `flai dashboard stop` (unregisters, only stops the container if it was the last project).

Tick the acceptance criteria against what was actually seen, not the design. Confirm the operator's real `flai serve` and the host's real `flaiover` container were untouched throughout (they are on a different daemon entirely, but check anyway). Remove every scratch resource afterward, including the isolated daemon container itself.

## Done when
All five checks pass against real Docker, inside the isolated daemon; the story's three acceptance criteria are ticked from what was verified; scratch resources, including the isolated daemon, are confirmed gone.

## Notes
Depends on T-0322–T-0324. This is where "an upgrade that fails leaves the previous container running" gets proven against real Docker, not just the unit tests' fake runner.
