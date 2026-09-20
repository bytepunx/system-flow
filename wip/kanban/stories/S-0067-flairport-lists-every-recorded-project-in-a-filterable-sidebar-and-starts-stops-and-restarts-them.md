---
id: S-0067
type: story
nature: feature
title: flairport lists every recorded project in a filterable sidebar and starts, stops, and restarts them
status: ready
parent: E-0007
owner: alex
created: 2026-09-20T01:18:59Z
updated: 2026-09-20T03:31:59Z
transitions:
  - to: ready
    at: 2026-09-20T03:31:59Z
    by: alex
tags: [flairport]
touches: [flairport/src, flai/cmd, system-flow.yaml]
---
# S-0067 flairport lists every recorded project in a filterable sidebar and starts, stops, and restarts them

## Goal
flairport exists: one page on the host that lists every flai project recorded in flai's home, shows which dashboards are running, and lets the operator start, stop, and restart them and find one quickly, so that managing several system-flow projects starts from a single place.

## Acceptance criteria
- [ ] flairport is a new sub-project, `flairport/`, listed in `system-flow.yaml` with its kind and tags so it is a release component, with design in `design/system/flairport.md`, its technology in `design/tech`, scripts and Makefile targets like flaiover's, and CI; how it runs is decided first and recorded in an ADR: a process on the host started by flai (it needs docker, flai, flai's home, and the project directories), not a container that is handed the docker socket, unless the ADR finds a better answer
- [ ] It reads the instance records (S-0065) and shows a sidebar of projects: name, key, state (running, stopped, unreachable, stale record), port, and flaiover version, with the state checked against docker and the dashboard's `/_health`, refreshed without a reload
- [ ] A filter bar above the list narrows it as the operator types, by name, key, and path, works from the keyboard, and keeps working with fifty projects
- [ ] Start, stop, and restart act through flai in the project's directory (`flai dashboard`, `flai dashboard stop`), never by driving docker directly, show flai's own message when it refuses (a push key with a passphrase, a port taken, an image that cannot be pulled), and the list follows the result
- [ ] "Manage" covers at least: open the project's dashboard, see what its container holds (a push key or no credential, whether git is read-only in it, from `flai dashboard status --json`), and forget a stale record; anything more is listed in the notes as later work rather than built
- [ ] flairport has no authentication of its own, as the epic says, and is safe to run that way: it binds to loopback only by default, refuses to start on another address without an explicit flag that names the risk, and its documentation says that it holds every project's token and can start and stop containers, so whatever is put in front of it (the tunnel's identity check, authstar) is the only lock. Project tokens are never sent to the browser
- [ ] It looks like flaiover's sibling: the same theme tokens, light and dark, and the same accessibility bar; tried in a browser with at least three scratch projects, never the operator's real `~/.flai`
- [ ] Tests for the record reader, the state check, the filter, and the actions with a fake flai; `docs/users` gains a flairport page and `docs/operators/index.md` a section on running it

## Tasks

## Notes
From E-0007, the operator's third bullet: "Read the currently open flai/flaiover sessions from the metadata files and present a sidebar interface (which should include a filter bar to support finding projects quickly) so that the operator using flairport can start, stop, restart, and manage all projects using this system-flow approach."

The epic's note: "It is not flairports job to provide an authentication/authorization mechanism; that is additional infrastructure that could be supplied either by the tunnel vendor or by authstar." The sixth criterion is what that costs: a page with no lock of its own that can reach every project must not be reachable by accident.

Depends on S-0065 (the records) and is much better after S-0066 (so several dashboards can run at once without the operator choosing ports). The embedded view of a selected project is S-0068; this story may simply open the project's dashboard in a new tab.

Relation to ADR-0024. That ADR sketched a hub that projects dial out to over a websocket, for projects on different machines. flairport is the single-host case the operator asked for first: it reads a local registry and needs no dial-out. The ADR for this story should say how the two relate, and whether flairport is that hub's first form or a separate thing.

`CLAUDE.md` says not to create top-level folders beyond those listed and the sub-projects in `system-flow.yaml`, and that flai owns the manifest's `projects`: adding `flairport` there is part of this story and is done the way the manifest allows, with the operator's word recorded, since it adds a release component.

Stack. flaiover is SvelteKit with the node adapter, Tailwind, and a token layer in `layout.css`; reusing that stack and sharing the theme is the obvious choice, and the ADR should say so or say why not. Whether flairport's server is Node (like flaiover) or a `flai` subcommand that serves a static front end is the main thing to settle: the second needs no Node on the host.
