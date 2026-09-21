---
id: E-0003
type: epic
nature: feature
title: flaiover dashboard
status: backlog
owner: alex
created: 2026-09-15T16:09:00Z
updated: 2026-09-21T04:09:25Z
transitions: []
tags: []
---

# E-0003 flaiover dashboard

## Outcome
A SvelteKit and Tailwind dashboard, published as a Docker image, that renders all documentation, searches design and wip, shows the kanban board, and charts flow metrics for a mounted conforming repo.

## Stories
- S-0011 SvelteKit scaffold and repo reader API
- S-0012 Documentation explorer and search
- S-0013 Kanban board view with transitions
- S-0014 Flow metric charts
- S-0015 Docker image build and publish
- S-0032 flaiover observability: request logs, health, metrics, traces
- S-0035 flai dashboard runs from a private registry or a local build
- S-0044 Dashboard theme from the brand palette
- S-0045 Charts do not update when switching chart type
- S-0071 Research how the dashboard and a flai process on the host can talk, so that actions in the dashboard trigger actions on the host
- S-0072 flai on the host dials the dashboard and the two speak JSON-RPC: the channel and nothing else
- S-0073 The dashboard reads work items through the channel: project, board, items, threads, and file changes
- S-0074 The dashboard reads documents, ADRs, narratives, the inbox, and search through the channel, and parses no project file
- S-0075 Every write the dashboard makes goes through the channel, and the image no longer carries flai, git, or ssh
- S-0076 flai serves MCP itself on the host, over HTTP as well as stdio, and the dashboard's /mcp endpoint goes
- S-0077 flai dashboard no longer mounts the clone, and everything that existed to make the mount safe goes with it
- S-0078 An acceptance made from the board is pushed and published by flai on the host, when the operator has enabled it
- S-0079 A story moved to ready starts an agent on the host, with the command the operator configured
- S-0080 One dashboard serves every project the host flai serves
- S-0081 The dashboard can restart, upgrade, and stop itself through flai on the host
- S-0082 Checks for a story in review are run on the host and shown on the review page
- S-0083 Logging in to the dashboard works from any address, not only localhost: the session cookie is not marked Secure over plain HTTP
- S-0084 One flai serve per repository: a second one steps back instead of taking the dashboard's connection
- S-0085 Stories and epics are editable from the dashboard: title, nature, tags, touches, parent, and body, each change made by flai on the host
- S-0086 The host flai connection banner says 'not connected' even when a flai is connected but too old
- S-0087 Moving a story to done merges it; a publish button on the board tags and pushes everything accumulated since the last one
- S-0088 The Review page doesn't render markdown in the narrative pane
- S-0089 Inbox should not continue showing open items for stories that moved to review
- S-0090 Operators should be able to answer open questions in via the dashboard
- S-0091 The flai CLI should install itself to a user's home directory or a configurable path

## Notes
Defined from the brief in the root CLAUDE.md.
