---
id: S-0106
type: story
nature: feature
title: The flai host command should become a process-managing shell
status: backlog
owner: alex
created: 2026-09-24T01:23:37Z
updated: 2026-09-24T01:23:37Z
transitions: []
tags: [cli]
touches: [flai/cmd]
---
# S-0106 The flai host command should become a process-managing shell

## Goal

A new command, `flai host`, is a bootstrap/loader that manages the other processes. For this to work, there needs to be a line of communication from `flai serve` process back to the the host.

The goal is to enable the addition of extensions to the dashboard's host panel that show the status of the other managed processes (flai serve and flai MCP) and enable management commands that can update the flai version and restart those sub-processes.

## Acceptance criteria
- [ ] flai host starts and manages child processes for flai serve and flai MCP
- [ ] flai serve is able to communicate commands back to the host in order to get status, version, issue upgrade commands, issue restart commands
- [ ] flai serve is no longer responsible for managing MCP or dashboard processes directly but instead sends the signal up to the host process which then manages the other sub-processes
- [ ] the flai host process stops all flai sub-processes when it crashes or stops (no zombie processes)
- [ ] the flai dashboard command no longer starts flai serve or flai MCP processes but instead boots the host process which handles the other two
- [ ] one flai host process can run per machine

## Tasks

## Notes
