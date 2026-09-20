---
id: S-0077
type: story
nature: improvement
title: flai dashboard no longer mounts the clone, and everything that existed to make the mount safe goes with it
status: ready
parent: E-0003
owner: alex
created: 2026-09-20T07:26:51Z
updated: 2026-09-20T07:28:22Z
transitions:
  - to: ready
    at: 2026-09-20T07:28:22Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flaiover/src, flaiover/Dockerfile, design/system, design/adrs, docs/operators]
---
# S-0077 flai dashboard no longer mounts the clone, and everything that existed to make the mount safe goes with it

## Goal
With reads, writes, and MCP off the mount, `flai dashboard` stops giving the container the project. The container gets a port, its login token, and the agent credential, and nothing else of the host.

## Acceptance criteria
- [ ] `flai dashboard` passes no project volume, no guard mounts, no git identity, no excludes file, no push key, no generated passwd, and no `PROJECT_DIR`; the container no longer needs to run as the host user, and whether it still does is decided and recorded
- [ ] `dashboard.push_key` and `dashboard.push_known_hosts` are retired: set, they produce a message that says what replaced them; `flai dashboard status` drops the lines about them and about read-only git paths
- [ ] Tried with the published image run the new way: every page and every write works, a process in the container cannot see any file of the project, and the routes ADR-0027 closed and the ones it left open (tracked files, refs, ignored files the host executes) no longer exist
- [ ] An ADR supersedes ADR-0022, ADR-0026, and ADR-0027, and says what of ADR-0018 still holds; I-0022's successor risks are reviewed and closed or restated
- [ ] The operators' documentation is rewritten where it explains mounts, the push key, and what the container can write; the upgrade note says what an operator with a push key configured must do

## Tasks

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the writes story and the MCP story.
