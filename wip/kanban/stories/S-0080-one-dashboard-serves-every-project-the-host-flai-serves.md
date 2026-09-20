---
id: S-0080
type: story
nature: feature
title: One dashboard serves every project the host flai serves
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:52Z
updated: 2026-09-20T12:10:11Z
transitions:
  - to: ready
    at: 2026-09-20T12:10:11Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, design/adrs, docs/operators]
---
# S-0080 One dashboard serves every project the host flai serves

## Goal
With no mount, nothing in the container belongs to one project. One dashboard container per user shows every project registered with the host flai: a project switcher, routes by project key, one login. This is the single pane E-0007 wanted, without a block of ports, a registry of tokens, or embedded pages.

## Acceptance criteria
- [ ] The host flai announces the projects it serves and the dashboard routes every page and API call by project key; a key flai does not serve is refused by flai, not only by the dashboard
- [ ] `flai dashboard` in any project makes sure the one container and the host flai run and that the project is registered, and prints that project's address; `flai dashboard stop` in a project unregisters it and stops the container only when it was the last
- [ ] One login token per user, kept in flai's home and not in each project's `.flai-cache`; an ADR refines ADR-0018 and says what happens to existing project tokens
- [ ] A project list with a filter, each project's state at a glance (stories in review, threads awaiting the designer, whether an agent is attending), and the current project remembered per browser
- [ ] Notifications, SSE, and caches are per project; one project's watcher or a slow answer does not stall another's pages
- [ ] The operators' and users' documentation are updated, including the upgrade from one container per project

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the mount story. E-0007 and its stories were cancelled on 2026-09-20; this replaces them.
