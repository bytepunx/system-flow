---
id: E-0007
type: epic
nature: feature
title: flairport - the interface for managing multiple flai projects
status: backlog
owner: alex
created: 2026-09-20T01:14:04Z
updated: 2026-09-20T03:35:30Z
transitions: []
tags: [flairport]
touches: [flairport/src]
---
# E-0007 flairport - the interface for managing multiple flai projects

## Outcome

- Gather all running flai-managed project instances and their web dashboard containers running on a single host and present them via a "single pane of glass" so the operator can have one page to go to and click through all the projects they're managing

- Expose this surface via something like cloudflared or a similar tunnel so that it can be accessed from outside the network over HTTPS (https terminated at the tunnel provider's network edge) with a simple passkey/identity check at the edge

## Stories

- S-0065 Extend flai and flaiover to write and manage a metadata file about each running instance into ~/.flai (home directory belongs to the current operator) in YAML format - this includes the port for each flaiover instance as well as the necessary token to access it

- S-0066 flai/flaiover will need to negotiate a port depending on what's available. this can be part of the metadata read/write pattern to a ~/.flai folder such that flai reserves a port range block of ~50-100 ports that are unlikely to collide with common OSS/system ports

- S-0067 Read the currently open flai/flaiover sessions from the metadata files and present a sidebar interface (which should include a filter bar to support finding projects quickly) so that the operator using flairport can start, stop, restart, and manage all projects using this system-flow approach

- S-0068 On selection of a project, present it, embedded, in the browser screen as a "integration at the glass" approach
- S-0069 Research tunnel options for reaching flairport from outside the network, ending in a story for flai CLI integration if one suits

## Notes

It is not flairports job to provide an authentication/authorization mechanism; that is additional infrastructure that could be supplied either by the tunnel vendor or by authstar.
