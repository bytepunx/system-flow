---
id: S-0079
type: story
nature: feature
title: A story moved to ready starts an agent on the host, with the command the operator configured
status: done
parent: E-0003
owner: alex
created: 2026-09-20T07:26:51Z
updated: 2026-09-20T19:06:55Z
transitions:
  - to: ready
    at: 2026-09-20T12:10:05Z
    by: alex
  - to: in-progress
    at: 2026-09-20T15:47:05Z
    by: system-flow
  - to: review
    at: 2026-09-20T16:01:21Z
    by: system-flow
  - to: done
    at: 2026-09-20T19:06:55Z
    by: alex
tags: [cli, dashboard]
touches: [flai/cmd, flai/internal, flaiover/src, design/system, design/conventions, template]
---
# S-0079 A story moved to ready starts an agent on the host, with the command the operator configured

## Goal
Ready means work starts. When a story becomes ready and no agent session is attending the project, flai on the host starts one with a command the operator wrote in the host's configuration. The dashboard supplies the story's ID and nothing else.

## Acceptance criteria
- [x] The action is off until enabled and has no default command; the command is an argument list in the host's flai configuration, run in the project's directory with `FLAI_AGENT` and the story's ID in its environment, never through a shell, and never built from anything the dashboard sent beyond a validated story ID
- [x] It fires when a story enters ready (from the board, the CLI, or an agent) and the WIP limit allows a pull; whether an agent is already attending is judged from something flai can know (an MCP session holding `wait_for_events`, a fresh narrative log entry), and the rule is written down
- [x] At most one agent is started per project at a time; a second ready story waits for the first session to end or for a limit the operator sets; a command that fails to start is reported on the board and journalled
- [x] The board shows that an agent was started for a story, by which command's name, and when; stopping a started agent from the dashboard is out of scope and said so
- [x] The conventions say what an agent started this way does first; baselines change in `template/` first
- [x] Tried with a scratch project and a stub command; never with the operator's real agent command

## Tasks
- T-0301 The agent action and its command live in the host's configuration, off and empty by default
- T-0302 flai serve starts the agent when a story enters ready, nobody is attending, and the limit allows a pull
- T-0303 The board shows that an agent was started for a story, by which command, and when, or why not
- T-0304 The conventions say what an agent started by flai does first, and the operators' page says what enabling it means
- T-0305 Tried with a scratch project and a stub command

## Notes
From S-0071's finding, `design/system/dashboard-host-channel.md`, and ADR-0029. Depends on the channel and the work-item reads (for the change notifications). This is the gap the operator hit on 2026-09-20: a move to ready changed a file and nothing on the host was told.
